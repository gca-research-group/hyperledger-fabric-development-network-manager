#!/bin/bash
source ./.scripts/config/_colors.sh
source ./.scripts/config/_icons.sh

verifyIfConfigTxFileExists() {
  echo -e "${PROCESSING_ICON} Verifying if the config file exists."
  if [[ ! -f $1 ]]; then
    echo -e "${RED}The configtx file has not been found:${NO_COLOR} $1"
    exit 1
  fi
}

verifyIfCryptoMaterialsExist() {
  echo -e "${PROCESSING_ICON} Verifying if the crypto materials exist."

  if [[ ! -d $1 ]]; then
    echo -e "${RED}Crypto materials do not exists:${NO_COLOR} $1"
    exit 1
  fi
}
