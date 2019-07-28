# nuage-vwan

AUTHENTICATING

Let’s start ! However, before we proceed, we first need to somehow authenticate agains Azure API. For Azure this means creating a service principal account that our program will use to authenticate and assume a role with permissions needed to execute API actions.

To create a service principal, let’s use Azure CLI, as shown below. The command will output an authentication file with information such as client id, client secrets and bunch of information needed to connect to Azure. Remember to keep it secure !

az ad sp create-for-rbac —sdk-auth > my.auth

AZURE ENVIRONMENT\n

Initialize authentication token\n

export AZURE_AUTH_LOCATION="/Users/henderiw/my.auth"\n

Initialize your resource group name and location\n

export AZURE_RG_NAME="vWAN"\n
export AZURE_RG_LOCATION="northeurope"\n

Initialize the VSD variables\n

export VSD_URL = "<your URL>"\n
    example your URL> https://10.0.0.1:8443\n
export VSD_USER = "<user>"\n
export VSD_PASSWORD = "<password>"\n
export VSD_ENTERPRISE = "<enterprise>"\n

Initialize the following variables using the storage account\n

export AZURE_STORAGE_ACCOUNT="<your storage account>"\n
export AZURE_STORAGE_ACCESS_KEY="<your storage access key>"\n

YML FILES\n



