package application

import (
	"context"
	pilotv1alpha1 "github.com/neo9/pilot-operator/pkg/apis/pilot/v1alpha1"
	"k8s.io/api/batch/v1beta1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func (r *ReconcileApplication) CronJobReconcile(request reconcile.Request, application *pilotv1alpha1.Application) (reconcile.Result, error) {
	reqLogger := getLogger(request.Namespace, request.Namespace, "CronJob")
	// Check if this Deployment already exists
	found := &v1beta1.CronJobList{}
	err := r.client.Get(context.TODO(), types.NamespacedName{Name: application.Name, Namespace: application.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating a new CronJob", "CronJob.Namespace", application.Namespace, "CronJob.Name", application.Name)
		cronJobsList := r.newCronJobListForCR(application)
		err = r.client.Create(context.TODO(), cronJobsList)
		if err != nil {
			return reconcile.Result{}, err
		}

		return reconcile.Result{}, nil
	} else if err != nil {
		return reconcile.Result{}, err
	}

	reqLogger.Info("Skip reconcile: CronJob already exists", "CronJob.Namespace", application.Namespace, "CronJob.Name", application.Name)
	return reconcile.Result{}, nil
}
