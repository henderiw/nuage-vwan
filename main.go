package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"os"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/network/mgmt/network"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/azure/auth"
	"github.com/Azure/go-autorest/autorest/to"
	azurewrapper "github.com/henderiw/azure-wrapper"
	nuagewrapper "github.com/henderiw/nuage-wrapper"
	"github.com/nuagenetworks/go-bambou/bambou"
	"github.com/nuagenetworks/vspk-go/vspk"
)

// VSD credential parameters
var vsdURL, vsdUser, vsdPass, vsdEnterprise string

// Usr is a user
var Usr *vspk.Me

// Authenticate with the Azure services using file-based authentication
func init() {
	var err error
	azurewrapper.Authorizer, err = auth.NewAuthorizerFromFile(azure.PublicCloud.ResourceManagerEndpoint)
	if err != nil {
		log.Fatalf("Failed to get OAuth config: %v", err)
	}

	authInfo, err := azurewrapper.ReadJSON(os.Getenv("AZURE_AUTH_LOCATION"))
	if err != nil {
		log.Fatalf("Failed to get data from AZURE_AUTH_LOCATION env: %v", err)
	}
	azurewrapper.ClientData.SubscriptionID = (*authInfo)["subscriptionId"].(string)
	azurewrapper.ClientData.ResourceGroupName = os.Getenv("AZURE_RG_NAME")
	azurewrapper.ClientData.ResourceGroupLocation = os.Getenv("AZURE_RG_LOCATION")
	vsdURL = os.Getenv("VSD_URL")
	vsdUser = os.Getenv("VSD_USER")
	vsdPass = os.Getenv("VSD_PASSWORD")
	vsdEnterprise = os.Getenv("VSD_ENTERPRISE")

}

func createVwanWorkflow(vwanName, vwanHubName, vpnGWName, vhubLocation string) {
	vwan, _ := azurewrapper.CreateVwan(vwanName)
	log.Printf("vwan Created Name: %s \n", to.String(vwan.Name))

	vhub, _ := azurewrapper.CreateVhub(vwanHubName, to.String(vwan.ID), "10.1.1.0/24", vhubLocation)
	log.Printf("vhub Created Name: %s \n", to.String(vhub.Name))
	log.Printf("vhub Created ID: %s \n", to.String(vhub.ID))

	var nsgConf azurewrapper.NsgConfYML
	vpnGW, _ := azurewrapper.CreateVpnGateway(vpnGWName, to.String(vhub.ID), "", vhubLocation, nsgConf)
	log.Printf("vpnGW Created ID: %s \n", to.String(vpnGW.ID))
	log.Printf("vpnGW Created Name: %s \n", to.String(vpnGW.Name))

}

func deleteVwanWorkflow(vwanName, vwanHubName, vpnGWName string) {
	vpnGW, _ := azurewrapper.GetVpnGateway(vpnGWName)
	log.Printf("vpnGW Name: %s \n", to.String(vpnGW.Name))
	log.Printf("vpnGW ID: %s \n", to.String(vpnGW.ID))

	if vpnGW.ID != nil {
		if vpnGW.VpnGatewayProperties.Connections != nil {
			idx := 0
			for i, connection := range *vpnGW.VpnGatewayProperties.Connections {
				log.Printf("Before deleting vWAN delete this Site first: %s \n", to.String(connection.RemoteVpnSite.ID))
				idx = i
			}
			if idx != 0 {
				log.Printf("Site Connection exists on GW")
				return
			}

		}
	}

	azurewrapper.DeleteVpnGateway(vpnGWName)
	azurewrapper.DeleteVhub(vwanHubName)
	azurewrapper.DeleteVwan(vwanName)
}

func getVwanWorkflow(vwanName, vwanHubName, vpnGWName string) {
	vwan, _ := azurewrapper.GetVwan(vwanName)
	log.Printf("VWAN Name: %s \n", to.String(vwan.Name))
	log.Printf("VWAN ID: %s \n", to.String(vwan.ID))

	if vwan.ID != nil {
		log.Printf("VWAN DisableVpnEncryption: %#v \n", to.Bool(vwan.VirtualWanProperties.DisableVpnEncryption))
		if vwan.VirtualWanProperties.VpnSites != nil {
			for _, site := range *vwan.VirtualWanProperties.VpnSites {
				log.Printf("Site: %s \n", to.String(site.ID))
			}
		}

		if vwan.VirtualWanProperties.VpnSites != nil {
			for _, hub := range *vwan.VirtualWanProperties.VirtualHubs {
				log.Printf("Hub: %s \n", to.String(hub.ID))
			}
		}
	}
	vhub, _ := azurewrapper.GetVhub(vwanHubName)
	log.Printf("VHUB Name: %s \n", to.String(vhub.Name))
	log.Printf("VHUB ID: %s \n", to.String(vhub.ID))

	vpnGW, _ := azurewrapper.GetVpnGateway(vpnGWName)
	log.Printf("vpnGW Name: %s \n", to.String(vpnGW.Name))
	log.Printf("vpnGW ID: %s \n", to.String(vpnGW.ID))

	if vpnGW.ID != nil {
		if vpnGW.VpnGatewayProperties.Connections != nil {
			for _, connection := range *vpnGW.VpnGatewayProperties.Connections {
				log.Printf("Site: %s \n", to.String(connection.RemoteVpnSite.ID))
			}
		}
	}
}

func createVhubWorkflow(vwanName, vwanHubName, vpnGWName, ipSubnet, vhubLocation string) {
	vwan, _ := azurewrapper.GetVwan(vwanName)
	log.Printf("vwan Name: %s \n", to.String(vwan.Name))

	vhub, _ := azurewrapper.CreateVhub(vwanHubName, to.String(vwan.ID), ipSubnet, vhubLocation)
	log.Printf("vhub Created Name: %s \n", to.String(vhub.Name))
	log.Printf("vhub Created ID: %s \n", to.String(vhub.ID))

	var nsgConf azurewrapper.NsgConfYML
	log.Printf("vhub.ID ID: %s \n", to.String(vhub.ID))
	log.Printf("nsgConf: %s \n", nsgConf)
	vpnGW, _ := azurewrapper.CreateVpnGateway(vpnGWName, to.String(vhub.ID), "", vhubLocation, nsgConf)
	log.Printf("vpnGW Created ID: %s \n", to.String(vpnGW.ID))
	log.Printf("vpnGW Created Name: %s \n", to.String(vpnGW.Name))

}

func deleteVhubWorkflow(vwanName, vwanHubName, vpnGWName string) {
	vpnGW, _ := azurewrapper.GetVpnGateway(vpnGWName)
	log.Printf("vpnGW Name: %s \n", to.String(vpnGW.Name))
	log.Printf("vpnGW ID: %s \n", to.String(vpnGW.ID))

	if vpnGW.ID != nil {
		if vpnGW.VpnGatewayProperties.Connections != nil {
			idx := 0
			for i, connection := range *vpnGW.VpnGatewayProperties.Connections {
				log.Printf("Before deleting vWAN delete this Site first: %s \n", to.String(connection.RemoteVpnSite.ID))
				idx = i
			}
			if idx != 0 {
				log.Printf("Site Connection exists on GW")
				return
			}

		}
	}
	azurewrapper.DeleteVpnGateway(vpnGWName)
	azurewrapper.DeleteVhub(vwanHubName)
}

func getVhubWorkflow(vwanName, vwanHubName, vpnGWName string) {
	vwan, _ := azurewrapper.GetVwan(vwanName)
	log.Printf("VWAN Name: %s \n", to.String(vwan.Name))
	log.Printf("VWAN ID: %s \n", to.String(vwan.ID))

	vhub, _ := azurewrapper.GetVhub(vwanHubName)
	log.Printf("VHUB Name: %s \n", to.String(vhub.Name))
	log.Printf("VHUB ID: %s \n", to.String(vhub.ID))

	vpnGW, _ := azurewrapper.GetVpnGateway(vpnGWName)
	log.Printf("vpnGW Name: %s \n", to.String(vpnGW.Name))
	log.Printf("vpnGW ID: %s \n", to.String(vpnGW.ID))

	if vpnGW.ID != nil {
		if vpnGW.VpnGatewayProperties.Connections != nil {
			for _, connection := range *vpnGW.VpnGatewayProperties.Connections {
				log.Printf("Site: %s \n", to.String(connection.RemoteVpnSite.ID))
			}
		}
	}
}

func addVwanSiteWorkflow(vwanName, vwanHubName, vpnGWName, vhubLocation string, nsgConf azurewrapper.NsgConfYML) {
	vwan, _ := azurewrapper.GetVwan(vwanName)
	log.Printf("VWAN Name: %s \n", to.String(vwan.Name))
	log.Printf("VWAN ID: %s \n", to.String(vwan.ID))

	vhub, _ := azurewrapper.GetVhub(vwanHubName)
	log.Printf("VHUB Name: %s \n", to.String(vhub.Name))
	log.Printf("VHUB ID: %s \n", to.String(vhub.ID))

	vsite, _ := azurewrapper.CreateVpnSite(nsgConf.NsgData.NsgName, to.String(vwan.ID), vhubLocation, nsgConf)
	log.Printf("vsite1 Created Name: %s \n", nsgConf.NsgData.NsgName)

	log.Printf("vpnGWName: %s \n", vpnGWName)
	log.Printf("vhub ID: %s \n", to.String(vhub.ID))
	log.Printf("vsite ID: %s \n", to.String(vsite.ID))
	log.Printf("nsg CONF: %s \n", nsgConf)

	vpnGW, _ := azurewrapper.CreateVpnGateway(vpnGWName, to.String(vhub.ID), to.String(vsite.ID), vhubLocation, nsgConf)
	log.Printf("vpnGW Created ID: %s \n", to.String(vpnGW.ID))
	log.Printf("vpnGW Created Name: %s \n", to.String(vpnGW.Name))

	blobName := nsgConf.NsgData.NsgName

	url, _ := azurewrapper.CreateStorageSasURL("vpnconfig", blobName)

	azurewrapper.DownloadVpnSitesConfig(to.String(vwan.Name), to.String(vsite.ID), url)

	azurewrapper.DownloadFileFromURL(url, blobName)

	log.Println(blobName + " saved!")

	rcvdAzureVWanData, readErr := ioutil.ReadFile(blobName)
	if readErr != nil {
		log.Fatal(readErr)
	}

	log.Printf("File contents: %s\n", rcvdAzureVWanData)

	// init the empty structure
	var cfg []azurewrapper.AzureVWanCfg

	// unmarshal (deserialize) the json and save the result in the struct &cfg
	err := json.Unmarshal([]byte(rcvdAzureVWanData), &cfg)
	if err != nil {
		log.Fatal(err)
	}

	vwanHubIP1 := cfg[0].VpnSiteConnections[0].GatewayConfiguration.IPAddresses.Instance0
	vwanHubIP2 := cfg[0].VpnSiteConnections[0].GatewayConfiguration.IPAddresses.Instance1
	vwanSite1Name := cfg[0].VpnSiteConfiguration.Name

	log.Printf("vwanHubIP1: %s\n", vwanHubIP1)
	log.Printf("vwanHubIP2: %s\n", vwanHubIP2)
	log.Printf("vwanSiteName: %s\n", vwanSite1Name)

}

func deleteVwanSiteWorkflow(vwanName, vwanHubName, vpnGWName, vhubLocation string, nsgConf azurewrapper.NsgConfYML) {
	vpnSite, _ := azurewrapper.GetVpnSite(nsgConf.NsgData.NsgName)
	log.Printf("vpnSite Name: %s \n", to.String(vpnSite.Name))
	log.Printf("vpnSite ID: %s \n", to.String(vpnSite.ID))

	vpnGW, _ := azurewrapper.GetVpnGateway(vpnGWName)
	log.Printf("vpnGW Name: %s \n", to.String(vpnGW.Name))
	log.Printf("vpnGW ID: %s \n", to.String(vpnGW.ID))

	var newVpnConnections []network.VpnConnection
	log.Println("Before newVpnConnections:", newVpnConnections)

	if vpnGW.ID != nil {
		if vpnGW.VpnGatewayProperties.Connections != nil {
			log.Println("Before:", *vpnGW.VpnGatewayProperties.Connections)
			for _, connection := range *vpnGW.VpnGatewayProperties.Connections {
				//log.Printf("Before deleting vWAN delete this Site first: %s \n", to.String(connection.RemoteVpnSite.ID))
				if to.String(vpnSite.ID) != to.String(connection.RemoteVpnSite.ID) {
					newVpnConnections = append(newVpnConnections, connection)
				}

			}
			log.Println("After:", newVpnConnections)
			*vpnGW.VpnGatewayProperties.Connections = newVpnConnections
			log.Println("After:", *vpnGW.VpnGatewayProperties.Connections)
			for _, connection := range *vpnGW.VpnGatewayProperties.Connections {
				log.Printf("New VPN GW Connections: %s \n", to.String(connection.ID))
				log.Printf("New VPN GW Connections: %s \n", connection.ProvisioningState)
				log.Printf("New VPN GW Connections: %s \n", to.String(connection.SharedKey))
				log.Printf("New VPN GW Connections: %s \n", to.String(connection.ID))
			}
		}
	}
	vpnGW, _ = azurewrapper.UpdateVpnGateway(vpnGWName, to.String(vpnGW.VirtualHub.ID), vhubLocation, newVpnConnections)

	if vpnGW.ID != nil {
		if vpnGW.VpnGatewayProperties.Connections != nil {
			for _, connection := range *vpnGW.VpnGatewayProperties.Connections {
				log.Printf("Sites remaining on VPN GW: %s \n", to.String(connection.RemoteVpnSite.ID))
			}
		}
	}
	azurewrapper.DeleteVpnSite(nsgConf.NsgData.NsgName)
}

func addNuageSiteWorkflow(vhub string, nsgConf azurewrapper.NsgConfYML) {
	rcvdAzureVWanData, readErr := ioutil.ReadFile(nsgConf.NsgData.NsgName)
	if readErr != nil {
		log.Fatal(readErr)
	}

	log.Printf("File contents: %s\n", rcvdAzureVWanData)

	// init the empty structure
	var cfg []azurewrapper.AzureVWanCfg

	// unmarshal (deserialize) the json and save the result in the struct &cfg
	err := json.Unmarshal([]byte(rcvdAzureVWanData), &cfg)
	if err != nil {
		log.Fatal(err)
	}

	vwanHubIP1 := cfg[0].VpnSiteConnections[0].GatewayConfiguration.IPAddresses.Instance0
	vwanHubIP2 := cfg[0].VpnSiteConnections[0].GatewayConfiguration.IPAddresses.Instance1
	//vwanSiteName := cfg[0].VpnSiteConfiguration.Name
	//vwanSiteIP := cfg[0].VpnSiteConfiguration.IPAddress
	//vwanSiteBGPEnabled := cfg[0].VpnSiteConnections[0].ConnectionConfiguration.IsBgpEnabled
	vwanSitePSK := cfg[0].VpnSiteConnections[0].ConnectionConfiguration.PSK
	//vwanSiteSAsize := cfg[0].VpnSiteConnections[0].ConnectionConfiguration.IPsecParameters.SADataSizeInKilobytes
	//vwanSiteSALifeTime := cfg[0].VpnSiteConnections[0].ConnectionConfiguration.IPsecParameters.SALifeTimeInSeconds

	//Start session to VSD
	var s *bambou.Session
	s, Usr = vspk.NewSession(vsdUser, vsdPass, vsdEnterprise, vsdURL)
	if err := s.Start(); err != nil {
		log.Println("Unable to connect to Nuage VSD: " + err.Description)
		os.Exit(1)
	}

	//find enterprise
	enterpriseCfg := map[string]interface{}{
		"Name": nsgConf.NsgData.Enterprise,
	}

	enterprise := nuagewrapper.NuageEnterprise(enterpriseCfg, Usr)

	//create or get PSK

	ikePSKCfg := map[string]interface{}{
		"Name":           vhub + "vspkdemoAzure",
		"Description":    vhub + "vspkdemoAzure",
		"UnencryptedPSK": vwanSitePSK,
	}

	ikePSK := nuagewrapper.NuageCreateIKEPSK(ikePSKCfg, enterprise)

	//create IKE Gateway(s)

	ikeGatewayCfg1 := map[string]interface{}{
		"Name":        vhub + "vspkdemoAzureIKEGatewayName1",
		"Description": vhub + "vspkdemoAzureIKEGatewayName1",
		"IKEVersion":  "V2",
		"IPAddress":   vwanHubIP1,
	}

	ikeGateway1 := nuagewrapper.NuageCreateIKEGateway(ikeGatewayCfg1, enterprise)

	ikeGatewayCfg2 := map[string]interface{}{
		"Name":        vhub + "vspkdemoAzureIKEGatewayName2",
		"Description": vhub + "vspkdemoAzureIKEGatewayName2",
		"IKEVersion":  "V2",
		"IPAddress":   vwanHubIP2,
	}

	ikeGateway2 := nuagewrapper.NuageCreateIKEGateway(ikeGatewayCfg2, enterprise)

	//Create IKE Encryption Profile

	ikeEncryptionProfileCfg := map[string]interface{}{
		"Name":                              vhub + "vspkdemoAzureIKEEncryptionProfile",
		"Description":                       vhub + "vspkdemoAzureIKEEncryptionProfile",
		"DPDMode":                           "REPLY_ONLY",
		"ISAKMPAuthenticationMode":          "PRE_SHARED_KEY",
		"ISAKMPDiffieHelmanGroupIdentifier": "GROUP_2_1024_BIT_DH",
		"ISAKMPEncryptionAlgorithm":         "AES256",
		"ISAKMPEncryptionKeyLifetime":       28800,
		"ISAKMPHashAlgorithm":               "SHA256",
		"IPsecEnablePFS":                    true,
		"IPsecEncryptionAlgorithm":          "AES256",
		"IPsecPreFragment":                  true,
		"IPsecSALifetime":                   3600,
		"IPsecAuthenticationAlgorithm":      "HMAC_SHA256",
		"IPsecSAReplayWindowSize":           "WINDOW_SIZE_64",
	}

	ikeEncryptionProfile := nuagewrapper.NuageCreateIKEEncryptionProfile(ikeEncryptionProfileCfg, enterprise)

	//Create IKE Gateway Profile

	ikeGatewayProfileCfg1 := map[string]interface{}{
		"Name":                             vhub + "vspkdemoAzureIKEGatewayProfile1",
		"Description":                      vhub + "vspkdemoAzureIKEGatewayProfile1",
		"AssociatedIKEAuthenticationID":    ikePSK.ID,
		"IKEGatewayIdentifier":             vwanHubIP1,
		"IKEGatewayIdentifierType":         "ID_IPV4_ADDR",
		"AssociatedIKEGatewayID":           ikeGateway1.ID,
		"AssociatedIKEEncryptionProfileID": ikeEncryptionProfile.ID,
	}

	ikeGatewayProfile1 := nuagewrapper.NuageCreateIKEGatewayProfile(ikeGatewayProfileCfg1, enterprise)

	ikeGatewayProfileCfg2 := map[string]interface{}{
		"Name":                             vhub + "vspkdemoAzureIKEGatewayProfile2",
		"Description":                      vhub + "vspkdemoAzureIKEGatewayProfile2",
		"AssociatedIKEAuthenticationID":    ikePSK.ID,
		"IKEGatewayIdentifier":             vwanHubIP2,
		"IKEGatewayIdentifierType":         "ID_IPV4_ADDR",
		"AssociatedIKEGatewayID":           ikeGateway2.ID,
		"AssociatedIKEEncryptionProfileID": ikeEncryptionProfile.ID,
	}

	ikeGatewayProfile2 := nuagewrapper.NuageCreateIKEGatewayProfile(ikeGatewayProfileCfg2, enterprise)

	nsGatewayCfg := map[string]interface{}{
		"Name":                  nsgConf.NsgData.NsgName,
		"TCPMSSEnabled":         true,
		"TCPMaximumSegmentSize": 1330,
		"NetworkAcceleration":   "NONE",
	}

	nsGateway := nuagewrapper.NuageNSG(nsGatewayCfg, enterprise)

	nsPortCfg := map[string]interface{}{
		"Name": nsgConf.NsgData.NsgPort,
	}

	nsPort := nuagewrapper.NuageNSGPort(nsPortCfg, nsGateway)

	nsVlanCfg := map[string]interface{}{
		"Value": 0,
	}

	nsVlan := nuagewrapper.NuageVlan(nsVlanCfg, nsPort)

	ikeGatewayConnCfg1 := map[string]interface{}{
		"Name":                          vhub + "vspkdemoAzureIKEGatewayConnection1",
		"Description":                   vhub + "vspkdemoAzureIKEGatewayConnection1",
		"NSGIdentifier":                 nsgConf.NsgData.NsgName,
		"NSGIdentifierType":             "ID_KEY_ID",
		"NSGRole":                       "INITIATOR",
		"AllowAnySubnet":                true,
		"AssociatedIKEGatewayProfileID": ikeGatewayProfile1.ID,
		"AssociatedIKEAuthenticationID": ikePSK.ID,
	}

	ikeGatewayConn1 := nuagewrapper.NuageCreateIKEGatewayConnection(ikeGatewayConnCfg1, nsVlan)
	log.Println(ikeGatewayConn1)

	ikeGatewayConnCfg2 := map[string]interface{}{
		"Name":                          vhub + "vspkdemoAzureIKEGatewayConnection2",
		"Description":                   vhub + "vspkdemoAzureIKEGatewayConnection2",
		"NSGIdentifier":                 nsgConf.NsgData.NsgName,
		"NSGIdentifierType":             "ID_KEY_ID",
		"NSGRole":                       "INITIATOR",
		"AllowAnySubnet":                true,
		"AssociatedIKEGatewayProfileID": ikeGatewayProfile2.ID,
		"AssociatedIKEAuthenticationID": ikePSK.ID,
	}

	ikeGatewayConn2 := nuagewrapper.NuageCreateIKEGatewayConnection(ikeGatewayConnCfg2, nsVlan)
	log.Println(ikeGatewayConn2)

}

func deleteNuageSiteWorkflow(vhub string, nsgConf azurewrapper.NsgConfYML) {

	//Start session to VSD
	var s *bambou.Session
	s, Usr = vspk.NewSession(vsdUser, vsdPass, vsdEnterprise, vsdURL)
	if err := s.Start(); err != nil {
		log.Println("Unable to connect to Nuage VSD: " + err.Description)
		os.Exit(1)
	}

	//find enterprise
	enterpriseCfg := map[string]interface{}{
		"Name": nsgConf.NsgData.Enterprise,
	}

	enterprise := nuagewrapper.NuageEnterprise(enterpriseCfg, Usr)

	nsGatewayCfg := map[string]interface{}{
		"Name": nsgConf.NsgData.NsgName,
	}

	nsGateway := nuagewrapper.NuageNSG(nsGatewayCfg, enterprise)

	nsPortCfg := map[string]interface{}{
		"Name": nsgConf.NsgData.NsgPort,
	}

	nsPort := nuagewrapper.NuageNSGPort(nsPortCfg, nsGateway)

	nsVlanCfg := map[string]interface{}{
		"Value": 0,
	}

	nsVlan := nuagewrapper.NuageVlan(nsVlanCfg, nsPort)

	ikeGatewayConnCfg1 := map[string]interface{}{
		"Name": vhub + "vspkdemoAzureIKEGatewayConnection1",
	}

	err1 := nuagewrapper.NuageDeleteIKEGatewayConnection(ikeGatewayConnCfg1, nsVlan)
	log.Println(err1)

	ikeGatewayConnCfg2 := map[string]interface{}{
		"Name": vhub + "vspkdemoAzureIKEGatewayConnection2",
	}

	err2 := nuagewrapper.NuageDeleteIKEGatewayConnection(ikeGatewayConnCfg2, nsVlan)
	log.Println(err2)
}

func main() {

	var vwan, vhub, vhubSubnet, vhubLocation, enterprise, nsg, operation string
	flag.StringVar(&vwan, "vwan", "test", "vwan name")
	flag.StringVar(&vhub, "vhub", "1", "vhub id")
	flag.StringVar(&vhubSubnet, "vhubSubnet", "10.1.1.0/24", "vhub Subnnet")
	flag.StringVar(&vhubLocation, "vhubLocation", "northeurope", "vhub Location")
	flag.StringVar(&enterprise, "enterprise", "vspk_public", "enterprise name")
	flag.StringVar(&nsg, "nsg", "nsg-site1.yml", "nsg name or group or all nsgs in the enterprise")
	flag.StringVar(&operation, "o", "getVWAN", "createVWAN, deleteVWAN, getVWAN, createVHUB, deleteVHUB, getVHUB, addVWANSite, deleteVWANSite, addNuageSite, deleteNuageSite")

	flag.Parse()

	log.Println("vwan:", vwan)
	log.Println("vhub:", vhub)
	log.Println("vhub Subnet:", vhubSubnet)
	log.Println("vhub Location:", vhubLocation)
	log.Println("enterprise:", enterprise)
	log.Println("nsg:", nsg)
	log.Println("operation:", operation)

	vwanName := "vwan" + vwan
	vwanHubName := "vwanHub" + vwan + vhub
	vpnGWName := "vpnGw" + vwan + vhub

	switch operation {
	case "createVWAN":
		log.Println("Create VWAN Workflow")
		createVwanWorkflow(vwanName, vwanHubName, vpnGWName, vhubLocation)
		break
	case "deleteVWAN":
		log.Println("Delete VWAN Workflow")
		deleteVwanWorkflow(vwanName, vwanHubName, vpnGWName)
		break
	case "getVWAN":
		log.Println("get VWAN Workflow")
		getVwanWorkflow(vwanName, vwanHubName, vpnGWName)
		break
	case "createVHUB":
		log.Println("Create VHUB Workflow")
		createVhubWorkflow(vwanName, vwanHubName, vpnGWName, vhubSubnet, vhubLocation)
		break
	case "deleteVHUB":
		log.Println("Delete VHUB Workflow")
		deleteVhubWorkflow(vwanName, vwanHubName, vpnGWName)
		break
	case "getVHUB":
		log.Println("get VHUB Workflow")
		getVhubWorkflow(vwanName, vwanHubName, vpnGWName)
		break
	case "addVWANSite":
		log.Println("add VWAN Site Workflow")
		if nsg != "all" {
			var nsgConf azurewrapper.NsgConfYML
			nsgConf.GetConf(nsg)
			log.Println(nsgConf)
			addVwanSiteWorkflow(vwanName, vwanHubName, vpnGWName, vhubLocation, nsgConf)
		}
		break
	case "deleteVWANSite":
		log.Println("delete VWAN Site Workflow")
		if nsg != "all" {
			var nsgConf azurewrapper.NsgConfYML
			nsgConf.GetConf(nsg)
			deleteVwanSiteWorkflow(vwanName, vwanHubName, vpnGWName, vhubLocation, nsgConf)
		}
		break
	case "addNuageSite":
		log.Println("add Nuage Site Workflow")
		if nsg != "all" {
			var nsgConf azurewrapper.NsgConfYML
			nsgConf.GetConf(nsg)
			log.Println(nsgConf)
			addNuageSiteWorkflow(vwanHubName, nsgConf)
		}
		break
	case "deleteNuageSite":
		log.Println("delete Nuage Site Workflow")
		if nsg != "all" {
			var nsgConf azurewrapper.NsgConfYML
			nsgConf.GetConf(nsg)
			deleteNuageSiteWorkflow(vwanHubName, nsgConf)
		}
		break
	default:
		log.Fatalln("Wrong Operation Input (create, delete, add or get)")

	}
}
