package e2e

import (
	"context"
	"testing"
	"time"

	"github.com/operator-framework/operator-sdk/pkg/test"
	"github.com/operator-framework/operator-sdk/pkg/test/e2eutil"

	"github.com/neo9/pilot-operator/test/e2e/applicationv1alpha1"
	"github.com/neo9/pilot-operator/test/e2e/helpers"
)

func TestSimpleWeb(t *testing.T) {
	namespace, ctx := helpers.GetClusterContext(t)

	// get global framework variables
	f := test.Global
	application := applicationv1alpha1.GetSampleWebApplication(namespace)
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