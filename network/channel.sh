#!/bin/bash
source ./_variables.sh
source ./config/_colors.sh
source ./config/_icons.sh
source ./_utils.sh

DOCKER_COMPOSE_FILE=./.docker/hyperledger-fabric/hyperledger-fabric-network.yml

ORGS="org1 org2 org3"

ORDERER_HOST=orderer.example.com:7050

BASE_PATH=/etc/hyperledger/fabric

GENESIS_BLOCK=$BASE_PATH/genesis.block

for org in $ORGS; do
    echo -e "${PROCESSING_ICON} Joining peer to the channel: ${org}."
    
    container="peer0.${org}.example.com"

    result=$(docker exec -it $container bash -c "peer channel list")

    if [[ "$result" == *"$CHANNEL_ID"* ]]; then
        echo -e "${SUCCESS_ICON} Peer has already joined the channel. No action needed."
    else
        CORE_PEER_MSPCONFIGPATH="$BASE_PATH/crypto-config/users/Admin@${org}.example.com/msp"

        COMMAND="export CORE_PEER_MSPCONFIGPATH="$CORE_PEER_MSPCONFIGPATH" && peer channel join -o $ORDERER_HOST -b $GENESIS_BLOCK"
        result=$(docker exec -it $container bash -c "$COMMAND")

        if [[ "$result" == *"Error"* ]]; then
            echo -e "${FAIL_ICON} Failed to join the channel: ${RED}$result${NO_COLOR}"
        else
            echo -e "${SUCCESS_ICON} Peer Joined."
        fi

        CORE_PEER_MSPCONFIGPATH=$BASE_PATH/msp
        COMMAND="export CORE_PEER_MSPCONFIGPATH="$CORE_PEER_MSPCONFIGPATH""
        result=$(docker exec -it $container bash -c "$COMMAND")
    fi
done

echo -e "${SUCCESS_ICON} Finished succesfully."
exit 0