#!/bin/bash

configs=(
  "01-empty-output"
  "02-no-organizations"
  "03-empty-org-name"
  "04-empty-org-domain"
  "05-invalid-ca-port"
  "06-empty-peer-name"
  "07-empty-peer-subdomain"
  "08-invalid-peer-port"
  "09-empty-orderer-name"
  "10-empty-orderer-subdomain"
  "11-invalid-orderer-port"
  "12-empty-chaincode-name"
  "13-empty-chaincode-path"
  "14-empty-chaincode-version"
  "15-empty-profile-orgs"
  "16-empty-channel-profile"
  "17-empty-channel-name"
)

for config in "${configs[@]}"; do
  echo "******* ${config} *******"
  go run cmd/cli/main.go artifacts generate --config samples/syntactic-constraints/${config}.yml
done