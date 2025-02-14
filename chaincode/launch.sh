echo -e "========= Clean ========="
./chaincode/clean.sh

echo -e "\n========= Package ========="
./chaincode/package.sh

echo -e "\n========= Install ========="
./chaincode/install.sh

echo -e "\n========= Approve ========="
./chaincode/approve.sh

echo -e "\n========= Commit ========="
./chaincode/commit.sh