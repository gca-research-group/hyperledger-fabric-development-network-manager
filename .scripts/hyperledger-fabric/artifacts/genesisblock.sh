#!/bin/bash
source ./.scripts/config/_colors.sh
source ./.scripts/config/_icons.sh
source ./.scripts/hyperledger-fabric/_variables.sh
source ./.scripts/hyperledger-fabric/artifacts/_utils.sh

CONTAINER_NAME=hyperledger-fabric-tools

OUTPUT_BLOCK=$BASE_PATH/channel/$SYSTEM_CHANNEL_ID.block
PROFILE=MultiOrdererGenesis
CONFIG_TX_COMMAND="configtxgen -outputBlock $OUTPUT_BLOCK -profile $PROFILE -channelID $SYSTEM_CHANNEL_ID -configPath $BASE_PATH"

generateGenesisBlock() {
  echo -e "${PROCESSING_ICON} Generating genesis block."
  result=$(docker exec -it $1 bash -c "$2" 2>&1)
  echo -e "${SUCCESS_ICON} Genesis block generated."
}

verifyIfConfigTxFileExists $CONFIG_TX_FILE
verifyIfCryptoMaterialsExist $CRYPTO_CONFIG_FOLDER
generateGenesisBlock $CONTAINER_NAME "$CONFIG_TX_COMMAND"

echo -e "${SUCCESS_ICON} Finished succesfully."
exit 0