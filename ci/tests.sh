#!/usr/bin/env sh

namespace="test-$(git rev-parse --short HEAD)-$TRAVIS_BUILD_ID"

kubectl create namespace $namespace
operator-sdk test local ./test/e2e/applicationv1alpha1 --up-local --namespace $namespace
exit_code=$?

kubectl delete namespace $namespace

exit $exit_code
