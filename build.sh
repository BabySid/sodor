#!/bin/bash

set -e

FAIL=1
OK=0

WORKSPACE=$(cd $(dirname $0) && pwd -P)

MIN_VERSION="1144"
function check_go_version() {
    local version=$(go version | awk '{print $3}' | sed s/[go[:space:].]//g)
    if [[ "${version}" -le "${MIN_VERSION}" ]]; then
      echo "version not matched. min version expect: " ${MIN_VERSION} " actual: " ${version}
      return ${FAIL}
    fi
    return ${OK}
}


function reset_env() {
    export GOPROXY="https://goproxy.io/,direct"
    export GO111MODULE=auto
}

BuildTime=$(date "+%F %T")
GoVersion=$(go version)

function begin() {
    rm -rf output
    mkdir -p output/bin
    mkdir -p output/logs
}

function finish() {
    chmod +x output/bin/*

    cd output
    tar -zcf sodor.tar.gz *
    cd ..
}

function build_fat_ctrl() {
    AppName="fat_ctrl"
    echo "########## build ${AppName} ##########"

    AppVersion="${AppName}"_$(date "+%F %T" | awk '{print $1"_"$2}')
    go build -ldflags "-X 'main.AppVersion=${AppVersion}'" -o ${AppName} ./fat_controller

    if [[ $? != 0 ]];then
        echo "compile ${AppName} failed"
        return ${FAIL}
    fi

    mv ${AppName} output/bin
}

function build_thomas() {
    AppName="thomas"
    echo "########## build ${AppName} ##########"

    AppVersion="${AppName}"_$(date "+%F %T" | awk '{print $1"_"$2}')
    go build -ldflags "-X 'main.AppVersion=${AppVersion}'" -o ${AppName} ./thomas

    if [[ $? != 0 ]];then
        echo "compile ${AppName} failed"
        return ${FAIL}
    fi

    mv ${AppName} output/bin
}

function main() {
    cd ${WORKSPACE}

    check_go_version
    if [[ $? -ne ${OK} ]]; then
        return ${FAIL}
    fi

    reset_env

    begin
    build_fat_ctrl
    if [[ $? -ne ${OK} ]]; then
        return ${FAIL}
    fi
    build_thomas
    if [[ $? -ne ${OK} ]]; then
        return ${FAIL}
    fi
    finish
    if [[ $? -ne ${OK} ]]; then
        return ${FAIL}
    fi

    return $?
}

main "$@"
