package application

import (
	"context"
	v1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	pilotv1alpha1 "github.com/neo9/pilot-operator/pkg/apis/pilot/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
)

func (r *ReconcileApplication) ServiceReconcile(request reconcile.Request, application *pilotv1alpha1.Application) (reconcile.Result, error) {
	reqLogger := log.WithValues("ApplicationKind", "Service", "Request.Namespace", request.Namespace, "Request.Name", request.Name)
	// Check if this Deployment already exists
	found := &v1.Service{}
	err := r.client.Get(context.TODO(), types.NamespacedName{Name: application.Name, Namespace: application.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating a new Service", "Service.Namespace", application.Namespace, "Service.Name", application.Name)
		dep := r.newServiceForCR(application)
		err = r.client.Create(context.TODO(), dep)
		if err != nil {
			return reconcile.Result{}, err
		}

		return reconcile.Result{}, nil
	} else if err != nil {
		return reconcile.Result{}, err
	}

	// Check for changes
	if isServiceUpdated(found, application) {
		err = r.client.Update(context.TODO(), found)
		if err != nil {
			reqLogger.Error(err, "Could not update the Service")
			return reconcile.Result{}, err
		}

		reqLogger.Info("Skip reconcile: Service updated", "Service.Namespace", found.Namespace, "Service.Name", found.Name)
		return reconcile.Result{Requeue: true}, nil
	}

	reqLogger.Info("Skip reconcile: Service already exists", "Service.Namespace", found.Namespace, "Service.Name", found.Name)
	return reconcile.Result{}, nil
}
