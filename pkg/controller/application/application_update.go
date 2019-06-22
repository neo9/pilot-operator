package application

import (
	"fmt"
	appsv1 "k8s.io/api/apps/v1"
	pilotv1alpha1 "github.com/neo9/pilot-operator/pkg/apis/pilot/v1alpha1"
)

func isUpdated(found *appsv1.Deployment, application *pilotv1alpha1.Application) bool {
	replicas := application.Spec.Replicas
	if *found.Spec.Replicas != replicas {
		found.Spec.Replicas = &replicas
		log.Info(fmt.Sprintf("Replicas should be updated: %d -> %d", *found.Spec.Replicas, replicas))
		return true
	}

	if !isEnvEqual(found, application) {
		log.Info("Secrets envs differs")
		found.Spec.Template.Spec.Containers[0].Env = getDeploymentContainerEnvs(application)
		return true
	}

	return false
}

func isEnvEqual(found *appsv1.Deployment, application *pilotv1alpha1.Application) bool {
	env := found.Spec.Template.Spec.Containers[0].Env
	secrets := application.Spec.Secrets
	if len(env) != len(secrets) + 1 {
		log.Info("Secret length")
		return false
	}

	// TODO: compare values

	return true
}


