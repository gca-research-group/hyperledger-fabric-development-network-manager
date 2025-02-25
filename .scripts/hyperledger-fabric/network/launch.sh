echo -e "========= Clean ========="
./.scripts/hyperledger-fabric/network/down.sh

echo -e "\n========= Up ========="
./.scripts/hyperledger-fabric/network/up.sh

echo -e "\n========= Channel ========="
./.scripts/hyperledger-fabric/network/channel.sh