#!/bin/bash
echo -e "********** Containers **********"
./down.sh

echo -e "\n********** Artifacts **********\n"
./artifacts/launch.sh

echo -e "\n********** Network **********\n"
./up.sh
./channel.sh

echo -e "\n********** Chaincode **********\n"
./chaincode/launch.sh
