package application

import (
	"context"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	pilotv1alpha1 "github.com/neo9/pilot-operator/pkg/apis/pilot/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	appsv1 "k8s.io/api/apps/v1"
)

func (r *ReconcileApplication) DeploymentReconcile(request reconcile.Request, application *pilotv1alpha1.Application) (reconcile.Result, error) {
	reqLogger := log.WithValues("ApplicationKind", "Deployment", "Request.Namespace", request.Namespace, "Request.Name", request.Name)
	// Check if this Deployment already exists
	found := &appsv1.Deployment{}
	err := r.client.Get(context.TODO(), types.NamespacedName{Name: application.Name, Namespace: application.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating a new Deployment", "Deployment.Namespace", application.Namespace, "Deployment.Name", application.Name)
		dep := r.newDeploymentForCR(application)
		err = r.client.Create(context.TODO(), dep)
		if err != nil {
			return reconcile.Result{}, err
		}

		return reconcile.Result{}, nil
	} else if err != nil {
		return reconcile.Result{}, err
	}

	// Check for changes
	if isDeploymentUpdated(found, application) {
		err = r.client.Update(context.TODO(), found)
		if err != nil {
			reqLogger.Error(err, "Could not update the deployment")
			return reconcile.Result{}, err
		}

		reqLogger.Info("Skip reconcile: Deployment updated", "Deployment.Namespace", found.Namespace, "Deployment.Name", found.Name)
		return reconcile.Result{Requeue: true}, nil
	}

	reqLogger.Info("Skip reconcile: Deployment already exists", "Deployment.Namespace", found.Namespace, "Deployment.Name", found.Name)
	return reconcile.Result{}, nil
}
