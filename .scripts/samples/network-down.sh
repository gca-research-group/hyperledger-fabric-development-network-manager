#!/bin/bash

configs=(
  "capability-v3_0"
  "cross-version-compatibility"
  "minimal-network"
  "multi-channel-network-multi-consensus"
  "multi-channel-network"
  "multi-organisation-consortium"
  "network-with-chaincode"
)

for config in "${configs[@]}"; do
  echo "******* ${config} *******"
  go run cmd/cli/main.go network down --config samples/${config}.yml
done