#!/bin/sh

set -e

APPNAME=$(basename $(PWD))
TODAY=$(date -u +"%Y-%m-%d")

DIR_BIN=pkg/bin
DIR_ZIP=pkg/zip

rm -rf pkg/*
mkdir -p ${DIR_BIN}/${TODAY}
mkdir -p ${DIR_ZIP}/${TODAY}

gox \
  -output "${DIR_BIN}/${TODAY}/{{.OS}}_{{.Arch}}/${APPNAME}" \
  -osarch "linux/amd64" \
  -osarch "darwin/amd64"

for PLATFORM in $(find ${DIR_BIN}/${TODAY} -mindepth 1 -maxdepth 1 -type d); do
  OSARCH=$(basename ${PLATFORM})
  pushd $PLATFORM >/dev/null 2>&1
  zip ../../../../${DIR_ZIP}/${TODAY}/${APPNAME}_${TODAY}_${OSARCH}.zip ./*
  popd >/dev/null 2>&1
done

pushd ${DIR_ZIP}/${TODAY} >/dev/null 2>&1
pwd
shasum -a256 *.zip > ${APPNAME}_${TODAY}_SHA256SUMS
popd >/dev/null 2>&1

pwd
aws s3 cp ${DIR_ZIP}/${TODAY}/ s3://${BIN_BUCKET_NAME}/projects/${APPNAME}/${TODAY}/ --recursive #--dryrun
