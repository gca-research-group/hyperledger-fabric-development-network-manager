#!/bin/bash
source ./_variables.sh
source ./config/_colors.sh
source ./config/_icons.sh

ORGS="org1 org2 org3"

BASE_PATH=/etc/hyperledger/fabric
VERSION=1.0
SEQUENCE=1

for org in $ORGS; do
    container="peer0.${org}.example.com"

    echo -e "${PROCESSING_ICON} Approving the chaincode: ${container}."

    CORE_PEER_MSPCONFIGPATH="$BASE_PATH/crypto-config/users/Admin@${org}.example.com/msp"
    COMMAND="export CORE_PEER_MSPCONFIGPATH="$CORE_PEER_MSPCONFIGPATH" && peer lifecycle chaincode queryinstalled"
    INSTALLED_CHAINCODES=$(docker exec -it $container bash -c "$COMMAND")
    PACKAGE_ID=$(echo "$INSTALLED_CHAINCODES" | grep -oP '(?<=Package ID:)[^,]+' | awk '{print $1}')

    COMMAND="export CORE_PEER_MSPCONFIGPATH="$CORE_PEER_MSPCONFIGPATH" && peer lifecycle chaincode approveformyorg -n $CHAINCODE_NAME -v $VERSION -C $CHANNEL_ID --sequence $SEQUENCE --package-id $PACKAGE_ID"
    result=$(docker exec -it $container bash -c "$COMMAND")

    if [[ "$result" == *"Error"* ]]; then
        echo -e "${FAIL_ICON} Failed to approve the chaincode: ${RED}$result ${NO_COLOR}"
    else
        echo -e "${SUCCESS_ICON} Approved."
    fi
done

echo -e "${SUCCESS_ICON} Finished succesfully."
exit 0