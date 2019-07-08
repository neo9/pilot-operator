package helpers

import (
	"github.com/operator-framework/operator-sdk/pkg/test"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/neo9/pilot-operator/pkg/apis"
	"github.com/operator-framework/operator-sdk/pkg/test/e2eutil"
	"testing"
	"time"
)

var (
	RetryInterval        = time.Second * 5
	Timeout              = time.Second * 60
	CleanupRetryInterval = time.Second * 1
	CleanupTimeout       = time.Second * 5
	operatorName         = "pilot-operator"
)

func GetClusterContext(t *testing.T, obj runtime.Object) (string, *test.TestCtx) {
	crdError := test.AddToFrameworkScheme(apis.AddToScheme, obj)
	if crdError != nil {
		t.Fatalf("failed to add custom resource scheme to framework: %v", crdError)
	}
	ctx := test.NewTestCtx(t)
	crError := ctx.InitializeClusterResources(&test.CleanupOptions{
		TestContext: ctx,
		Timeout: CleanupTimeout,
		RetryInterval: CleanupRetryInterval,
	})
	if crError != nil {
		t.Fatalf("failed to initialize cluster resources: %v", crError)
	}

	t.Log("Initialized cluster resources")
	namespace, err := ctx.GetNamespace()
	if err != nil {
		t.Fatal(err)
	}

	f := test.Global
	err = e2eutil.WaitForOperatorDeployment(t, f.KubeClient, namespace, operatorName, 1, RetryInterval, Timeout)
	if err != nil {
		t.Fatal(err)
	}

	return namespace, ctx
}
