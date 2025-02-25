#!/bin/bash
source ./.scripts/config/_colors.sh
source ./.scripts/config/_icons.sh
source ./.scripts/hyperledger-fabric/_utils.sh
source ./.scripts/hyperledger-fabric/_variables.sh

echo -e "${PROCESSING_ICON} Initializing Orderer."
runContainer $HYPERLEDGER_FABRIC_NETWORK $ORDERER

echo -e "${PROCESSING_ICON} Initializing Org1."
runContainer $HYPERLEDGER_FABRIC_NETWORK $PEER_ORG1 $CA_ORG1

echo -e "${PROCESSING_ICON} Initializing Org2."
runContainer $HYPERLEDGER_FABRIC_NETWORK $PEER_ORG2 $CA_ORG2

echo -e "${PROCESSING_ICON} Initializing Org3."
runContainer $HYPERLEDGER_FABRIC_NETWORK $PEER_ORG3 $CA_ORG3
echo -e "${SUCCESS_ICON} Network initialized."
exit 0
