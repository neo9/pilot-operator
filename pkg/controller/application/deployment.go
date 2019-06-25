package application

import (
	pilotv1alpha1 "github.com/neo9/pilot-operator/pkg/apis/pilot/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	appsv1 "k8s.io/api/apps/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func  (r *ReconcileApplication) newDeploymentForCR(application *pilotv1alpha1.Application) *appsv1.Deployment {
	dep := getDeployment(application)
	controllerutil.SetControllerReference(application, dep, r.scheme)
	return dep
}

func getDeployment(application *pilotv1alpha1.Application) *appsv1.Deployment {
	labels := map[string]string{
		"app": application.Name,
		"controller": "pilot",
	}

	var replicas int32 = 1
	if application.Spec.Replicas > 0 {
		replicas = application.Spec.Replicas
	}

	dep :=  &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      application.Name,
			Namespace: application.ObjectMeta.Namespace,
			Labels:    labels,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Containers: getDeploymentContainers(application),
				},
			},
		},
	}

	return dep
}

func getDeploymentContainers(application *pilotv1alpha1.Application) []corev1.Container {
	probe := getDeploymentProbe(application)

	return []corev1.Container{
		{
			Name:    application.Name,
			Image:   application.Spec.Repository + ":" + getDeploymentTag(application),
			Env: getDeploymentContainerEnvs(application),
			ReadinessProbe: &probe,
			LivenessProbe: probe.DeepCopy(),
			ImagePullPolicy: corev1.PullIfNotPresent,
			Resources: getDeploymentResources(application),
			Ports: getDeploymentPorts(application),
		},
	}
}

func getDeploymentPorts(application *pilotv1alpha1.Application) []corev1.ContainerPort {
	var port int32 = 0
	if application.Spec.Service.TargetPort != 0 {
		port = application.Spec.Service.TargetPort
	}

	return []corev1.ContainerPort{
		{
			Protocol: corev1.ProtocolTCP,
			Name: "http",
			ContainerPort: port,
		},
	}
}

func getDeploymentProbe(application *pilotv1alpha1.Application) corev1.Probe {
	var port int32 = 80
	if application.Spec.Service.TargetPort != 0 {
		port = application.Spec.Service.TargetPort
	} else if application.Spec.Service.Port != 0 {
		port = application.Spec.Service.Port
	}

	path := "/ping"
	if application.Spec.HealthCheck.Path != "" {
		path = application.Spec.HealthCheck.Path
	} else if application.Spec.Type == pilotv1alpha1.JAVA {
		path = "/actuator/health"
	} else if application.Spec.Type == pilotv1alpha1.WEB {
		path = "/health"
	}

	return corev1.Probe{
		Handler: corev1.Handler{
			HTTPGet: &corev1.HTTPGetAction{
				Path:  path,
				Port:  intstr.IntOrString{
					Type: intstr.Int,
					IntVal: port,
				},
				Scheme: "HTTP",
			},
		},
		FailureThreshold: 3,
		InitialDelaySeconds: 30,
		PeriodSeconds: 8,
		SuccessThreshold: 1,
		TimeoutSeconds: 1,
	}
}

func getDeploymentTag(application *pilotv1alpha1.Application) string {
	if application.Spec.InitTag == "" && application.Spec.Tag == "" {
		return "latest"
	} else if application.Spec.Tag != "" {
		return application.Spec.Tag
	}

	return application.Spec.InitTag
}

func getDeploymentResources(application *pilotv1alpha1.Application) corev1.ResourceRequirements {
	resources := application.Spec.Resources

	if resources.Memory == "" || resources.CPU == "" {
		if application.Spec.Type == pilotv1alpha1.WEB {
			resources.Memory = "20M"
			resources.CPU = "5m"
		} else if application.Spec.Type == pilotv1alpha1.NODEJS {
			resources.Memory = "186M"
			resources.CPU = "10m"
		} else {
			resources.Memory = "300M"
			resources.CPU = "10m"
		}
	}

	return corev1.ResourceRequirements{
		Requests: corev1.ResourceList{
			"cpu": resource.MustParse(resources.CPU),
			"memory": resource.MustParse(resources.Memory),
		},
		Limits: corev1.ResourceList{
			// CFS BUG: no CPU limit to avoid unnecessary throttling
			"memory": resource.MustParse(resources.Memory),
		},
	}
}

func getDeploymentContainerEnvs(application *pilotv1alpha1.Application) []corev1.EnvVar {
	namespaceEnv := corev1.EnvVar{
		Name: getDeploymentNamespaceEnvName(application.Spec.Type),
		ValueFrom: &corev1.EnvVarSource{
			FieldRef: &corev1.ObjectFieldSelector{
				APIVersion: "v1",
				FieldPath: "metadata.namespace",
			},
		},
	}

	envs := []corev1.EnvVar{namespaceEnv}
	for _, secret := range application.Spec.Pod.Secrets {
		envs = append(envs, corev1.EnvVar{
			Name: secret.Key,
			ValueFrom: &corev1.EnvVarSource{
				SecretKeyRef:  &corev1.SecretKeySelector{
					LocalObjectReference: corev1.LocalObjectReference{Name: secret.Name},
					Key:                  secret.Key,
				},
			},
		})
	}

	return envs
}

func getDeploymentNamespaceEnvName(deploymentType pilotv1alpha1.ApplicationType) string {
	if deploymentType == pilotv1alpha1.WEB || deploymentType == pilotv1alpha1.NODEJS {
		return "NODE_ENV"
	}

	return "PROFILE"
}
