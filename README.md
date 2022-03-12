# rc-contract-go

This is the smart contract based on golang.

The smart contract can be tested with the following steps.

## Important

Everytime the contract is updated, run the following commands to update the local modules.

```go
go get -u
go mod tidy
go mod vendor
```

### Bring up the test network

Open a terminal from the /test-network dictionary, run

```shell
./network.sh down
./network.sh up createChannel -ca
```

Note that establishing gRPC connection requires the usage of CA in the latest release.

List the docker containers you created above

```shell
docker ps -a
```

The smart contract can be installed using the script provided by the fabric samples.

```shell
./network.sh deployCC -ccn basic -ccp ../../rc-contract-go/ -ccl go
```

Or it can be installed manually using the following commands.

Add the binaries and config filepath to the CLI path

```shell
export PATH=${PWD}/../bin:$PATH
export FABRIC_CFG_PATH=$PWD/../config/
```

Package chaincode (the chaincode path should be changed accordingly)

```shell
peer lifecycle chaincode package rc.tar.gz --path ../../rc-contract-go/ --lang golang --label rc_1.0
```

### Install and approve the chaincode as Org1

Operate the peer CLI as the Org1 admin user

```shell
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_LOCALMSPID="Org1MSP"
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
export CORE_PEER_ADDRESS=localhost:7051
```

From the /test-network terminal, run the following command to install chaincode to Org1 peer node

```shell
peer lifecycle chaincode install rc.tar.gz
```

If the chaincode intallation failed with github connection problem, open a terminal from /rc-contract-go and run

```go
go mod tidy
go mod vendor
```

Query the installed chaincode id (the package ID should be changed according to the ID returned by the command "queryinstalled")

```shell
peer lifecycle chaincode queryinstalled
```

Save the id to a variable (the chaincode id should be changed according to the output of the last command)

```shell
export CC_PACKAGE_ID=rc_1.0:ef3824ede27add5526cc37d47e6a8e7122a22129c570db19e98bf5f34b04ff47
```

Approve as Org1 admin

```shell
peer lifecycle chaincode approveformyorg -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --channelID mychannel --name basic --version 1.0 --package-id $CC_PACKAGE_ID --sequence 1 --tls --cafile ${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem
```

Check if the chaincode is ready to be committed to the channel (currently it should not be ready as the chaincode is not installed on Org2 peer and approved by Org2)

```shell
peer lifecycle chaincode checkcommitreadiness --channelID mychannel --name basic --version 1.0 --sequence 1 --tls --cafile ${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem --output json
```

Try to commit the chaincode to the channel (orderer service), it should fail here

```shell
peer lifecycle chaincode commit -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --channelID mychannel --name basic --version 1.0 --sequence 1 --tls --cafile ${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem --peerAddresses localhost:7051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt --peerAddresses localhost:9051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt
```

### Install and approve the chaincode as Org2

Operate the peer CLI as the Org2 admin user

```shell
export CORE_PEER_LOCALMSPID="Org2MSP"
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org2.example.com/users/Admin@org2.example.com/msp
export CORE_PEER_ADDRESS=localhost:9051
```

Install chaincode to Org2 peer node

```shell
peer lifecycle chaincode install rc.tar.gz
```

Approve chaincode definition as Org2 admin

```shell
peer lifecycle chaincode approveformyorg -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --channelID mychannel --name basic --version 1.0 --package-id $CC_PACKAGE_ID --sequence 1 --tls --cafile ${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem
```

Check if the chaincode is ready to be committed to the channel (it should be ready here as **BOTH** the Orgs has approved the chaincode definition)

```shell
peer lifecycle chaincode checkcommitreadiness --channelID mychannel --name basic --version 1.0 --sequence 1 --tls --cafile ${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem --output json
```

Committing the chaincode to the channel (orderer service)

```shell
peer lifecycle chaincode commit -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --channelID mychannel --name basic --version 1.0 --sequence 1 --tls --cafile ${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem --peerAddresses localhost:7051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt --peerAddresses localhost:9051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt
```

Check the commit result

```shell
peer lifecycle chaincode querycommitted --channelID mychannel --name basic --cafile ${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem
```

### Invoke chaincode

Invoking the method instantiate, which does nothing

```shell
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile ${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C mychannel -n basic --peerAddresses localhost:7051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt --peerAddresses localhost:9051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt -c '{"function":"Instantiate","Args":[]}'
```

Issue a new response credit,

```shell
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile ${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C mychannel -n basic --peerAddresses localhost:7051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt --peerAddresses localhost:9051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt -c '{"function":"Issue","Args":["000001","VPPO","2022-02-07"]}'
```

Querying the world state for a credit with credit number and issuer will

```shell
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile ${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C mychannel -n basic --peerAddresses localhost:7051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt --peerAddresses localhost:9051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt -c '{"function":"Query","Args":["000001","VPPO"]}'
```

Waiting for future updates...

### Clean up

```shell
./network.sh down
```
