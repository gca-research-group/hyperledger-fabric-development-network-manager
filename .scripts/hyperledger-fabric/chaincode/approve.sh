#!/bin/bash
source ./.scripts/config/_colors.sh
source ./.scripts/config/_icons.sh
source ./.scripts/hyperledger-fabric/_variables.sh

ORGS="Org1 Org2 Org3"
PEERS="peer0"

VERSION=1.0
SEQUENCE=1

CHAINCODE_PATH=$BASE_PATH/chaincode
CONTAINER=hyperledger-fabric-tools

for ORG in $ORGS; do
  for PEER in $PEERS; do
    echo -e "${PROCESSING_ICON} Approving the chaincode on the ${ORG}, ${PEER}."

    CORE_PEER_MSPCONFIGPATH="$BASE_PATH/crypto-materials/peerOrganizations/${ORG,,}.example.com/users/Admin@${ORG,,}.example.com/msp"
    CORE_PEER_ADDRESS="${PEER}.${ORG,,}.example.com:7051"
    CORE_PEER_LOCALMSPID="${ORG}MSP"
    CORE_PEER_TLS_ROOTCERT_FILE="$BASE_PATH/crypto-materials/peerOrganizations/${ORG,,}.example.com/peers/${PEER}.${ORG,,}.example.com/tls/ca.crt"
    
    COMMAND="CORE_PEER_MSPCONFIGPATH=$CORE_PEER_MSPCONFIGPATH CORE_PEER_ADDRESS=$CORE_PEER_ADDRESS CORE_PEER_LOCALMSPID=$CORE_PEER_LOCALMSPID CORE_PEER_TLS_ROOTCERT_FILE=$CORE_PEER_TLS_ROOTCERT_FILE peer lifecycle chaincode queryinstalled"
    INSTALLED_CHAINCODES=$(docker exec -it $CONTAINER bash -c "$COMMAND")
    PACKAGE_ID=$(echo "$INSTALLED_CHAINCODES" | grep -oP '(?<=Package ID:)[^,]+' | awk '{print $1}')
    
    COMMAND="cd $CHAINCODE_PATH && CORE_PEER_MSPCONFIGPATH=$CORE_PEER_MSPCONFIGPATH CORE_PEER_ADDRESS=$CORE_PEER_ADDRESS CORE_PEER_LOCALMSPID=$CORE_PEER_LOCALMSPID CORE_PEER_TLS_ROOTCERT_FILE=$CORE_PEER_TLS_ROOTCERT_FILE peer lifecycle chaincode approveformyorg -n $CHAINCODE_NAME -v $VERSION -C $CHANNEL_ID --sequence $SEQUENCE --package-id $PACKAGE_ID --tls --cafile $ORDERER_CA"
    result=$(docker exec $CONTAINER bash -c "$COMMAND")
  done
done

echo -e "${SUCCESS_ICON} Finished succesfully."
exit 0