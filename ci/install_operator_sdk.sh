#!/usr/bin/env sh

OPERATOR_SDK_VERSION=v0.8.1

set -e

curl -OJL https://github.com/operator-framework/operator-sdk/releases/download/${OPERATOR_SDK_VERSION}/operator-sdk-${OPERATOR_SDK_VERSION}-x86_64-linux-gnu

curl -OJL https://github.com/operator-framework/operator-sdk/releases/download/${OPERATOR_SDK_VERSION}/operator-sdk-${OPERATOR_SDK_VERSION}-x86_64-linux-gnu.asc

gpg --keyserver keyserver.ubuntu.com --recv-key "0CF50BEE7E4DF6445E08C0EA9AFDE59E90D2B445"

gpg --verify operator-sdk-${OPERATOR_SDK_VERSION}-x86_64-linux-gnu.asc

chmod +x operator-sdk-${OPERATOR_SDK_VERSION}-x86_64-linux-gnu
sudo cp operator-sdk-${OPERATOR_SDK_VERSION}-x86_64-linux-gnu /usr/local/bin/operator-sdk

rm operator-sdk-${OPERATOR_SDK_VERSION}-x86_64-linux-gnu

operator-sdk version
