#!/usr/bin/env bash

BASE_PATH=$PWD

# versions
CMP_ALERTS="1.0."
CMP_CONFIG="1.0."
CMP_DASHBOARD="1.0."
CMP_DATABASE="1.0."
CMP_GRPC="1.0."
CMP_LOGGER="1.0."
CMP_MAIL="1.0."
CMP_METRICS="1.1."
CMP_PROFILING="1.0."
CMP_WORKERS="1.1."

for CMP in `find components -maxdepth 1 -type d ! -path components`
do
    cd $BASE_PATH/$CMP

    CMP_NAME=`basename $PWD | tr "[:lower:]" "[:upper:]"`
    CMP_VERSION="1.0."

    CMP_VAR="CMP_${CMP_NAME}"
    if [ -n "${!CMP_VAR}" ]; then
        CMP_VERSION=${!CMP_VAR}
    fi

    CMP_BUILD_NUMBER=`git log component.go | wc -l`
    CMP_VERSION="${CMP_VERSION}${CMP_BUILD_NUMBER}"
    CMP_PACKAGE=`go list -e -f '{{.Name}}' ./`

cat << EOF > version.go
package ${CMP_PACKAGE}

const (
	ComponentVersion = "${CMP_VERSION}"
)
EOF
done

# formatting
cd $BASE_PATH
goimports -w ./

BINDATA="go-bindata-assetfs"
# BINDATA="go-bindata-assetfs -debug"

cd $BASE_PATH/components/alerts && $BINDATA -pkg=alerts templates/...
cd $BASE_PATH/components/dashboard && $BINDATA -pkg=dashboard templates/... public/...
cd $BASE_PATH/components/grpc && $BINDATA -pkg=grpc templates/... && protoc --go_out=plugins=grpc:. *.proto
cd $BASE_PATH/components/mail && $BINDATA -pkg=mail templates/...
cd $BASE_PATH/components/metrics && $BINDATA -pkg=metrics templates/...
cd $BASE_PATH/components/profiling && $BINDATA -pkg=profiling templates/... public/...
cd $BASE_PATH/components/workers && $BINDATA -pkg=workers templates/... public/...

cd $BASE_PATH
easyjson components/alerts/handler_ajax.go
easyjson components/workers/handler_ajax.go