source ./.scripts/hyperledger-fabric/_variables.sh
source ./.scripts/hyperledger-fabric/_utils.sh

removeContainersInExecution $HYPERLEDGER_FABRIC_NETWORK $HYPERLEDGER_FABRIC_TOOLS

echo -e "${SUCCESS_ICON} Finished succesfully."
exit 0