package application

import (
	appsv1 "k8s.io/api/apps/v1"
	pilotv1alpha1 "github.com/neo9/pilot-operator/pkg/apis/pilot/v1alpha1"
	"fmt"
	"reflect"
)

func isUpdated(current *appsv1.Deployment, application *pilotv1alpha1.Application) bool {
	stateModifications := 0
	state := getDeployment(application)

	replicas := *state.Spec.Replicas
	if *current.Spec.Replicas != replicas {
		log.Info(fmt.Sprintf("Replicas should be updated: %d -> %d", *current.Spec.Replicas, replicas))
		*current.Spec.Replicas = replicas
		stateModifications++
	}

	stateContainer := state.Spec.Template.Spec.Containers[0]
	currentContainer := &current.Spec.Template.Spec.Containers[0]

	if stateContainer.Image != currentContainer.Image {
		log.Info(fmt.Sprintf("Image should be updated: %s -> %s", currentContainer.Image, stateContainer.Image))
		currentContainer.Image = stateContainer.Image
		stateModifications++
	}

	if !reflect.DeepEqual(stateContainer.LivenessProbe, currentContainer.LivenessProbe) {
		log.Info("Probes should be updated")
		currentContainer.LivenessProbe = stateContainer.LivenessProbe
		currentContainer.ReadinessProbe = stateContainer.ReadinessProbe
		stateModifications++
	}

	if !reflect.DeepEqual(stateContainer.Resources, currentContainer.Resources) {
		log.Info(
			fmt.Sprintf("Resources should be updated: CPU (%s -> %s), Memory (%s -> %s)",
				currentContainer.Resources.Requests.Cpu().String(),
				stateContainer.Resources.Requests.Cpu().String(),
				currentContainer.Resources.Requests.Memory().String(),
				stateContainer.Resources.Requests.Memory().String(),
			))
		currentContainer.Resources = stateContainer.Resources
		stateModifications++
	}

	if !reflect.DeepEqual(stateContainer.Env, currentContainer.Env) {
		log.Info("Containers env variables differs")
		currentContainer.Env = stateContainer.Env
		stateModifications++
	}

	return stateModifications > 0
}

