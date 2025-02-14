echo -e "========= Clean ========="
./chaincode/clean.sh

echo -e "\n========= Package ========="
./chaincode/package.sh

echo -e "\n========= Install ========="
./chaincode/install.sh
