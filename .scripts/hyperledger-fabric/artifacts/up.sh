source ./.scripts/hyperledger-fabric/_variables.sh
source ./.scripts/hyperledger-fabric/_utils.sh

CONTAINER_NAME=hyperledger-fabric-tools

verifyIfDockerComposeFileExists $HYPERLEDGER_FABRIC_TOOLS
verifyIfDockerComposeFileExists $HYPERLEDGER_FABRIC_NETWORK
runContainer $HYPERLEDGER_FABRIC_NETWORK $HYPERLEDGER_FABRIC_TOOLS

echo -e "${SUCCESS_ICON} Finished succesfully."
exit 0
