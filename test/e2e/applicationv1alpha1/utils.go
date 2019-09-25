package applicationv1alpha1

import (
	"context"
	"github.com/neo9/pilot-operator/pkg/apis/pilot/v1alpha1"
	"github.com/operator-framework/operator-sdk/pkg/test"
	"k8s.io/apimachinery/pkg/types"
	appsv1 "k8s.io/api/apps/v1"
)

func getDeployment(f *test.Framework, application v1alpha1.Application) *appsv1.Deployment {
	deployment := &appsv1.Deployment{}
	f.Client.Get(context.TODO(), types.NamespacedName{Name: application.Name, Namespace: application.Namespace}, deployment)
	return deployment
}
