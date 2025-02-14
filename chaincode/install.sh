#!/bin/bash
source ./config/_colors.sh
source ./config/_icons.sh

CONTAINERS="peer0.org1.example.com peer0.org2.example.com peer0.org3.example.com"

CC_LABEL=asset.1.0-1.0

BASE_PATH=/etc/hyperledger/fabric
CHAINO_COLORODE_PATH=$BASE_PATH/chaincode
CC_PACKAGE_FILE=$CC_LABEL.tar.gz

for container in $CONTAINERS; do
    echo -e "${PROCESSING_ICON} Installing the chaincode: ${container}."
    echo -e "${PROCESSING_ICON} Verifying if the chaincode is already installed."
    COMMAND='peer lifecycle chaincode queryinstalled | grep "asset"'
    result=$(docker exec -it $container bash -c "$COMMAND" 2>&1)

    if [[ -n "$result" ]]; then
        echo -e "${SUCCESS_ICON} The chaincode is already installed." 
        echo -e "${SUCCESS_ICON} Finished succesfully."  
        exit 0
    fi

    echo -e "${PROCESSING_ICON} Installing."
    COMMAND="cd $CHAINO_COLORODE_PATH && peer lifecycle chaincode install $CC_PACKAGE_FILE"
    docker exec -it $container bash -c "$COMMAND"
    echo -e "${SUCCESS_ICON} Installed."
done

echo -e "${SUCCESS_ICON} Finished succesfully."
exit 0
