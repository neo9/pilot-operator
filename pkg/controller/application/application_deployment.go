package application

import (
	pilotv1alpha1 "github.com/neo9/pilot-operator/pkg/apis/pilot/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	appsv1 "k8s.io/api/apps/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func  (r *ReconcileApplication) newDeploymentForCR(application *pilotv1alpha1.Application) *appsv1.Deployment {
	labels := map[string]string{
		"app": application.Name,
		"controller": "pilot",
	}

	var replicas int32 = 1


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
	controllerutil.SetControllerReference(application, dep, r.scheme)
	return dep
}

func getDeploymentContainers(application *pilotv1alpha1.Application) []corev1.Container {
	return []corev1.Container{
		{
			Name:    application.Name,
			Image:   application.Spec.Repository + ":" + application.Spec.Tag,
			Env: getDeploymentContainerEnvs(application),
		},
	};
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
	for _, secret := range application.Spec.Secrets {
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
