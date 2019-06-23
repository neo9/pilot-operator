package e2e

import (
	"github.com/neo9/pilot-operator/pkg/apis"
	"context"
	"testing"
	"time"

	"github.com/operator-framework/operator-sdk/pkg/test"
	"github.com/operator-framework/operator-sdk/pkg/test/e2eutil"

	"github.com/neo9/pilot-operator/test/e2e/applicationv1alpha1"
)

var (
	retryInterval        = time.Second * 5
	timeout              = time.Second * 60
	cleanupRetryInterval = time.Second * 1
	cleanupTimeout       = time.Second * 5
)

func TestSimple(t *testing.T) {
	t.Logf("Testing nginx pod creation")
	sampleList := applicationv1alpha1.GetSampleList()
	//noinspection GoTypesCompatibility
	crdError := test.AddToFrameworkScheme(apis.AddToScheme, &sampleList)
	if crdError != nil {
		t.Fatalf("failed to add custom resource scheme to framework: %v", crdError)
	}
	ctx := test.NewTestCtx(t)
	defer ctx.Cleanup()
	crError := ctx.InitializeClusterResources(&test.CleanupOptions{TestContext: ctx, Timeout: cleanupTimeout, RetryInterval: cleanupRetryInterval})
	if crError != nil {
		t.Fatalf("failed to initialize cluster resources: %v", crError)
	}
	t.Log("Initialized cluster resources")
	namespace, err := ctx.GetNamespace()
	if err != nil {
		t.Fatal(err)
	}
	// get global framework variables
	f := test.Global
	// wait for operator to be ready
	err = e2eutil.WaitForOperatorDeployment(t, f.KubeClient, namespace, "pilot-operator", 1, retryInterval, timeout)
	if err != nil {
		t.Fatal(err)
	}
	application := applicationv1alpha1.GetSampleNginxApplication(namespace)
	//noinspection GoTypesCompatibility
	err = f.Client.Create(context.TODO(), &application, &test.CleanupOptions{TestContext: ctx, Timeout: timeout, RetryInterval: retryInterval})

	if err != nil {
		t.Fatal(err)
	}

	err = e2eutil.WaitForDeployment(t, f.KubeClient, namespace, "nginx", 1, time.Second*5, time.Second*30)

	application.Spec.Replicas = 4
	//noinspection GoTypesCompatibility
	err = f.Client.Update(context.TODO(), &application)
	if err != nil {
		t.Fatal(err)
	}

	err = e2eutil.WaitForDeployment(t, f.KubeClient, namespace, "nginx", 4, time.Second*5, time.Second*30)

	if err != nil {
		t.Fatal(err)
	}
}