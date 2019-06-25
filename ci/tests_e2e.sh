#!/usr/bin/env sh

namespace="pilot-op-test-2e2-$(git rev-parse --short HEAD)"

kubectl create namespace $namespace
operator-sdk test local ./test/e2e --namespace $namespace
exit_code=$?

kubectl delete namespace $namespace

exit $exit_code
