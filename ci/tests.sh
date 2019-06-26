#!/usr/bin/env sh

namespace="test-$(git rev-parse --short HEAD)"

kubectl create namespace $namespace
operator-sdk test local ./test/e2e --up-local --namespace $namespace
exit_code=$?

kubectl delete namespace $namespace

exit $exit_code
