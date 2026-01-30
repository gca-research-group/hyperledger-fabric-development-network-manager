CHAINCODE_NAME=test

CHANNEL_ID=orgschannel
SYSTEM_CHANNEL_ID=systemchannel

BASE_PATH=/opt/gopath/src/github.com/hyperledger/fabric
CRYPTO_CONFIG_FOLDER=./.docker/hyperledger-fabric/artifacts/crypto-materials
CRYPTO_CONFIG_FILE=./.docker/hyperledger-fabric/artifacts/crypto-config.yml
CONFIG_TX_FILE=./.docker/hyperledger-fabric/artifacts/configtx.yml

HYPERLEDGER_FABRIC_NETWORK=./.docker/hyperledger-fabric/network.yml
HYPERLEDGER_FABRIC_TOOLS=./.docker/hyperledger-fabric/tools.yml

PEER0_ORG1=./.docker/hyperledger-fabric/orgs/org1/peer0.yml
PEER1_ORG1=./.docker/hyperledger-fabric/orgs/org1/peer1.yml
CA_ORG1=./.docker/hyperledger-fabric/orgs/org1/ca.yml

PEER0_ORG2=./.docker/hyperledger-fabric/orgs/org2/peer0.yml
PEER1_ORG2=./.docker/hyperledger-fabric/orgs/org2/peer1.yml
CA_ORG2=./.docker/hyperledger-fabric/orgs/org2/ca.yml

PEER0_ORG3=./.docker/hyperledger-fabric/orgs/org3/peer0.yml
PEER1_ORG3=./.docker/hyperledger-fabric/orgs/org3/peer1.yml
CA_ORG3=./.docker/hyperledger-fabric/orgs/org3/ca.yml

ORDERER=./.docker/hyperledger-fabric/orderer.yml
ORDERER_HOST=orderer.example.com:7050
ORDERER_CA=$BASE_PATH/crypto-materials/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem

