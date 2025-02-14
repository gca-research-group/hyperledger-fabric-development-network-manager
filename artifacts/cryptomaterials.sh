#!/bin/bash
source ./config/_colors.sh
source ./config/_icons.sh
source ./artifacts/_utils.sh

CRYPTO_CONFIG_FOLDER=./.docker/hyperledger-fabric/artifacts/crypto-config
CRYPTO_CONFIG_FILE=./.docker/hyperledger-fabric/artifacts/crypto-config.yml
GITIGNORE_FILE=$(dirname $CRYPTO_CONFIG_FOLDER)/.gitignore
DOCKER_COMPOSE_FILE=./.docker/hyperledger-fabric/hyperledger-fabric-tools.yml
CONTAINER_NAME=hyperledger-fabric-tools
CRYPTO_CONFIG_COMMAND="cryptogen generate --config=/etc/hyperledger/fabric/crypto-config.yml --output=/etc/hyperledger/fabric/crypto-config"

verifyIfTheCryptoConfigFileExists() {
  echo -e "${PROCESSING_ICON} Verifying if the config file exists."
  if [[ ! -f $1 ]]; then
    echo -e "${RED}The crypto-config file has not been found:${NO_COLOR} $1"
    exit 1
  fi
}

generateCryptoMaterials() {
  echo -e "${PROCESSING_ICON} Generating Crypto Materials."
  exec_command=$(docker exec -it $1 bash -c "$2" 2>&1)
  if [[ $? -ne 0 ]]; then
    echo -e "${RED}Failed to execute command in the container:${NO_COLOR}"
    echo "$exec_command"
    exit 1
  fi

  echo -e "${SUCCESS_ICON} Crypto materials generated."
}

verifyIfTheCryptoConfigFileExists $CRYPTO_CONFIG_FILE
verifyIfDockerComposeFileExists $DOCKER_COMPOSE_FILE
removeContainersInExecution $DOCKER_COMPOSE_FILE
runContainer $DOCKER_COMPOSE_FILE
verifyIfContainerIsRunning $CONTAINER_NAME
generateCryptoMaterials $CONTAINER_NAME "$CRYPTO_CONFIG_COMMAND"
removeContainersInExecution $DOCKER_COMPOSE_FILE

echo -e "${SUCCESS_ICON} Finished succesfully."
exit 0