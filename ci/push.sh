#!/usr/bin/env sh

if [ $TRAVIS_BRANCH != 'master' ]; then
  echo "Not on master branch, no push"
  exit 0
fi

echo "Master branch detected. Pushing latest image"


echo "$DOCKER_PASSWORD" | docker login -u "$DOCKER_LOGIN" --password-stdin
docker push neo9sas/pilot-operator:latest
