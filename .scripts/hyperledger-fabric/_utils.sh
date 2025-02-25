#!/bin/bash
source ./.scripts/config/_colors.sh
source ./.scripts/config/_icons.sh

verifyIfDockerIfRunning() {
  echo -e "${PROCESSING_ICON} Verifying if Docker is running."
  if ! docker info > /dev/null 2>&1; then
    echo -e "${RED}Docker is not running. Please start Docker and try again.${NO_COLOR}"
    exit 1
  fi
  echo -e "${SUCCESS_ICON} Docker is running."
}

verifyIfContainerIsRunning() {
  echo -e "${PROCESSING_ICON} Verifying if the container is running."
  container_status=$(docker inspect -f '{{.State.Running}}' $1 2>&1)

  if [[ "$container_status" != "true" ]]; then
    echo -e "${RED}The container is not running:${NO_COLOR}"
    echo "$container_status"
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

  COMMAND=""

  for arg in "$@"; do
    COMMAND="$COMMAND -f $arg"
  done;

  errors=$(docker compose $COMMAND down 2>&1 > /dev/null)

  if [[ -n "$errors" && "$errors" != *"Removed"* ]]; then
    echo -e "${RED}Failed to remove the containers:"
    echo -e "$errors${NO_COLOR}"
    exit 1
  fi

  echo -e "${SUCCESS_ICON} Containers deleted."
}

runContainer() {
  echo -e "${PROCESSING_ICON} Starting the container."

  COMMAND=""

  for arg in "$@"; do
    COMMAND="$COMMAND -f $arg"
  done;

  errors=$(docker compose $COMMAND up --build -d 2>&1 > /dev/null)

  echo -e "${SUCCESS_ICON} Container started."
}