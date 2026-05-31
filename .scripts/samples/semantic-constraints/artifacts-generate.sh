#!/bin/bash

configs=(
  "01-unsupported-channel-capability"
  "02-unsupported-application-capability"
  "03-unsupported-orderer-capability"
  "04-duplicate-org-name"
  "05-invalid-peer-version"
  "06-invalid-orderer-version"
  "07-no-orderer-topology"
  "08-multiple-bootstrap-orgs"
  "09-invalid-consensus-type"
  "10-profile-references-undefined"
  "11-exposed-port-conflict"
  "12-invalid-channel-name"
)

for config in "${configs[@]}"; do
  echo "******* ${config} *******"
  go run cmd/cli/main.go artifacts generate --config samples/semantic-constraints/${config}.yml
done