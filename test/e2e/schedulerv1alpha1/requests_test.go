package schedulerv1alpha1

import (
	"context"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"testing"
	"time"

	"github.com/operator-framework/operator-sdk/pkg/test"

	"github.com/neo9/pilot-operator/test/e2e/helpers"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestSimpleScheduler(t *testing.T) {
	list := getSampleList()
	namespace, ctx := helpers.GetClusterContext(t, &list)
	defer ctx.Cleanup()

	// get global framework variables
	f := test.Global
	scheduler := getSimpleScheduler(namespace)
	err := f.Client.Create(context.TODO(), &scheduler, &test.CleanupOptions{
		TestContext: ctx,
		Timeout: helpers.Timeout,
		RetryInterval: helpers.RetryInterval,
	})
	if err != nil {
		t.Fatal(err)
	}

	for _, request := range scheduler.Spec.Requests {
		err = waitForCronJob(t, f.KubeClient, namespace, scheduler.Name, request.Name, time.Second, time.Second * 5)
		if err != nil {
			t.Fatal(err)
		}
	}

	ctx.Cleanup()
}

func waitForCronJob(t *testing.T, kubeclient kubernetes.Interface, namespace string, name string, requestName string, retryInterval, timeout time.Duration) error {
	err := wait.Poll(retryInterval, timeout, func() (done bool, err error) {
		_, err = kubeclient.BatchV1beta1().CronJobs(namespace).Get(name + "-" + requestName, metav1.GetOptions{})
		if err != nil {
			if apierrors.IsNotFound(err) {
				t.Logf("Waiting for availability of %s CronJob\n", name)
				return false, nil
			}
			return false, err
		}

		return true, nil
	})

	if err != nil {
		return err
	}

	t.Logf("CronJob available %s\n", name)
	return nil
}

