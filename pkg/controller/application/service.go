package application

import (
	pilotv1alpha1 "github.com/neo9/pilot-operator/pkg/apis/pilot/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	v1 "k8s.io/api/core/v1"
)

func  (r *ReconcileApplication) newServiceForCR(application *pilotv1alpha1.Application) *v1.Service {
	dep := getService(application)
	controllerutil.SetControllerReference(application, dep, r.scheme)
	return dep
}

func getService(application *pilotv1alpha1.Application) *v1.Service {
	labels := map[string]string{
		"app": application.Name,
		"controller": "pilot",
	}

	return &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      application.Name,
			Namespace: application.ObjectMeta.Namespace,
			Labels:    getMergedLabels(labels, application.Spec.Labels),
		},
		Spec: v1.ServiceSpec{
			Type: v1.ServiceTypeClusterIP,
			Ports: []v1.ServicePort{
				{
					Name: "http",
					Port: getPort(application.Spec),
					TargetPort: getTargetPort(application.Spec),
					Protocol: v1.ProtocolTCP,
				},
			},
			Selector: labels,
		},
	}
}

func getPort(applicationSpec pilotv1alpha1.ApplicationSpec) int32 {
	if applicationSpec.Service.Port != 0 {
		return applicationSpec.Service.Port
	}

	return 80
}

func getTargetPort(applicationSpec pilotv1alpha1.ApplicationSpec) intstr.IntOrString {
	var targetPort int32 = 8080
	if applicationSpec.Service.TargetPort != 0 {
		targetPort = applicationSpec.Service.TargetPort
	}

	return intstr.IntOrString{
		Type: intstr.Int,
		IntVal: targetPort,
	}
}
