package integration

import (
	"fmt"
	//"os"
	"path"
	"testing"
	"time"

	"github.com/hyperledger/fabric-sdk-go/api/apiconfig"
	ca "github.com/hyperledger/fabric-sdk-go/api/apifabca"
	fab "github.com/hyperledger/fabric-sdk-go/api/apifabclient"
	"github.com/hyperledger/fabric-sdk-go/api/apitxn"
	chmgmt "github.com/hyperledger/fabric-sdk-go/api/apitxn/chmgmtclient"
	resmgmt "github.com/hyperledger/fabric-sdk-go/api/apitxn/resmgmtclient"
	//packager "github.com/hyperledger/fabric-sdk-go/pkg/fabric-client/ccpackager/gopackager"
	//pb "github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/protos/peer"

	deffab "github.com/hyperledger/fabric-sdk-go/def/fabapi"
	"github.com/hyperledger/fabric-sdk-go/def/fabapi/opt"
	//"github.com/hyperledger/fabric-sdk-go/pkg/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/errors"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabric-client/events"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabric-client/orderer"
	//"github.com/hyperledger/fabric-sdk-go/pkg/fabric-client/peer"
	//"github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/common/cauthdsl"
)

// BaseSetupImpl implementation of BaseTestSetup
type BaseSetupImpl struct {
	Client          fab.FabricClient
	Channel         fab.Channel
	EventHub        fab.EventHub
	ConnectEventHub bool
	ConfigFile      string
	OrgID           string
	ChannelID       string
	ChainCodeID     string
	Initialized     bool
	ChannelConfig   string
	AdminUser       ca.User
}

// GetChannel initializes and returns a channel based on config
func (setup *BaseSetupImpl) GetChannel(client fab.FabricClient, channelID string, orgs []string) (fab.Channel, error) {

	channel, err := client.NewChannel(channelID)
	if err != nil {
		return nil, errors.WithMessage(err, "NewChannel failed")
	}

	ordererConfig, err := client.Config().RandomOrdererConfig()
	if err != nil {
		return nil, errors.WithMessage(err, "RandomOrdererConfig failed")
	}

	orderer, err := orderer.NewOrdererFromConfig(ordererConfig, client.Config())
	if err != nil {
		return nil, errors.WithMessage(err, "NewOrderer failed")
	}
	err = channel.AddOrderer(orderer)
	if err != nil {
		return nil, errors.WithMessage(err, "adding orderer failed")
	}

	for _, org := range orgs {
		peerConfig, err := client.Config().PeersConfig(org)
		if err != nil {
			return nil, errors.WithMessage(err, "reading peer config failed")
		}
		for _, p := range peerConfig {
			endorser, err := deffab.NewPeerFromConfig(&apiconfig.NetworkPeer{PeerConfig: p}, client.Config())
			if err != nil {
				return nil, errors.WithMessage(err, "NewPeer failed")
			}
			err = channel.AddPeer(endorser)
			if err != nil {
				return nil, errors.WithMessage(err, "adding peer failed")
			}
		}
	}

	return channel, nil
}

// HasPrimaryPeerJoinedChannel checks whether the primary peer of a channel
// has already joined the channel. It returns true if it has, false otherwise,
// or an error
func HasPrimaryPeerJoinedChannel(client fab.FabricClient, channel fab.Channel) (bool, error) {
	foundChannel := false
	primaryPeer := channel.PrimaryPeer()
	response, err := client.QueryChannels(primaryPeer)
	if err != nil {
		return false, errors.WithMessage(err, "failed to query channel for primary peer")
	}
	for _, responseChannel := range response.Channels {
		if responseChannel.ChannelId == channel.Name() {
			foundChannel = true
		}
	}

	return foundChannel, nil
}

func (setup *BaseSetupImpl) setupEventHub(client fab.FabricClient) error {
	eventHub, err := setup.getEventHub(client)
	if err != nil {
		return err
	}

	if setup.ConnectEventHub {
		if err := eventHub.Connect(); err != nil {
			return errors.WithMessage(err, "eventHub connect failed")
		}
	}
	setup.EventHub = eventHub

	return nil
}

// getEventHub initilizes the event hub
func (setup *BaseSetupImpl) getEventHub(client fab.FabricClient) (fab.EventHub, error) {
	eventHub, err := events.NewEventHub(client)
	if err != nil {
		return nil, errors.WithMessage(err, "NewEventHub failed")
	}
	foundEventHub := false
	peerConfig, err := client.Config().PeersConfig(setup.OrgID)
	if err != nil {
		return nil, errors.WithMessage(err, "PeersConfig failed")
	}
	for _, p := range peerConfig {
		if p.URL != "" {
			fmt.Printf("EventHub connect to peer (%s)\n", p.URL)
			serverHostOverride := ""
			if str, ok := p.GRPCOptions["ssl-target-name-override"].(string); ok {
				serverHostOverride = str
			}
			eventHub.SetPeerAddr(p.EventURL, p.TLSCACerts.Path, serverHostOverride)
			foundEventHub = true
			break
		}
	}

	if !foundEventHub {
		return nil, errors.New("event hub configuration not found")
	}

	return eventHub, nil
}

func (setup *BaseSetupImpl) ensureJoinChannel(sdk *deffab.FabricSDK) error {
	// Channel management client is responsible for managing channels (create/update)
	chMgmtClient, err := sdk.NewChannelMgmtClientWithOpts("Admin",
		&deffab.ChannelMgmtClientOpts{OrgName: "ordererorg"})
	if err != nil {
		//t.Fatalf("Failed to create new channel management client: %s", err)
		return errors.WithMessage(err, "Failed to create new channel management client")
	}

	var resMgmtClient resmgmt.ResourceMgmtClient

	// Resource management client is responsible for managing resources (joining channels, install/instantiate/upgrade chaincodes)
	resMgmtClient, err = sdk.NewResourceMgmtClient("Admin")
	if err != nil {
		//t.Fatalf("Failed to create new resource management client: %s", err)
		return errors.WithMessage(err, "Failed to create new channel resource management client")
	}

	// Check if primary peer has joined channel
	alreadyJoined, err := HasPrimaryPeerJoinedChannel(setup.Client, setup.Channel)
	if err != nil {
		return errors.WithMessage(err, "failed while checking if primary peer has already joined channel")
	}

	if !alreadyJoined {

		// Channel config signing user (has to belong to one of channel orgs)
		org1Admin, err := sdk.NewPreEnrolledUser("Org1", "Admin")
		if err != nil {
			return errors.WithMessage(err, "failed getting Org1 admin user")
		}

		// Create channel (or update if it already exists)
		req := chmgmt.SaveChannelRequest{
			ChannelID:     setup.ChannelID,
			ChannelConfig: setup.ChannelConfig,
			SigningUser:   org1Admin,
		}

		if err = chMgmtClient.SaveChannel(req); err != nil {
			return errors.WithMessage(err, "SaveChannel failed")
		}

		time.Sleep(time.Second * 3)

		if err = setup.Channel.Initialize(nil); err != nil {
			return errors.WithMessage(err, "channel init failed")
		}

		if err = resMgmtClient.JoinChannel(setup.ChannelID); err != nil {
			return errors.WithMessage(err, "JoinChannel failed")
		}
	}

	return nil
}

// Initialize reads configuration from file and sets up client, channel and event hub
func (setup *BaseSetupImpl) Initialize() error {
	// Create SDK setup for the integration tests
	sdkOptions := deffab.Options{
		ConfigFile: setup.ConfigFile,
	}

	sdk, err := deffab.NewSDK(sdkOptions)
	if err != nil {
		return errors.WithMessage(err, "SDK init failed")
	}

	session, err := sdk.NewPreEnrolledUserSession(setup.OrgID, "Admin")
	if err != nil {
		return errors.WithMessage(err, "failed getting admin user session for org")
	}

	sc, err := sdk.NewSystemClient(session)
	if err != nil {
		return errors.WithMessage(err, "NewSystemClient failed")
	}

	setup.Client = sc
	setup.AdminUser = session.Identity()

	channel, err := setup.GetChannel(setup.Client, setup.ChannelID, []string{setup.OrgID})
	if err != nil {
		return errors.Wrapf(err, "create channel (%s) failed: %v", setup.ChannelID)
	}
	setup.Channel = channel

	//	if err = setup.ensureJoinChannel(sdk); err != nil {
	//		return err
	//	}

	if err := setup.setupEventHub(sc); err != nil {
		return err
	}

	setup.Initialized = true

	return nil
}

//2017-12-26 14:46:18,212 [8] DEBUG DigitalAssetSignerTool.MainForm - sign message is: create_account{"addr":"hc1BQMxicuFF7WVJHgRXnzi35ai7mSNLr3QB","ts":131587443770429007,"note":""}{"pubkey":"0463f80f1dfc879fc44eb84cca38220971c837457b94cf2023944476f3f135bc2f139b3eba475587f19c3fc38dff29905d5618c07c796e1fa3682052349bc96286"}
//2017-12-26 14:46:18,215 [8] DEBUG DigitalAssetSignerTool.MainForm - sign is: 1cfd8c7f775e5288e7aef2ce74a4a697310b9762aef077cc399dba219f7ceec4c24881b7e753e5c7b0d20caf6dc3a477078b9db51f06c5d052d139172b2e2fa670

//2017-12-26 15:19:49,208 [8] DEBUG DigitalAssetSignerTool.MainForm - sign message is: create_account{"addr":"hc148QruQaND3ALxbtrqKR3btDge9rxYadbf","ts":131587463879193870,"note":""}{"pubkey":"048f1c17d07d90fa3a4c9e245e1b9322bfe33a3206f65638d64b60c3a875b792a7a890977e93f5eb9c342f0e07ecf382c837377914df0406e44fd13ed22d9620a5"}
//2017-12-26 15:19:49,210 [8] DEBUG DigitalAssetSignerTool.MainForm - sign is: 1cffc62c9373d32ff4e20096634f7fdf2201a690d947296b0c164afb6847566ac45225f603641f9eaa433390ed59f3bc269aa7ddaceafae00545787e07cf2cff8d

func TestChannelClient(t *testing.T) {
	fmt.Printf("start test\n")

	testSetup := BaseSetupImpl{
		ConfigFile:      "./fixtures/config/" + "config_test.yaml",
		ChannelID:       "mychannel",
		OrgID:           "Org1",
		ChannelConfig:   path.Join("./fixtures/", "channel", "mychannel.tx"),
		ConnectEventHub: true,
		ChainCodeID:     "dasset14",
	}

	if err := testSetup.Initialize(); err != nil {
		t.Fatalf(err.Error())
	}

	// Create SDK setup for the integration tests
	sdkOptions := deffab.Options{
		ConfigFile: testSetup.ConfigFile,
		StateStoreOpts: opt.StateStoreOpts{
			Path: "/tmp/enroll_user",
		},
	}

	sdk, err := deffab.NewSDK(sdkOptions)
	if err != nil {
		t.Fatalf("Failed to create new SDK: %s", err)
	}

	chClient, err := sdk.NewChannelClient(testSetup.ChannelID, "User1")
	if err != nil {
		t.Fatalf("Failed to create new channel client: %s", err)
	}

	// Synchronous transaction
	invokeArgs := [][]byte{
		[]byte("create_account"),
		[]byte(`{"addr":"hc148QruQaND3ALxbtrqKR3btDge9rxYadbf","ts":131587463879193870,"note":""}`),
		[]byte(`{"pubkey":"048f1c17d07d90fa3a4c9e245e1b9322bfe33a3206f65638d64b60c3a875b792a7a890977e93f5eb9c342f0e07ecf382c837377914df0406e44fd13ed22d9620a5"}`),
		[]byte("1cffc62c9373d32ff4e20096634f7fdf2201a690d947296b0c164afb6847566ac45225f603641f9eaa433390ed59f3bc269aa7ddaceafae00545787e07cf2cff8d"),
	}
	request := apitxn.ExecuteTxRequest{
		ChaincodeID: testSetup.ChainCodeID,
		Fcn:         "invoke",
		Args:        invokeArgs,
	}
	result1, err := chClient.ExecuteTx(request)
	if err != nil {
		//t.Fatalf("Failed to create account: %s", err)
		fmt.Printf("Failed to create account: %s\n", err)
	}
	fmt.Printf("result1=%s\n", result1)

	// Synchronous query
	//testQuery("200", testSetup.ChainCodeID, chClient, t)
	var queryArgs = [][]byte{[]byte("query_account"), []byte("hc148QruQaND3ALxbtrqKR3btDge9rxYadbf")}
	result2, err := chClient.Query(
		apitxn.QueryRequest{ChaincodeID: testSetup.ChainCodeID, Fcn: "invoke", Args: queryArgs})
	if err != nil {
		t.Fatalf("Failed to invoke example cc: %s", err)
	}
	fmt.Printf("result2=%s\n", result2)

}
