#!/usr/bin/env bash

RED='\033[0;31m'
NC='\033[0m' # No Color
GREEN='\033[0;32m'

if [ -z "$AZ_USER" ]; then
        echo -e "${RED}AZ_USER not set${NC}"
        exit 1
fi

if [ -z "$MASTER_FQDN" ]; then
        echo -e "${RED}MASTER_FQDN not set${NC}"
        exit 1
fi

if [ -z "$KEY_VAULT" ]; then
        echo -e "${RED}KEY_VAULT not set${NC}"
        exit 1
fi

if [ -z "$VAULT_SECRET_NAME" ]; then
        echo -e "${RED}VAULT_SECRET_NAME not set${NC}"
        exit 1
fi

if [ -z "$SP_USER" ]; then
        echo -e "${RED}SP_USER not set${NC}"
        exit 1
fi

if [ -z "$SP_PASS" ]; then
        echo -e "${RED}SP_PASS not set${NC}"
        exit 1
fi

if [ -z "$SP_TENANT" ]; then
        echo -e "${RED}SP_TENANT not set${NC}"
        exit 1
fi

if [ -z "$HELM_REPOSITORY" ]; then
        echo -e "${RED}HELM_REPOSITORY not set${NC}"
        exit 1
fi

if [ -z "$HELM_REPO_NAME" ]; then
        echo -e "${RED}HELM_REPO_NAME not set${NC}"
        exit 1
fi


#login
az login --service-principal --username ${SP_USER} --password ${SP_PASS} --tenant ${SP_TENANT}

RET=$?

#get private key so we can download the kubectl config
if [[ "${RET}" == "0" ]]; then
    rm -r -f ~/.ssh/
    mkdir -p ~/.ssh
    touch ~/.ssh/k8
    #decode base64 instead?
    PRIVATE_KEY_DATA=$(az keyvault secret show --name ${VAULT_SECRET_NAME} --vault-name ${KEY_VAULT} --query "value")
    echo -e "${PRIVATE_KEY_DATA}" | sed -e 's/\"//g' > ~/.ssh/k8
    chmod 400 ~/.ssh/k8
    rm -f -r ~/.kube/config
    scp -oStrictHostKeyChecking=no -i ~/.ssh/k8 ${AZ_USER}@${MASTER_FQDN}:.kube/config ~/.kube
fi

#setup helm and add datamountaineer repo
helm init
helm repo add datamountaineer https://datamountaineer.github.io/helm-charts/
helm repo add ${HELM_REPO_NAME} ${HELM_REPOSITORY}
