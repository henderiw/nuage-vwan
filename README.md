# nuage-vwan

AUTHENTICATING
Let’s start ! However, before we proceed, we first need to somehow authenticate agains Azure API. For Azure this means creating a service principal account that our program will use to authenticate and assume a role with permissions needed to execute API actions.

To create a service principal, let’s use Azure CLI, as shown below. The command will output an authentication file with information such as client id, client secrets and bunch of information needed to connect to Azure. Remember to keep it secure !

az ad sp create-for-rbac —sdk-auth > my.auth

AZURE ENVIRONMENT

Initialize authentication token
export AZURE_AUTH_LOCATION="/Users/henderiw/my.auth"

Initialize your resource group name and location
export AZURE_RG_NAME="vWAN"
export AZURE_RG_LOCATION="northeurope"

Initialize the VSD variables
export VSD_URL = "<your URL>"
    example your URL> https://10.0.0.1:8443
export VSD_USER = "<user>"
export VSD_PASSWORD = "<password>"
export VSD_ENTERPRISE = "<enterprise>"

Initialize the following variables using the storage account
export AZURE_STORAGE_ACCOUNT="<your storage account>"
export AZURE_STORAGE_ACCESS_KEY="<your storage access key>"