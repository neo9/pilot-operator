package scheduler

import (
	"context"
	"fmt"
	"github.com/go-logr/logr"
	"k8s.io/api/batch/v1beta1"
	"k8s.io/apimachinery/pkg/labels"
	"reflect"

	pilotv1alpha1 "github.com/neo9/pilot-operator/pkg/apis/pilot/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_scheduler")

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new Scheduler Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileScheduler{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("scheduler-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	err = c.Watch(&source.Kind{Type: &pilotv1alpha1.Scheduler{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	err = c.Watch(&source.Kind{Type: &v1beta1.CronJob{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &pilotv1alpha1.Scheduler{},
	})
	if err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcileScheduler implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileScheduler{}

// ReconcileScheduler reconciles a Scheduler object
type ReconcileScheduler struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

func getRequestLogger(request reconcile.Request) logr.Logger {
	return log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
}

func (r *ReconcileScheduler) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := getRequestLogger(request)
	reqLogger.Info("Reconciling Scheduler")

	// Fetch the Scheduler instance
	instance := &pilotv1alpha1.Scheduler{}
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}


	cronJobs := r.newCronJobListForCR(instance)
	create, err := r.createCronJobs(cronJobs, reqLogger)
	if err != nil || create {
		return reconcile.Result{Requeue: create}, err
	}

	result, err := r.updateCronJobs(cronJobs, instance, reqLogger)
	if err != nil {
		return reconcile.Result{}, err
	}

	return result, err
}

func (r *ReconcileScheduler) updateCronJobs(cronJobs []v1beta1.CronJob, instance *pilotv1alpha1.Scheduler, reqLogger logr.Logger) (reconcile.Result, error) {
	cronJobList := &v1beta1.CronJobList{}
	err := r.client.List(context.TODO(), &client.ListOptions{
		Namespace: instance.Namespace,
        LabelSelector: labels.SelectorFromSet(getRequestLabelCronJob(instance)),
	}, cronJobList)
	if err != nil {
		return reconcile.Result{}, err
	}

	reqLogger.Info("Checking update")

	for i := 0; i < len(cronJobList.Items); i++ {
		cronJob := cronJobList.Items[i]
		requestPos := getRequestInstancePosition(instance, cronJob.Name)
		if requestPos == -1 {
			reqLogger.Info(fmt.Sprintf("Delete CronJob request %s", cronJob.Name))
            err = r.client.Delete(context.TODO(), &cronJob)
            if err != nil && !errors.IsNotFound(err) {
				reqLogger.Error(err, fmt.Sprintf("Cannot delete CronJob request %s", cronJob.Name))
            	return reconcile.Result{}, err
			}
		} else {
			request := instance.Spec.Requests[requestPos]
			newCronJob := getRequestCronJob(instance, requestPos)
			// TODO: code refactor

			if cronJob.Spec.Schedule != request.Schedule {
				reqLogger.Info(fmt.Sprintf("Update schedule '%s' -> '%s'", cronJob.Spec.Schedule, request.Schedule), "Scheduler.Request", cronJob.Name)
				err = r.updateCronJob(&newCronJob)
				if err != nil {
					return reconcile.Result{}, err
				}
			} else if reflect.DeepEqual(cronJob.Spec.JobTemplate.Spec.Template.Spec.Containers[0].Args, newCronJob) {
				reqLogger.Info("Update command", "Scheduler.Request", cronJob.Name)
				err = r.updateCronJob(&newCronJob)
				if err != nil {
					return reconcile.Result{}, err
				}
			}
		}
	}

	return reconcile.Result{}, nil
}

func (r *ReconcileScheduler) updateCronJob(job *v1beta1.CronJob) error {
	return r.client.Update(context.TODO(), job)
}

func getRequestInstancePosition(instance *pilotv1alpha1.Scheduler, name string) int {
	for i := 0; i < len(instance.Spec.Requests); i++ {
		if getCronJobName(instance, i) == name {
			return i
		}
	}

	return -1
}

func (r *ReconcileScheduler) createCronJobs(cronJobs []v1beta1.CronJob, reqLogger logr.Logger) (bool, error) {
	create := false

	for i := 0; i < len(cronJobs); i++ {
		cronJob := cronJobs[i]
		found := &v1beta1.CronJob{}
		err := r.client.Get(context.TODO(), types.NamespacedName{Name: cronJob.Name, Namespace: cronJob.Namespace}, found)

		if err != nil && errors.IsNotFound(err) {
			reqLogger.Info("Creating a new CronJob", "CronJob.Namespace", cronJob.Namespace, "CronJob.Name", cronJob.Name)
			err = r.client.Create(context.TODO(), &cronJob)
			if err != nil {
				return false, err
			} else {
				create = true
			}
		}
	}

	return create, nil
}

