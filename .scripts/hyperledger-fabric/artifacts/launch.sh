#!/bin/bash
echo -e "========= Down ========="
./.scripts/hyperledger-fabric/artifacts/down.sh

echo -e "========= Clean ========="
./.scripts/hyperledger-fabric/artifacts/clean.sh

echo -e "========= Up ========="
./.scripts/hyperledger-fabric/artifacts/up.sh

echo -e "\n========= Crypto Materials ========="
./.scripts/hyperledger-fabric/artifacts/cryptomaterials.sh

echo -e "\n========= Genesis Block ========="
./.scripts/hyperledger-fabric/artifacts/genesisblock.sh

echo -e "\n========= Channel ========="
./.scripts/hyperledger-fabric/artifacts/channeltx.sh
