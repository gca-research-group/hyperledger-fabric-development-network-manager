#!/bin/bash
source ./.scripts/config/_colors.sh
source ./.scripts/config/_icons.sh
source ./.scripts/hyperledger-fabric/artifacts/_utils.sh
source ./.scripts/hyperledger-fabric/_variables.sh

CONTAINER_NAME=hyperledger-fabric-tools

OUTPUT_CHANNEL=./channel/$CHANNEL_ID.tx
PROFILE=MultiChannel
CONFIG_TX_COMMAND="configtxgen -outputCreateChannelTx $OUTPUT_CHANNEL -profile $PROFILE -channelID $CHANNEL_ID -configPath $BASE_PATH"

generateChannel() {
  echo -e "${PROCESSING_ICON} Generating the channel."
  result=$(docker exec -it $1 bash -c "$2" 2>&1)
  echo -e "${SUCCESS_ICON} Channel generated."
}

verifyIfConfigTxFileExists $CONFIG_TX_FILE
verifyIfCryptoMaterialsExist $CRYPTO_CONFIG_FOLDER
generateChannel $CONTAINER_NAME "$CONFIG_TX_COMMAND"

echo -e "${SUCCESS_ICON} Finished succesfully."
exit 0