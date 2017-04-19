# Azure CLI, Helm, Kubectl and the Landscaper

This docker contains the [azure cli](https://docs.microsoft.com/en-us/cli/azure/install-azure-cli),
 [kubectl](https://kubernetes.io/docs/tasks/kubectl/install/), [Helm](https://helm.sh/) and 
 Eneco's [Landscaper](https://github.com/Eneco/landscaper) to managing and deploying ``landscapes`` via Helm.
 
 The ``configure.sh`` script initializes helm and pulls the kubectl config from the target cluster. 
 
 The script requires the following enviroment variables:
 
*   AZ_USER - azure user for the kubernetes master node vm
*   MASTER_FQDN - the fqdn of the master vm
*   KEY_VAULT - the key vault to retrieve the private key from
*   VAULT_SECRET_NAME - the secret the name for the private key
*   SP_USER - the service principal to fetch the private key with
*   SP_PASS - the service principal password
*   SP_TENANT - the service principal tenant id
*   HELM_REPO_NAME - the name of the repo to add to helm (helm repo add ${HELM_REPO_NAME} ${HELM_REPOSITORY})
*   HELM_REPOSITORY - the helm repo url to add

The script will also add https://datamountaineer.github.io/helm-charts/ to helm.