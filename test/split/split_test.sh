#!/bin/bash

set -eu
cd $(dirname $0)
go build -i ../..

function concatenate() {
    local -r dir=$1
    local -r tempfile=$(mktemp)
    # Use tempfile to avoid 'too many arguments error' when using `cat *` directory
    for i in ${dir}/*; do
        cat $i >> ${tempfile}
    done
    cat ${tempfile}
}

function execute_tests() {
    local -r bulk_passed=$1
    local -ar split_patterns=(1 2 500 10000)
    for testcase in *; do
        if ! [ -d $testcase ]; then
            continue
        fi
        if [[ ${bulk_passed} > 0 ]] && [[ $testcase =~ bulk.* ]]; then
            continue
        fi
        local result_dir=./${testcase}/result
        for split_pattern in ${split_patterns[@]}; do
            mkdir -p ${result_dir}
            echo "[TEST] ${testcase} (${split_pattern})"
            ./psql-splitter split ./${testcase}/data/target.sql -o ${result_dir} -n ${split_pattern}
            concatenate ${result_dir} |diff ./${testcase}/data/target.sql -
            rm -r ${result_dir}
        done

        if [[ $? != 0 ]]; then
            echo '[FAILED] Difference found'
            exit 1
        fi
        echo '[PASSED]'
    done

    echo 'ALL TESTS PASSED'
    exit 0
}

execute_tests $1