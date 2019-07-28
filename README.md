# nuage-vwan

BACKGROUND

https://docs.microsoft.com/en-us/azure/virtual-wan/

AUTHENTICATING

Let’s start ! However, before we proceed, we first need to somehow authenticate agains Azure API. For Azure this means creating a service principal account that our program will use to authenticate and assume a role with permissions needed to execute API actions.

To create a service principal, let’s use Azure CLI, as shown below. The command will output an authentication file with information such as client id, client secrets and bunch of information needed to connect to Azure. Remember to keep it secure !

az ad sp create-for-rbac —sdk-auth > my.auth

AZURE ENVIRONMENT\n

Initialize authentication token

export AZURE_AUTH_LOCATION="/Users/henderiw/my.auth"

Initialize your resource group name and location

export AZURE_RG_NAME="vWAN"
export AZURE_RG_LOCATION="northeurope"

Initialize the VSD variables

export VSD_URL = "<your URL>"
    example your URL> https://10.0.0.1:8443\n
export VSD_USER = "<user>"
export VSD_PASSWORD = "<password>"
export VSD_ENTERPRISE = "<enterprise>"

Initialize the following variables using the storage account\n

export AZURE_STORAGE_ACCOUNT="<your storage account>"
export AZURE_STORAGE_ACCESS_KEY="<your storage access key>"

YML FILES

Assist in adding sites on vWAN and Nuage

nsg_data:
  enterprise: vspk_public 
  nsg_name: vspkNsgE300WifiLte1
  nsg_port: lte0
  public_ip: 81.246.71.160
  bgp_enabled: false
  bgp_nsg_asn: 1111
  lan_subnet: 
    - 172.0.5.0/24

EXAMPLES

./nuage-vwan -h

create VWAN -> create vwan, vhub and vpnGW on AZURE

./nuage-vwan -o createVWAN

add VWAN site -> create VPN Site and update the VPN GW by adding the site on AZURE

./nuage-vwan -nsg nsg-site1.yml -o addVWANSite

add Nuage site -> create VWAN site on Nuage VSD based on nsg-site1.yml info

./nuage-vwan -nsg nsg-site1.yml -o addNuageSite



