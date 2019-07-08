package schedulerv1alpha1

import (
	"github.com/neo9/pilot-operator/pkg/apis/pilot/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func getSampleList() v1alpha1.SchedulerList {
	return v1alpha1.SchedulerList{
		TypeMeta: getTypeMeta(),
	}
}

func getSimpleScheduler(namespace string) v1alpha1.Scheduler {
	return v1alpha1.Scheduler{
		TypeMeta: getTypeMeta(),
		ObjectMeta: metav1.ObjectMeta{
			Name:      "simple-cronjob",
			Namespace: namespace,
		},
		Spec: v1alpha1.SchedulerSpec{
			Requests: []v1alpha1.SchedulerRequest{
				{
					Schedule: "* * * * *",
					Name: "default",
					Service: "google.com",
					Method: "get",
					Port: 80,
					Path: "/",
				},
			},
		},
	}
}

func getTypeMeta() metav1.TypeMeta {
	return metav1.TypeMeta{
		Kind:       "Scheduler",
		APIVersion: "pilot.neo9.fr/v1alpha1",
	}
}
