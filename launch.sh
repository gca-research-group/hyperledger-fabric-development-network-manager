#!/bin/bash
echo -e "********** Containers **********"
./network/down.sh

echo -e "\n********** Artifacts **********\n"
./artifacts/launch.sh

echo -e "\n********** Network **********\n"
./network/launch.sh

echo -e "\n********** Chaincode **********\n"
./chaincode/launch.sh
