package applicationv1alpha1

import (
	"github.com/neo9/pilot-operator/pkg/apis/pilot/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func getSampleList() v1alpha1.ApplicationList {
	return v1alpha1.ApplicationList{
		TypeMeta: getTypeMeta(),
	}
}

func getSampleWebApplication(namespace string) v1alpha1.Application {
	return v1alpha1.Application{
		TypeMeta: getTypeMeta(),
		ObjectMeta: metav1.ObjectMeta{
			Name:      "default-backend",
			Namespace: namespace,
		},
		Spec: v1alpha1.ApplicationSpec{
			Type:       v1alpha1.WEB,
			Repository: "k8s.gcr.io/defaultbackend-amd64",
			InitTag:    "1.5",
			HealthCheck: v1alpha1.ApplicationHealthCheck{
				Path: "/healthz",
			},
		},
	}
}

func getSampleNginxApplication(namespace string, version string) v1alpha1.Application {
	return v1alpha1.Application{
		TypeMeta: getTypeMeta(),
		ObjectMeta: metav1.ObjectMeta{
			Name:      "nginx",
			Namespace: namespace,
		},
		Spec: v1alpha1.ApplicationSpec{
			Type:       v1alpha1.WEB,
			Repository: "nginx",
			InitTag:    version,
			HealthCheck: v1alpha1.ApplicationHealthCheck{
				Path: "/",
			},
			Service: v1alpha1.ApplicationService{
				TargetPort: 80,
			},
		},
	}
}

func getTypeMeta() metav1.TypeMeta {
	return metav1.TypeMeta{
		Kind:       "Application",
		APIVersion: "pilot.neo9.fr/v1alpha1",
	}
}
