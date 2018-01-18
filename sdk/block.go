package sdk

import (
	"encoding/base64"
	"encoding/json"
	"strconv"
	"strings"

	"github.com/leaguej/fabric-cli/printer"

	fabricCommon "github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/protos/common"
)

// Base64URLDecode decodes the base64 string into a byte array
func Base64URLDecode(data string) ([]byte, error) {
	//check if it has padding or not
	if strings.HasSuffix(data, "=") {
		return base64.URLEncoding.DecodeString(data)
	}
	return base64.RawURLEncoding.DecodeString(data)
}

func (testSetup *BaseSetupImpl) QueryBlock(blockID string, bHash bool) (string, error) {
	var block *fabricCommon.Block

	if !bHash {
		num, err := strconv.Atoi(blockID)
		if err != nil {
			num = -1
		}
		block, err = testSetup.Channel.QueryBlock(num)
		if err != nil {
			return "", err
		}
	} else {
		hashBytes, err := Base64URLDecode(blockID)
		if err != nil {
			return "", err
		}

		block, err = testSetup.Channel.QueryBlockByHash(hashBytes)
		if err != nil {
			return "", err
		}
	}

	//	data, err := json.Marshal(block)
	//	if err != nil {
	//		return "", err
	//	}

	//	return string(data), nil
	p := printer.NewBlockPrinter(printer.JSON, printer.BUFFER)
	p.PrintBlock(block)

	str, err := p.ToString()

	return str, err
}

func (testSetup *BaseSetupImpl) QueryTx(txID string) (string, string, error) {
	tx, err := testSetup.Channel.QueryTransaction(txID)
	if err != nil {
		return "", "", err
	}

	//	data, err := json.Marshal(tx)
	//	if err != nil {
	//		return "", err
	//	}
	p := printer.NewBlockPrinter(printer.JSON, printer.BUFFER)
	p.PrintProcessedTransaction(tx)

	str, err := p.ToString()

	txData := &TransactionData{}
	err = json.Unmarshal([]byte(str), txData)
	if err != nil {
		return "", "", err
	}

	simpleData := &TxSimpleData{}
	simpleData.TxID = txData.Payload.Header.ChannelHeader.TxID
	resp := &txData.Payload.Data.Actions[0].Payload.Action.ProposalResponsePayload.Extension.Response
	simpleData.Response.Message = resp.Message
	simpleData.Response.Payload = resp.Payload
	simpleData.Response.Status = resp.Status

	nsRWs := txData.Payload.Data.Actions[0].Payload.Action.ProposalResponsePayload.Extension.Results.NsRWs
	for _, nsRW := range nsRWs {
		for _, write := range nsRW.KvRwSet.Writes {
			if strings.HasPrefix(write.Key, "DigitalAsset_last_request_") {
				simpleData.Request = write.Value
				break
			}
		}
	}

	simpleStr, _ := json.Marshal(simpleData)

	return str, string(simpleStr), err
}
