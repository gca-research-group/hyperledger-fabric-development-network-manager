#!/bin/bash
source ./config/_colors.sh
source ./config/_icons.sh

CONTAINERS="peer0.org1.example.com peer0.org2.example.com peer0.org3.example.com"

CC_LABEL=asset.1.0-1.0

BASE_PATH=/etc/hyperledger/fabric
CHAINCODE_PATH=$BASE_PATH/chaincode
CC_PACKAGE_FILE=$CC_LABEL.tar.gz
SMART_CONTRACT=asset.go

for container in $CONTAINERS; do
    echo -e "${PROCESSING_ICON} Packaging the chaincode: ${container}."

    echo -e "${PROCESSING_ICON} Installing dependencies."
    COMMAND="cd $CHAINCODE_PATH && go mod tidy"
    docker exec -it $container bash -c "$COMMAND"
    echo -e "${SUCCESS_ICON} Dependencies installed."

    echo -e "${PROCESSING_ICON} Packaging."
    COMMAND="cd $CHAINCODE_PATH && peer lifecycle chaincode package $CC_PACKAGE_FILE -p $SMART_CONTRACT --label $CC_LABEL"

    result=$(docker exec -it $container bash -c "$COMMAND")

    if [[ "$result" ==  *"Error"* ]]; then
        echo -e "${FAIL_ICON} Failed to package the chaincode: ${RED}$result${NO_COLOR}"
    else
        echo -e "${SUCCESS_ICON} Packaged."
    fi
done;

echo -e "${SUCCESS_ICON} Finished succesfully."
exit 0
