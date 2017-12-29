package sdk

import (
	"encoding/hex"
	"fmt"

	"github.com/hyperledger/fabric-sdk-go/api/apitxn"
)

func (testSetup *BaseSetupImpl) InvokeChainCode(method, header, content, sign string) (string, error) {

	invokeArgs := [][]byte{
		[]byte(method),
		[]byte(header),
		[]byte(content),
		[]byte(sign),
	}
	request := apitxn.ExecuteTxRequest{
		ChaincodeID: testSetup.ChainCodeID,
		Fcn:         "invoke",
		Args:        invokeArgs,
	}
	txID, err := testSetup.ChannelClient.ExecuteTx(request)
	if err != nil {
		fmt.Printf("Failed to create account: %s\n", err)
		return "", err
	}

	jsonResult := fmt.Sprintf(`{"txid":"%s", nonce="%s"}`, txID.ID, hex.EncodeToString(txID.Nonce))
	fmt.Printf("result1=%s\n", jsonResult)
	return jsonResult, nil
}

func (testSetup *BaseSetupImpl) QueryChainCode(method, filter string) (string, error) {

	invokeArgs := [][]byte{
		[]byte(method),
		[]byte(filter),
	}
	request := apitxn.QueryRequest{
		ChaincodeID: testSetup.ChainCodeID,
		Fcn:         "invoke",
		Args:        invokeArgs,
	}
	result1, err := testSetup.ChannelClient.Query(request)
	if err != nil {
		//fmt.Printf("Failed to create account: %s\n", err)
		return "", err
	}
	fmt.Printf("result1=%s\n", result1)
	return string(result1), nil
}
