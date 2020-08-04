#!/usr/bin/env bash

### 参数

### 生成证书

## 通道名称
CHANNEL_NAME="swanbychannel"


cryptogen generate --config=./crypto-config.yaml

sleep 5

## 生成创世块配置文件
configtxgen -profile SWANMultiNodeEtcdRaft   -channelID test-system-channel-name  -outputBlock  ./config/genesis.block
# configtxgen -profile SWANMultiNodeEtcdRaft     -outputBlock  ./config/genesis.block
if [ "$?" -ne 0 ]; then
  echo "Failed to generate orderer genesis block..."
  exit 1
fi


## 生成通道配置文件
configtxgen -profile TwoOrgsChannel -outputCreateChannelTx ./config/$CHANNEL_NAME.tx -channelID $CHANNEL_NAME
if [ "$?" -ne 0 ]; then
  echo "Failed to generate channel configuration transaction..."
  exit 1
fi

## 生成锚节点 org1
configtxgen -profile TwoOrgsChannel -outputAnchorPeersUpdate ./config/Org1MSPanchors.tx -channelID $CHANNEL_NAME -asOrg Org1MSP
if [ "$?" -ne 0 ]; then
  echo "Failed to generate anchor peer update for ProductOrg..."
  exit 1
fi

## 生成锚节点 org2
configtxgen -profile TwoOrgsChannel -outputAnchorPeersUpdate ./config/Org2MSPanchors.tx -channelID $CHANNEL_NAME -asOrg Org2MSP
if [ "$?" -ne 0 ]; then
  echo "Failed to generate anchor peer update for FactoringOrg..."
  exit 1
fi
