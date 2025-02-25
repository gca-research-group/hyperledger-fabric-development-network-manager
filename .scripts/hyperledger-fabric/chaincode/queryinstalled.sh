CORE_PEER_MSPCONFIGPATH="$BASE_PATH/crypto-materials/users/Admin@org1.example.com/msp"

export CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/crypto-materials/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp 
export CORE_PEER_ADDRESS=peer0.org1.example.com:7051
export CORE_PEER_LOCALMSPID=Org1MSP
peer lifecycle chaincode queryinstalled