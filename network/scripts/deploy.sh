#!/bin/bash
set -e

echo "=== 1. Generating Crypto Materials ==="
cryptogen generate --config=crypto-config.yaml --output=crypto-config

echo "=== 2. Generating Genesis Block ==="
configtxgen -profile OrdererGenesis -channelID system-channel -outputBlock channel-artifacts/genesis.block

echo "=== 3. Generating Channel Configuration ==="
configtxgen -profile ChannelConfig -outputCreateChannelTx channel-artifacts/channel.tx -channelID mychannel

echo "=== 4. Starting Network ==="
docker-compose up -d

echo "=== 5. Creating Channel ==="
docker exec peer0.org1.example.com peer channel create -o orderer.example.com:7050 -c mychannel -f channel-artifacts/channel.tx

echo "=== 6. Joining Peers to Channel ==="
docker exec peer0.org1.example.com peer channel join -b mychannel.block
docker exec peer1.org1.example.com peer channel join -b mychannel.block

echo "=== 7. Installing Chaincode ==="
docker exec peer0.org1.example.com peer lifecycle chaincode install education.tar.gz

echo "=== 8. Deploying Chaincode ==="
docker exec peer0.org1.example.com peer lifecycle chaincode approveformyorg -n education -v 1.0 -C mychannel
docker exec peer1.org1.example.com peer lifecycle chaincode approveformyorg -n education -v 1.0 -C mychannel
docker exec peer0.org1.example.com peer lifecycle chaincode commit -n education -v 1.0 -C mychannel

echo "=== Done! Network is ready ==="
