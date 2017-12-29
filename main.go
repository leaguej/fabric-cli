package main

import (
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/leaguej/fabric-cli/chaincode"
	"github.com/leaguej/fabric-cli/common"
	cliConfig "github.com/leaguej/fabric-cli/config"
)

func CreateAccount() error {
	var flags = &pflag.FlagSet{}
	cliConfig.InitChannelID(flags, CHANNEL_DIGITAL_ASSET)
	cliConfig.InitChaincodeID(flags, CHAINCODE_DIGITAL_ASSET)
	cliConfig.InitConfigFile(flags)
	cliConfig.InitLoggingLevel(flags)
	cliConfig.InitUserName(flags)
	cliConfig.InitUserPassword(flags)
	cliConfig.InitOrdererTLSCertificate(flags)
	//cliConfig.InitPrintFormat(flags, "json")
	cliConfig.InitWriter(flags)
	cliConfig.InitOrgIDs(flags)

	cliConfig.InitIterations(flags)
	cliConfig.InitSleepTime(flags)
	cliConfig.InitTimeout(flags)

	action, err := chaincode.NewInvokeAction(&pflag.FlagSet{})
	if err != nil {
		return err
	}

	defer action.Terminate()

	curve := btcec.S256()
	priKey, _ := btcec.NewPrivateKey(curve)
	pubKey := priKey.PubKey()

	addr, _ := PublicKeyToAddress(pubKey)

	method := "create_account"
	header := fmt.Sprintf(`{"addr":"%s","ts":%d,"note":"test create account"}`,
		addr, 123456789)
	content := fmt.Sprintf(`{"pubkey":"%s"}`,
		hex.EncodeToString(pubKey.SerializeUncompressed()))

	message := method + header + content
	hashed := sha256.Sum256([]byte(message))
	sign_data, _ := btcec.SignCompact(curve, priKey, hashed[:], false)
	sign := hex.EncodeToString(sign_data)

	//	beego.Info("private key is: " + hex.EncodeToString(priKey.Serialize()))
	//	beego.Info("public key is: " + hex.EncodeToString(pubKey.SerializeUncompressed()))
	//	beego.Info("address is: " + addr)
	//	beego.Info("header is: " + header)
	//	beego.Info("content is: " + content)
	//	beego.Info("message is: " + message)
	//	beego.Info("sign is: " + sign)

	args := &common.ArgStruct{
		Func: "invoke",
		Args: []string{method, header, content, sign},
	}

	err = action.Invoke(args)
	if err != nil {
		return err
	}

	message = "private key: " + hex.EncodeToString(priKey.Serialize()) +
		", addr: " + addr +
		", public key: " + hex.EncodeToString(pubKey.SerializeUncompressed())

	has_error := false
	if action.Response_Status < 0 {
		has_error = true
	}

	fmt.Printf("error=%s, message=%s, payload=%s\n", has_error,
		action.Response_Message+"\n"+message, action.Response_Payload)

	return nil
}

func main() {
	CreateAccount()
}
