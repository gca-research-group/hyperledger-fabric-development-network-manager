#!/bin/bash
source ./.scripts/config/_colors.sh
source ./.scripts/config/_icons.sh

BASE_FOLDER=./.docker/hyperledger-fabric/artifacts

echo -e "${PROCESSING_ICON} Removing crypto materials."
rm -rf $BASE_FOLDER/crypto-materials

echo -e "${PROCESSING_ICON} Removing channel artifacts."
rm -rf $BASE_FOLDER/channel

echo -e "${SUCCESS_ICON} Finished succesfully."
exit 0