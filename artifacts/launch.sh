#!/bin/bash
echo -e "========= Clean ========="
./artifacts/clean.sh

echo -e "\n========= Crypto Materials ========="
./artifacts/cryptomaterials.sh

echo -e "\n========= Genesis Block ========="
./artifacts/genesisblock.sh

echo -e "\n========= Channel ========="
./artifacts/channel.sh
