#!/bin/bash
source ./.scripts/config/_colors.sh
source ./.scripts/config/_icons.sh
source ./.scripts/hyperledger-fabric/_utils.sh
source ./.scripts/hyperledger-fabric/_variables.sh

verifyIfDockerIfRunning

removeContainersInExecution $HYPERLEDGER_FABRIC_NETWORK $PEER0_ORG1 $CA_ORG1
removeContainersInExecution $HYPERLEDGER_FABRIC_NETWORK $PEER1_ORG1 $CA_ORG1

removeContainersInExecution $HYPERLEDGER_FABRIC_NETWORK $PEER0_ORG2 $CA_ORG2
removeContainersInExecution $HYPERLEDGER_FABRIC_NETWORK $PEER1_ORG2 $CA_ORG2

removeContainersInExecution $HYPERLEDGER_FABRIC_NETWORK $PEER0_ORG3 $CA_ORG3
removeContainersInExecution $HYPERLEDGER_FABRIC_NETWORK $PEER1_ORG3 $CA_ORG3

removeContainersInExecution $HYPERLEDGER_FABRIC_NETWORK $ORDERER

echo -e "${SUCCESS_ICON} Finished succesfully."
exit 0
