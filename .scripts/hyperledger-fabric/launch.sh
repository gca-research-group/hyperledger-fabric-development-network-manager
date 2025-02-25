#!/bin/bash
echo -e "********** Containers **********"
./.scripts/hyperledger-fabric/network/down.sh

echo -e "\n********** Artifacts **********\n"
./.scripts/hyperledger-fabric/artifacts/launch.sh

echo -e "\n********** Network **********\n"
./.scripts/hyperledger-fabric/network/launch.sh

echo -e "\n********** Chaincode **********\n"
./.scripts/hyperledger-fabric/chaincode/launch.sh
