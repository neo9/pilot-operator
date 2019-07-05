package scheduler

import (
	"fmt"
	pilotv1alpha1 "github.com/neo9/pilot-operator/pkg/apis/pilot/v1alpha1"
	v1 "k8s.io/api/batch/v1"
	"k8s.io/api/batch/v1beta1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"strings"
)

func (r *ReconcileScheduler) newCronJobListForCR(scheduler *pilotv1alpha1.Scheduler) []v1beta1.CronJob {
	applicationRequests := scheduler.Spec.Requests

	cronJobs := make([]v1beta1.CronJob, len(applicationRequests))
	for i := 0; i < len(applicationRequests); i++ {
		cronJob := getRequestCronJob(scheduler, i)
		controllerutil.SetControllerReference(scheduler, &cronJob, r.scheme)
		cronJobs[i] = cronJob
	}

	return cronJobs
}

func getRequestLabelCronJob(scheduler *pilotv1alpha1.Scheduler) map[string]string{
	return map[string]string{
		"scheduler": scheduler.Name,
		"controller": "pilot",
	}
}

func getCronJobName(scheduler *pilotv1alpha1.Scheduler, index int) string {
	return fmt.Sprintf("%s-%s", scheduler.Name, scheduler.Spec.Requests[index].Name)
}

func getRequestCronJob(scheduler *pilotv1alpha1.Scheduler, index int) v1beta1.CronJob {
	request := scheduler.Spec.Requests[index]

	var historyLimit int32 = 5

	return v1beta1.CronJob{
		ObjectMeta: metav1.ObjectMeta{
			Name:      getCronJobName(scheduler, index),
			Namespace: scheduler.ObjectMeta.Namespace,
			Labels:    getRequestLabelCronJob(scheduler),
		},
		Spec: v1beta1.CronJobSpec{
			ConcurrencyPolicy: v1beta1.ForbidConcurrent,
			FailedJobsHistoryLimit: &historyLimit,
			SuccessfulJobsHistoryLimit: &historyLimit,
			Schedule: request.Schedule,
			JobTemplate: v1beta1.JobTemplateSpec{
				Spec: v1.JobSpec{
					Template: getRequestJobSpec(&request),
				},
			},
		},
	}
}

func getServicePort(requests *pilotv1alpha1.SchedulerRequest) int32 {
	var port int32 = 80
	if requests.Port != 0 {
		port = requests.Port
	}

	return port
}

func getRequestJobSpec(request *pilotv1alpha1.SchedulerRequest) corev1.PodTemplateSpec {
	method := "post"
	if request.Method != "" {
		method = request.Method
	}

	servicePort := getServicePort(request)

	return corev1.PodTemplateSpec{
		Spec: corev1.PodSpec{
			RestartPolicy: corev1.RestartPolicyOnFailure,
			Containers: []corev1.Container{
				{
					Name: request.Name,
					Image: "byrnedo/alpine-curl:0.1.7",
					Command: []string{
						"curl",
					},
					Args: []string{
						"-vis",
						"-X" + strings.ToUpper(method),
						fmt.Sprintf("http://%s:%d/%s", request.Service, servicePort, request.Path),
					},
					Resources: corev1.ResourceRequirements{
						Requests: corev1.ResourceList{
							"cpu":    resource.MustParse("1m"),
							"memory": resource.MustParse("20M"),
						},
						Limits: corev1.ResourceList{
							// CFS BUG: no CPU limit to avoid unnecessary throttling
							"memory": resource.MustParse("20M"),
						},
					},
				},
			},
		},
	}
}