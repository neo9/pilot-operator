#!/usr/bin/env sh

if [ $TRAVIS_BRANCH != 'master' ]; then
  echo "Not on master branch, no e2e possible"
  exit 0
fi


namespace="test-e2e-$(git rev-parse --short HEAD)-$TRAVIS_BUILD_ID"

kubectl create namespace $namespace
operator-sdk test local ./test/e2e/applicationv1alpha1 --namespace $namespace
exit_code=$?

kubectl delete namespace $namespace

exit $exit_code
