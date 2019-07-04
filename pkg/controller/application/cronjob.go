package application

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

func (r *ReconcileApplication) newCronJobListForCR(application *pilotv1alpha1.Application) *v1beta1.CronJobList {
	cronJobList := getRequestCronJobList(application)
	for i := 0; i < len(cronJobList.Items); i++ {
		controllerutil.SetControllerReference(application, &cronJobList.Items[i], r.scheme)
	}

	return cronJobList
}

func getRequestCronJobList(application *pilotv1alpha1.Application) *v1beta1.CronJobList {
	applicationRequests := application.Spec.CronJob.Requests

	cronJobs := make([]v1beta1.CronJob, len(applicationRequests))
	for i := 0; i < len(applicationRequests); i++ {
		cronJob := getRequestCronJob(application, i)
		cronJobs = append(cronJobs, cronJob)
	}

	return &v1beta1.CronJobList{
		Items: cronJobs,
	}
}

func getRequestCronJob(application *pilotv1alpha1.Application, index int) v1beta1.CronJob {
	request := application.Spec.CronJob.Requests[index]

	labels := map[string]string{
		"app":        application.Name,
		"controller": "pilot",
	}

	var historyLimit int32 = 5

	cronJob := v1beta1.CronJob{
		ObjectMeta: metav1.ObjectMeta{
			Name:      application.Name,
			Namespace: application.ObjectMeta.Namespace,
			Labels:    labels,
		},
		Spec: v1beta1.CronJobSpec{
			ConcurrencyPolicy: v1beta1.ForbidConcurrent,
			FailedJobsHistoryLimit: &historyLimit,
			SuccessfulJobsHistoryLimit: &historyLimit,
			Schedule: request.Schedule,
			JobTemplate: v1beta1.JobTemplateSpec{
                Spec: v1.JobSpec{
                    Template: getRequestJobSpec(&request, application.Name, getServicePort(application.Spec)),
				},
			},
		},
	}

	return cronJob
}

func getRequestJobSpec(request *pilotv1alpha1.ApplicationCronJobRequest, serviceName string, servicePort int32) corev1.PodTemplateSpec {
	method := "post"
	if request.Method != "" {
		method = request.Method
	}

    return corev1.PodTemplateSpec{
        Spec: corev1.PodSpec{
            RestartPolicy: corev1.RestartPolicyOnFailure,
            Containers: []corev1.Container{
            	{
         			Name: request.Name,
         			Image: "byrnedo/alpine-curl:0.1.7",
         			Command: []string{
         				"curl",
						"-vis",
         				"-X" + strings.ToUpper(method),
                        fmt.Sprintf("http://%s:%i/%s", serviceName, servicePort, request.Path),
					},
					Resources: corev1.ResourceRequirements{
						Requests: corev1.ResourceList{
							"cpu":    resource.MustParse("1m"),
							"memory": resource.MustParse("10M"),
						},
						Limits: corev1.ResourceList{
							// CFS BUG: no CPU limit to avoid unnecessary throttling
							"memory": resource.MustParse("10M"),
						},
					},
				},
			},
		},
	}
}

