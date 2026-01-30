#!/bin/bash
source ./.scripts/config/_colors.sh
source ./.scripts/config/_icons.sh
source ./.scripts/hyperledger-fabric/_variables.sh

ORGS="Org1 Org2 Org3"
PEERS="peer0 peer1"

CC_LABEL="$CHAINCODE_NAME.1.0-1.0"

CHAINCODE_PATH=$BASE_PATH/chaincode
CC_PACKAGE_FILE=$CC_LABEL.tar.gz
CONTAINER=hyperledger-fabric-tools

for ORG in $ORGS; do
    for PEER in $PEERS; do
      echo -e "${PROCESSING_ICON} Installing the chaincode on the ${ORG}, ${PEER}."

      CORE_PEER_MSPCONFIGPATH="$BASE_PATH/crypto-materials/peerOrganizations/${ORG,,}.example.com/users/Admin@${ORG,,}.example.com/msp"
      CORE_PEER_ADDRESS="${PEER}.${ORG,,}.example.com:7051"
      CORE_PEER_LOCALMSPID="${ORG}MSP"
      CORE_PEER_TLS_ROOTCERT_FILE="$BASE_PATH/crypto-materials/peerOrganizations/${ORG,,}.example.com/peers/${PEER}.${ORG,,}.example.com/tls/ca.crt"

      COMMAND="cd $CHAINCODE_PATH && CORE_PEER_MSPCONFIGPATH=$CORE_PEER_MSPCONFIGPATH CORE_PEER_ADDRESS=$CORE_PEER_ADDRESS CORE_PEER_LOCALMSPID=$CORE_PEER_LOCALMSPID CORE_PEER_TLS_ROOTCERT_FILE=$CORE_PEER_TLS_ROOTCERT_FILE peer lifecycle chaincode install $CC_PACKAGE_FILE"
      result=$(docker exec $CONTAINER bash -c "$COMMAND")
    done
done

echo -e "${SUCCESS_ICON} Finished succesfully."
exit 0
