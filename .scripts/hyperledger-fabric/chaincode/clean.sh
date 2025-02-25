#!/bin/bash
source ./.scripts/config/_colors.sh
source ./.scripts/config/_icons.sh

echo -e "${PROCESSING_ICON} Removing chaincode files."
rm -rf ./.docker/hyperledger-fabric/chaincode/go.sum
rm -rf ./.docker/hyperledger-fabric/chaincode/*.tar.gz

echo -e "${SUCCESS_ICON} Finished succesfully."
exit 0
