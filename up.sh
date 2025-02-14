#!/bin/bash
source ./config/_colors.sh
source ./config/_icons.sh

DOCKER_COMPOSE_FILE=./.docker/hyperledger-fabric/hyperledger-fabric-network.yml

echo -e "${PROCESSING_ICON} Initializing the network."
docker compose -f $DOCKER_COMPOSE_FILE up --build -d > /dev/null 2>&1
echo -e "${SUCCESS_ICON} Network initialized."
exit 0
