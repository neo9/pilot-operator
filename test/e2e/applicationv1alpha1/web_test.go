package applicationv1alpha1

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/operator-framework/operator-sdk/pkg/test"
	"github.com/operator-framework/operator-sdk/pkg/test/e2eutil"

	"github.com/neo9/pilot-operator/test/e2e/helpers"
)

func TestSimpleWeb(t *testing.T) {
	list := getSampleList()
	namespace, ctx := helpers.GetClusterContext(t, &list)
	defer ctx.Cleanup()

	// get global framework variables
	f := test.Global
	application := getSampleWebApplication(namespace)
	err := f.Client.Create(context.TODO(), &application, &test.CleanupOptions{
		TestContext: ctx,
		Timeout: helpers.Timeout,
		RetryInterval: helpers.RetryInterval,
	})
	if err != nil {
		t.Fatal(err)
	}

	err = e2eutil.WaitForDeployment(t, f.KubeClient, namespace, application.Name, 1, time.Second*5, time.Second*30)

	application.Spec.Replicas = 4
	err = f.Client.Update(context.TODO(), &application)
	if err != nil {
		t.Fatal(err)
	}

	err = e2eutil.WaitForDeployment(t, f.KubeClient, namespace, application.Name, 4, time.Second*5, time.Second*30)
	if err != nil {
		t.Fatal(err)
	}
	ctx.Cleanup()
}

func TestLabels(t *testing.T) {
	list := getSampleList()
	namespace, ctx := helpers.GetClusterContext(t, &list)
	defer ctx.Cleanup()

	// get global framework variables
	f := test.Global
	application := getSampleNginxApplication(namespace, "1.16.0")
	application.Spec.Labels = map[string]string{"prometheus": "true"}
	err := f.Client.Create(context.TODO(), &application, &test.CleanupOptions{
		TestContext:   ctx,
		Timeout:       helpers.Timeout,
		RetryInterval: helpers.RetryInterval,
	})
	if err != nil {
		t.Fatal(err)
	}

	err = e2eutil.WaitForDeployment(t, f.KubeClient, namespace, application.Name, 1, time.Second*5, time.Second*30)
	deployment := getDeployment(f, application)
	expectedLabels := map[string]string{
		"controller": "pilot-operator",
		"name": "nginx",
		"prometheus": "true",
	}

	if reflect.DeepEqual(expectedLabels, deployment.Labels) {
		t.Fatal(errors.New("Labels do not match"))
	}

	application.Spec.Labels = map[string]string{"prometheus": "false"}
	err = f.Client.Update(context.TODO(), &application)
	if err != nil {
		t.Fatal(err)
	}
	err = e2eutil.WaitForDeployment(t, f.KubeClient, namespace, application.Name, 1, time.Second*5, time.Second*30)
	if err != nil {
		t.Fatal(err)
	}
	deployment = getDeployment(f, application)
	expectedLabels = map[string]string{
		"controller": "pilot-operator",
		"name": "nginx",
		"prometheus": "true",
	}

	if reflect.DeepEqual(expectedLabels, deployment.Labels) {
		t.Fatal(errors.New("Labels do not match after update"))
	}
}

