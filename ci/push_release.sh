#!/usr/bin/env sh

if ! [ -z "$TRAVIS_TAG" ]; then
    echo "Release $TRAVIS_TAG"
    docker tag neo9sas/pilot-operator:latest neo9sas/pilot-operator:$TRAVIS_TAG
    docker push neo9sas/pilot-operator:$TRAVIS_TAG
fi
