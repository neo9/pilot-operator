#!/usr/bin/env sh

if [ $TRAVIS_BRANCH != 'master' ]; then
  echo "Not on master branch, no e2e possible"
  exit 0
fi

namespace="test-2e2-$(git rev-parse --short HEAD)"

kubectl create namespace $namespace
operator-sdk test local ./test/e2e --namespace $namespace
exit_code=$?

kubectl delete namespace $namespace

exit $exit_code
