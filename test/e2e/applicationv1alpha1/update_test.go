package applicationv1alpha1

import (
	"context"
	"testing"
	"time"

	"github.com/operator-framework/operator-sdk/pkg/test"
	"github.com/operator-framework/operator-sdk/pkg/test/e2eutil"

	"errors"
	"github.com/neo9/pilot-operator/pkg/apis/pilot/v1alpha1"
	"github.com/neo9/pilot-operator/test/e2e/helpers"
)

func TestSimpleUpdate(t *testing.T) {
	list := getSampleList()
	namespace, ctx := helpers.GetClusterContext(t, &list)
	defer ctx.Cleanup()

	// get global framework variables
	f := test.Global
	version := "1.16.0"
	targetVersion := "1.17.0"
	application := getSampleNginxApplication(namespace, version)
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

	application.Spec.HealthCheck.Path = "/"
	err = f.Client.Update(context.TODO(), &application)
	if err != nil {
		t.Fatal(err)
	}
	err = e2eutil.WaitForDeployment(t, f.KubeClient, namespace, application.Name, 1, time.Second*5, time.Second*30)
	if err != nil {
		t.Fatal(err)
	}
	deployment := getDeployment(f, application)
	if deployment.Spec.Template.Spec.Containers[0].LivenessProbe.HTTPGet.Path != "/" {
		t.Fatal(errors.New("health check was not updated"))
	}


	ctx.Cleanup()
}



func checkVersion(f *test.Framework, application v1alpha1.Application, version string) error {
	deployment := getDeployment(f, application)
	expectedImage := "nginx:" + version
	if deployment.Spec.Template.Spec.Containers[0].Image != expectedImage {
		return errors.New("Image should not have been updated to " + expectedImage)
	}

	return nil
}