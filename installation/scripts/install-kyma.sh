#!/usr/bin/env bash

set -o errexit

echo "Installing Kyma..."

LOCAL_ENV=${LOCAL_ENV:-false}

CURRENT_DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )
SCRIPTS_DIR="${CURRENT_DIR}/../scripts"
source $SCRIPTS_DIR/utils.sh
useMinikube

POSITIONAL=()
while [[ $# -gt 0 ]]
do
    key="$1"

    case ${key} in
        --kyma-release)
            checkInputParameterValue "${2}"
            KYMA_RELEASE="$2"
            shift
            shift
            ;;
         --kyma-installation)
            checkInputParameterValue "${2}"
            KYMA_INSTALLATION="$2"
            shift
            shift
            ;;
        --*)
            echo "Unknown flag ${1}"
            exit 1
        ;;
        *) # unknown option
            POSITIONAL+=("$1") # save it in an array for later
            shift # past argument
            ;;
    esac
done
set -- "${POSITIONAL[@]}" # restore positional parameters

ROOT_PATH=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )/../..

CERT=$(<$HOME/.minikube/ca.crt)
CERT="${CERT//$'\n'/\\\\n}"

INSTALLER_CR_PATH="${ROOT_PATH}"/installation/resources/kyma/installer-cr-kyma-minimal.yaml
OVERRIDES_KYMA_MINIMAL_CFG_LOCAL="${ROOT_PATH}"/installation/resources/kyma/installer-overrides-kyma-minimal-config-local.yaml

MINIMAL_OVERRIDES_FILENAME=override-local-minimal.yaml
MINIMAL_OVERRIDES_CONTENT=$(sed "s~\"__CERT__\"~\"$CERT\"~" "${OVERRIDES_KYMA_MINIMAL_CFG_LOCAL}")

>"${MINIMAL_OVERRIDES_FILENAME}" cat <<-EOF
$MINIMAL_OVERRIDES_CONTENT
EOF

INSTALLER_CR_FULL_PATH="${ROOT_PATH}"/installation/resources/kyma/installer-cr-kyma.yaml
OVERRIDES_KYMA_FULL_CFG_LOCAL="${ROOT_PATH}"/installation/resources/kyma/installer-overrides-kyma-full-config-local.yaml

FULL_OVERRIDES_FILENAME=override-local-full.yaml
FULL_OVERRIDES_CONTENT=$(sed "s~\"__CERT__\"~\"$CERT\"~" "${OVERRIDES_KYMA_FULL_CFG_LOCAL}")

>"${FULL_OVERRIDES_FILENAME}" cat <<-EOF
$FULL_OVERRIDES_CONTENT
EOF

trap "rm -f ${MINIMAL_OVERRIDES_FILENAME} ${FULL_OVERRIDES_FILENAME}" EXIT INT TERM

if [[ $KYMA_RELEASE == *PR-* ]]; then
  KYMA_TAG=$(curl -L https://storage.googleapis.com/kyma-development-artifacts/${KYMA_RELEASE}/kyma-installer-cluster.yaml | grep 'image: eu.gcr.io/kyma-project/kyma-installer:'| sed 's+image: eu.gcr.io/kyma-project/kyma-installer:++g' | tr -d '[:space:]')
  if [ -z "$KYMA_TAG" ]; then echo "ERROR: Kyma artifacts for ${KYMA_RELEASE} not found."; exit 1; fi
  KYMA_SOURCE="eu.gcr.io/kyma-project/kyma-installer:${KYMA_TAG}"
elif [[ $KYMA_RELEASE == main ]]; then
  KYMA_SOURCE="main"
elif [[ $KYMA_RELEASE == *main-* ]]; then
  KYMA_SOURCE=$(echo $KYMA_RELEASE | sed 's+main-++g' | tr -d '[:space:]')
else
  KYMA_SOURCE="${KYMA_RELEASE}"
fi

echo "Using Kyma source '${KYMA_SOURCE}'..."

echo "Installing Kyma..."
set -o xtrace
if [[ $KYMA_INSTALLATION == *full* ]]; then
  echo "Installing full Kyma"
  kyma install -c $INSTALLER_CR_FULL_PATH -o $FULL_OVERRIDES_FILENAME --source $KYMA_SOURCE
else
  echo "Installing minimal Kyma"
  kyma install -c $INSTALLER_CR_PATH -o $MINIMAL_OVERRIDES_FILENAME --source $KYMA_SOURCE
fi
set +o xtrace

# Kyma CLI uses the internal IP for the /etc/hosts override when docker driver is used. However on the host machine this should be localhost
USED_DRIVER=$(minikube profile list -o json | jq -r ".valid[0].Config.Driver")
if [[ $USED_DRIVER == "docker" ]]; then
  MINIKUBE_IP=$(minikube ssh egrep "minikube$" /etc/hosts | cut -f1)
  if [ "$(uname)" == "Darwin" ]; then #  this is the case when the script is ran on local Mac OSX machines, reference issue: https://stackoverflow.com/questions/4247068/sed-command-with-i-option-failing-on-mac-but-works-on-linux
    sudo sed -i "" "s/$MINIKUBE_IP/127.0.0.1/g" /etc/hosts
  else # this is the case when the script is ran on non-Mac OSX machines, ex. as part of remote PR jobs
    sudo sed -i "s/$MINIKUBE_IP/127.0.0.1/g" /etc/hosts
  fi
fi