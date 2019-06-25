package application

import (
	pilotv1alpha1 "github.com/neo9/pilot-operator/pkg/apis/pilot/v1alpha1"
	v1 "k8s.io/api/core/v1"
	"reflect"
)

func isServiceUpdated(current *v1.Service, application *pilotv1alpha1.Application) bool {
	stateModifications := 0
	state := getService(application)

	if !reflect.DeepEqual(state.Spec.Ports, current.Spec.Ports) {
		log.Info("Service ports should be updated")
		current.Spec.Ports = state.Spec.Ports
		stateModifications++
	}

	return stateModifications > 0
}

