#!/usr/bin/env sh

namespace="test-2e2-$(git rev-parse --short HEAD)-${TRAVIS_BRANCH}${TRAVIS_TAG}"

kubectl create namespace $namespace
operator-sdk test local ./test/e2e --namespace $namespace
exit_code=$?

kubectl delete namespace $namespace

exit $exit_code
