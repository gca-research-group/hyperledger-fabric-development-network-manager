#!/bin/bash
source ./.scripts/config/_colors.sh
source ./.scripts/config/_icons.sh
source ./.scripts/hyperledger-fabric/_variables.sh

CC_LABEL="$CHAINCODE_NAME.1.0-1.0"

CHAINCODE_PATH=$BASE_PATH/chaincode
CC_PACKAGE_FILE=$CC_LABEL.tar.gz
CONTAINER=hyperledger-fabric-tools

echo -e "${PROCESSING_ICON} Installing dependencies."
COMMAND="cd $CHAINCODE_PATH && go mod tidy"
result=(docker exec -it $CONTAINER bash -c "$COMMAND")
echo -e "${SUCCESS_ICON} Dependencies installed."

COMMAND="cd $CHAINCODE_PATH && peer lifecycle chaincode package $CC_PACKAGE_FILE -p $CHAINCODE_NAME.go --label $CC_LABEL"
result=(docker exec -it $CONTAINER bash -c "$COMMAND")

echo -e "${SUCCESS_ICON} Finished succesfully."
exit 0
