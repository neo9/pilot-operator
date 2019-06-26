package e2e

import (
	"context"
	"testing"
	"time"

	"github.com/operator-framework/operator-sdk/pkg/test"
	"github.com/operator-framework/operator-sdk/pkg/test/e2eutil"

	"github.com/neo9/pilot-operator/test/e2e/applicationv1alpha1"
	"github.com/neo9/pilot-operator/test/e2e/helpers"
	"k8s.io/apimachinery/pkg/types"
	appsv1 "k8s.io/api/apps/v1"
	"github.com/neo9/pilot-operator/pkg/apis/pilot/v1alpha1"
	"errors"
)

func TestSimpleUpdate(t *testing.T) {
	namespace, ctx := helpers.GetClusterContext(t)
	defer ctx.Cleanup()

	// get global framework variables
	f := test.Global
	version := "1.16.0"
	targetVersion := "1.17.0"
	application := applicationv1alpha1.GetSampleNginxApplication(namespace, version)
	err := f.Client.Create(context.TODO(), &application, &test.CleanupOptions{
		TestContext: ctx,
		Timeout: helpers.Timeout,
		RetryInterval: helpers.RetryInterval,
	})
	if err != nil {
		t.Fatal(err)
	}

	err = e2eutil.WaitForDeployment(t, f.KubeClient, namespace, application.Name, 1, time.Second*5, time.Second*30)

	// Init tag does not modify tag (init tag is latest for the first deployment)
	application.Spec.Tag = version
	application.Spec.InitTag = targetVersion
	err = f.Client.Update(context.TODO(), &application)
	if err != nil {
		t.Fatal(err)
	}
	err = e2eutil.WaitForDeployment(t, f.KubeClient, namespace, application.Name, 1, time.Second*5, time.Second*30)
	if err != nil {
		t.Fatal(err)
	}

	err = checkVersion(f, application, version)
	if err != nil {
		t.Fatal(err)
	}

	application.Spec.Tag = targetVersion
	err = f.Client.Update(context.TODO(), &application)
	if err != nil {
		t.Fatal(err)
	}
	err = e2eutil.WaitForDeployment(t, f.KubeClient, namespace, application.Name, 1, time.Second*5, time.Second*30)
	if err != nil {
		t.Fatal(err)
	}

	err = checkVersion(f, application, targetVersion)
	if err != nil {
		t.Fatal(err)
	}


	ctx.Cleanup()
}


func checkVersion(f *test.Framework, application v1alpha1.Application, version string) error {
	deployment := &appsv1.Deployment{}
	f.Client.Get(context.TODO(), types.NamespacedName{Name: application.Name, Namespace: application.Namespace}, deployment)
	expectedImage := "nginx:" + version
	if deployment.Spec.Template.Spec.Containers[0].Image != expectedImage {
		return errors.New("Image should not have been updated to " + expectedImage)
	}

	return nil
}