echo -e "========= Clean ========="
./.scripts/hyperledger-fabric/chaincode/clean.sh

echo -e "\n========= Package ========="
./.scripts/hyperledger-fabric/chaincode/package.sh

echo -e "\n========= Install ========="
./.scripts/hyperledger-fabric/chaincode/install.sh

echo -e "\n========= Approve ========="
./.scripts/hyperledger-fabric/chaincode/approve.sh

echo -e "\n========= Commit ========="
./.scripts/hyperledger-fabric/chaincode/commit.sh