package sdk

import (
	"encoding/hex"
	"encoding/json"
	//"errors"
	"fmt"
	//"time"

	"github.com/hyperledger/fabric-sdk-go/api/apitxn"
)

type InvokeReturnData struct {
	TxID  string `json:"txid"`
	Nonce string `json:"nonce"`
}

func (testSetup *BaseSetupImpl) InvokeChainCode(method, content, sign string) (string, error) {

	//	eventID := "test([a-zA-Z]+)"

	//	// Register chaincode event (pass in channel which receives event details when the event is complete)
	//	notifier := make(chan *apitxn.CCEvent)
	//	rce := testSetup.ChannelClient.RegisterChaincodeEvent(notifier, testSetup.ChainCodeID, eventID)

	invokeArgs := [][]byte{
		[]byte(method),
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
		//testSetup.ChannelClient.UnregisterChaincodeEvent(rce)
		return "", err
	}

	//	select {
	//	case ccEvent := <-notifier:
	//		fmt.Printf("Received cc event: %s", ccEvent)
	//		if ccEvent.TxID != txID.ID {
	//			msg := fmt.Sprintf("CCEvent(%s) and ExecuteTx(%s) transaction IDs don't match", ccEvent.TxID, txID.ID)
	//			err = errors.New(msg)
	//			fmt.Println(msg)
	//			testSetup.ChannelClient.UnregisterChaincodeEvent(rce)
	//			return "", err
	//		}
	//	case <-time.After(time.Second * 20):
	//		msg := fmt.Sprintf("Did NOT receive CC for eventId(%s)\n", eventID)
	//		err = errors.New(msg)
	//		fmt.Println(msg)
	//		testSetup.ChannelClient.UnregisterChaincodeEvent(rce)
	//		return "", err
	//	}

	//	// Unregister chain code event using registration handle
	//	err = testSetup.ChannelClient.UnregisterChaincodeEvent(rce)
	//	if err != nil {
	//		fmt.Printf("Unregister cc event failed: %s", err)
	//		return "", err
	//	}

	data := &InvokeReturnData{
		TxID:  txID.ID,
		Nonce: hex.EncodeToString(txID.Nonce),
	}
	jsonResult, _ := json.Marshal(data)

	//jsonResult := fmt.Sprintf(`{"txid":"%s", nonce="%s"}`, txID.ID, hex.EncodeToString(txID.Nonce))
	fmt.Printf("method=%s, result1=%s\n", method, jsonResult)
	return string(jsonResult), nil
}

func (testSetup *BaseSetupImpl) QueryChainCode(method, filter string) (string, error) {

	invokeArgs := [][]byte{
		[]byte(method),
		[]byte(filter),
	}
	request := apitxn.QueryRequest{
		ChaincodeID: testSetup.ChainCodeID,
		Fcn:         "query",
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

func (testSetup *BaseSetupImpl) QueryChainCodeWithParams(params [][]byte) (string, error) {

	//	invokeArgs := [][]byte{
	//		[]byte(method),
	//		[]byte(filter),
	//	}
	request := apitxn.QueryRequest{
		ChaincodeID: testSetup.ChainCodeID,
		Fcn:         "query",
		Args:        params,
	}
	result1, err := testSetup.ChannelClient.Query(request)
	if err != nil {
		//fmt.Printf("Failed to create account: %s\n", err)
		return "", err
	}
	fmt.Printf("result1=%s\n", result1)
	return string(result1), nil
}
