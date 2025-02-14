#!/bin/bash
source ./config/_colors.sh
source ./config/_icons.sh

BASE_FOLDER=./.docker/hyperledger-fabric/artifacts

echo -e "${PROCESSING_ICON} Removing crypto materials."
rm -rf $BASE_FOLDER/crypto-config

echo -e "${PROCESSING_ICON} Removing .tx artifacts."
rm -rf $BASE_FOLDER/*.tx

echo -e "${PROCESSING_ICON} Removing genesis block."
rm -rf $BASE_FOLDER/*.block

echo -e "${SUCCESS_ICON} Finished succesfully."
exit 0