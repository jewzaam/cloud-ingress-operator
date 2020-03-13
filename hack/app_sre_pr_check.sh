#!/bin/bash

# AppSRE team CD

set -exv

CURRENT_DIR=$(dirname "$0")

python "$CURRENT_DIR"/validate_yaml.py "$CURRENT_DIR"/../deploy/crds
if [ "$?" != "0" ]; then
    exit 1
fi

BASE_IMG="cloud-ingress-operator"
IMG="${BASE_IMG}:latest"

# verify the template has all parameters
NOT_FOUND=""
TEMPLATE_FILENAME="hack/olm-registry/olm-artifacts-template.yaml"
PARAMETERS=$(git grep "\${\([^}]*\)}" $TEMPLATE_FILENAME | tr '{' "\n" | tr '}' "\n" | grep -s -A1 "[$]" | grep -v -e "[$]" -e "^[-]" | sort -u)
for PARAM in $PARAMETERS; 
do 
    FOUND=$(cat $TEMPLATE_FILENAME | python -c 'import json, sys, yaml ; y=yaml.safe_load(sys.stdin.read()) ; json.dump(y, sys.stdout)' | jq -r ".parameters[] | select(.name == \"$PARAM\") | .name")
    if [ "$FOUND" == "" ];
    then 
        NOT_FOUND="$NOT_FOUND $PARAM"
    fi
done

if [ "$NOT_FOUND" != "" ];
then
    echo "FAILURE: the following parameters were not found in '$TEMPLATE_FILENAME': $NOT_FOUND"
    exit -1
fi

IMG="$IMG" make container-build
