#!/bin/bash
source ./config/_colors.sh
source ./config/_icons.sh

verifyIfContainerIsRunning() {
  echo -e "${PROCESSING_ICON} Verifying if the container is running."
  container_status=$(docker inspect -f '{{.State.Running}}' $1 2>&1)

  if [[ "$container_status" != "true" ]]; then
    echo -e "${RED}The container is not running:${NO_COLOR}"
    echo "$container_status"
    exit 1
  fi
}

verifyIfConfigTxFileExists() {
  echo -e "${PROCESSING_ICON} Verifying if the config file exists."
  if [[ ! -f $1 ]]; then
    echo -e "${RED}The configtx file has not been found:${NO_COLOR} $1"
    exit 1
  fi
}

verifyIfDockerComposeFileExists() {
  echo -e "${PROCESSING_ICON} Verifying if the docker compose file exists."

  if [[ ! -f $1 ]]; then
    echo -e "${RED}The docker compose file has not been found:${NO_COLOR} $1"
    exit 1
  fi
}

removeContainersInExecution() {
  echo -e "${PROCESSING_ICON} Deleting containers in execution."
  errors=$(docker compose -f $1 down 2>&1 > /dev/null)

  if [[ -n "$errors" && "$errors" != *"Removed"* ]]; then
    echo -e "${RED}Failed to remove the containers:${NO_COLOR}"
    echo "$errors"
    exit 1
  fi

  echo -e "${SUCCESS_ICON} Containers deleted."
}

runContainer() {
  echo -e "${PROCESSING_ICON} Starting the container."
  errors=$(docker compose -f $1 up -d 2>&1 > /dev/null)
  echo -e "${SUCCESS_ICON} Container started."
}

verifyIfCryptoMaterialsExist() {
  echo -e "${PROCESSING_ICON} Verifying if the crypto materials exist."

  if [[ ! -d $1 ]]; then
    echo -e "${RED}Crypto materials do not exists:${NO_COLOR} $1"
    exit 1
  fi
}