#!/bin/bash
source ./.scripts/config/_colors.sh
source ./.scripts/config/_icons.sh
source ./.scripts/hyperledger-fabric/_variables.sh
source ./.scripts/hyperledger-fabric/_utils.sh

sleep 5

ORGS="Org1 Org2 Org3"
PEERS="peer0 peer1"

GENESIS_BLOCK=$BASE_PATH/channel/$CHANNEL_ID.block
CHANNEL_TX=./$CHANNEL_ID.tx
CONTAINER=hyperledger-fabric-tools

echo -e "${PROCESSING_ICON} Creating channel."
COMMAND="cd $BASE_PATH/channel && peer channel create -o $ORDERER_HOST -c $CHANNEL_ID -f $CHANNEL_TX --tls true --cafile $ORDERER_CA"
result=$(docker exec -it $CONTAINER bash -c "$COMMAND")
echo -e "${SUCCESS_ICON} Channel created."

sleep 5

for ORG in $ORGS; do
  for PEER in $PEERS; do
    echo -e "${PROCESSING_ICON} Joining peer to the channel: ${ORG} ${PEER}."   

    CORE_PEER_MSPCONFIGPATH="$BASE_PATH/crypto-materials/peerOrganizations/${ORG,,}.example.com/users/Admin@${ORG,,}.example.com/msp"
    CORE_PEER_ADDRESS="${PEER}.${ORG,,}.example.com:7051"
    CORE_PEER_LOCALMSPID="${ORG}MSP"
    CORE_PEER_TLS_ROOTCERT_FILE="$BASE_PATH/crypto-materials/peerOrganizations/${ORG,,}.example.com/peers/${PEER}.${ORG,,}.example.com/tls/ca.crt"

    COMMAND="CORE_PEER_MSPCONFIGPATH=$CORE_PEER_MSPCONFIGPATH CORE_PEER_ADDRESS=$CORE_PEER_ADDRESS CORE_PEER_LOCALMSPID=$CORE_PEER_LOCALMSPID CORE_PEER_TLS_ROOTCERT_FILE=$CORE_PEER_TLS_ROOTCERT_FILE peer channel join -b $GENESIS_BLOCK"
    result=$(docker exec -it $CONTAINER bash -c "$COMMAND")
  done
done

echo -e "${SUCCESS_ICON} Finished succesfully."
exit 0