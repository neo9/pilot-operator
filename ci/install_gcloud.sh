#!/usr/bin/env sh

set -e

if [ ! -d "$HOME/google-cloud-sdk/bin" ]; then
  rm -rf $HOME/google-cloud-sdk
  export CLOUDSDK_CORE_DISABLE_PROMPTS=1

  curl https://sdk.cloud.google.com | bash
fi
