#!/bin/bash
source ./_variables.sh
source ./config/_colors.sh
source ./config/_icons.sh

ORGS="org1 org2 org3"

CC_LABEL="$CHAINCODE.1.0-1.0"

BASE_PATH=/etc/hyperledger/fabric
CHAINCODE_PATH=$BASE_PATH/chaincode
CC_PACKAGE_FILE=$CC_LABEL.tar.gz

for org in $ORGS; do
    container="peer0.${org}.example.com"

    echo -e "${PROCESSING_ICON} Installing the chaincode: ${container}."
    echo -e "${PROCESSING_ICON} Verifying if the chaincode is already installed."

    CORE_PEER_MSPCONFIGPATH="$BASE_PATH/crypto-config/users/Admin@${org}.example.com/msp"
    COMMAND="export CORE_PEER_MSPCONFIGPATH="$CORE_PEER_MSPCONFIGPATH" && peer lifecycle chaincode queryinstalled | grep "$CHAINCODE""
    result=$(docker exec -it $container bash -c "$COMMAND" 2>&1)

    if [[ "$result" == *"Error"* ]]; then
        echo -e "${FAIL_ICON} Failed to intall the chaincode: ${RED}$result${NO_COLOR}"
    elif [[ -n "$result" ]]; then
        echo -e "${SUCCESS_ICON} The chaincode is already installed." 
        echo -e "${SUCCESS_ICON} Finished succesfully: $result."  
    else
        echo -e "${PROCESSING_ICON} Installing."
        COMMAND="export CORE_PEER_MSPCONFIGPATH="$CORE_PEER_MSPCONFIGPATH" && cd $CHAINCODE_PATH && peer lifecycle chaincode install $CC_PACKAGE_FILE"
        docker exec -it $container bash -c "$COMMAND"
        echo -e "${SUCCESS_ICON} Installed."
    fi
done

echo -e "${SUCCESS_ICON} Finished succesfully."
exit 0
