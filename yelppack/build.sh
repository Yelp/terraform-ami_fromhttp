#!/bin/bash

project=$1
version=$2
iteration=$3

go get ${project}
pushd /go/src/github.com/hashicorp/terraform
git remote add bobtfish https://github.com/bobtfish/terraform.git
git fetch bobtfish
git fetch --tags bobtfish
git checkout 0.4.0-pre1
popd
pushd /go/bin
go build ../src/terraform-provider-yelpaws
popd
mkdir /dist && cd /dist
fpm -s dir -t deb --name ${project} \
    --iteration ${iteration} --version ${version} \
    /go/bin/${project}=/usr/bin/
