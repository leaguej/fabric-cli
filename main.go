package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/leaguej/fabric-cli/btcd/btcec"
	"github.com/leaguej/fabric-cli/digitalasset/util"
	"github.com/leaguej/fabric-cli/sdk"
)

func CreateAccount() error {

	setup, err := sdk.NewSdkClient()
	if err != nil {
		return err
	}

	defer setup.Close()

	curve := btcec.S256()
	priKey, _ := btcec.NewPrivateKey(curve)
	pubKey := priKey.PubKey()

	addr, _ := util.PublicKeyToAddress(pubKey.SerializeUncompressed())

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

	result, err := setup.InvokeChainCode(method, header, content, sign)
	if err != nil {
		return err
	}

	message = "private key: " + hex.EncodeToString(priKey.Serialize()) +
		", addr: " + addr +
		", public key: " + hex.EncodeToString(pubKey.SerializeUncompressed())

	fmt.Printf("message=%s, payload=%s\n",
		result+"\n"+message, result)

	return nil
}

func main() {
	CreateAccount()
}
