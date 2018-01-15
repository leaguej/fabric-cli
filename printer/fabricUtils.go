package printer

import (
	"fmt"

	"github.com/golang/protobuf/proto"

	cb "github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/protos/common"
	pb "github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/protos/peer"

	"github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/protos/ledger/rwset"
	"github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/protos/ledger/rwset/kvrwset"
)

// ExtractPayloadOrPanic retrieves the payload of a given envelope and unmarshals it -- it panics if either of these operations fail.
func ExtractPayloadOrPanic(envelope *cb.Envelope) *cb.Payload {
	payload, err := ExtractPayload(envelope)
	if err != nil {
		panic(err)
	}
	return payload
}

// ExtractPayload retrieves the payload of a given envelope and unmarshals it.
func ExtractPayload(envelope *cb.Envelope) (*cb.Payload, error) {
	payload := &cb.Payload{}
	if err := proto.Unmarshal(envelope.Payload, payload); err != nil {
		return nil, fmt.Errorf("Envelope does not carry a Payload: %s", err)
	}
	return payload, nil
}

// ExtractEnvelopeOrPanic retrieves the requested envelope from a given block and unmarshals it -- it panics if either of these operation fail.
func ExtractEnvelopeOrPanic(block *cb.Block, index int) *cb.Envelope {
	envelope, err := ExtractEnvelope(block, index)
	if err != nil {
		panic(err)
	}
	return envelope
}

// ExtractEnvelope retrieves the requested envelope from a given block and unmarshals it.
func ExtractEnvelope(block *cb.Block, index int) (*cb.Envelope, error) {
	if block.Data == nil {
		return nil, fmt.Errorf("No data in block")
	}

	envelopeCount := len(block.Data.Data)
	if index < 0 || index >= envelopeCount {
		return nil, fmt.Errorf("Envelope index out of bounds")
	}
	marshaledEnvelope := block.Data.Data[index]
	envelope, err := GetEnvelopeFromBlock(marshaledEnvelope)
	if err != nil {
		return nil, fmt.Errorf("Block data does not carry an envelope at index %d: %s", index, err)
	}
	return envelope, nil
}

// GetEnvelopeFromBlock gets an envelope from a block's Data field.
func GetEnvelopeFromBlock(data []byte) (*cb.Envelope, error) {
	//Block always begins with an envelope
	var err error
	env := &cb.Envelope{}
	if err = proto.Unmarshal(data, env); err != nil {
		return nil, fmt.Errorf("Error getting envelope(%s)", err)
	}

	return env, nil
}

// UnmarshalChannelHeader returns a ChannelHeader from bytes
func UnmarshalChannelHeader(bytes []byte) (*cb.ChannelHeader, error) {
	chdr := &cb.ChannelHeader{}
	err := proto.Unmarshal(bytes, chdr)
	if err != nil {
		return nil, fmt.Errorf("UnmarshalChannelHeader failed, err %s", err)
	}

	return chdr, nil
}

// GetSignatureHeader Get SignatureHeader from bytes
func GetSignatureHeader(bytes []byte) (*cb.SignatureHeader, error) {
	sh := &cb.SignatureHeader{}
	err := proto.Unmarshal(bytes, sh)
	return sh, err
}

// GetTransaction Get Transaction from bytes
func GetTransaction(txBytes []byte) (*pb.Transaction, error) {
	tx := &pb.Transaction{}
	err := proto.Unmarshal(txBytes, tx)
	return tx, err
}

// UnmarshalConfigUpdate attempts to unmarshal bytes to a *cb.ConfigUpdate
func UnmarshalConfigUpdate(data []byte) (*cb.ConfigUpdate, error) {
	configUpdate := &cb.ConfigUpdate{}
	err := proto.Unmarshal(data, configUpdate)
	if err != nil {
		return nil, err
	}
	return configUpdate, nil
}

// UnmarshalConfigUpdate attempts to unmarshal bytes to a *cb.ConfigUpdate or panics
func UnmarshalConfigUpdateOrPanic(data []byte) *cb.ConfigUpdate {
	result, err := UnmarshalConfigUpdate(data)
	if err != nil {
		panic(err)
	}
	return result
}

// GetChaincodeActionPayload Get ChaincodeActionPayload from bytes
func GetChaincodeActionPayload(capBytes []byte) (*pb.ChaincodeActionPayload, error) {
	cap := &pb.ChaincodeActionPayload{}
	err := proto.Unmarshal(capBytes, cap)
	return cap, err
}

// TxValidationFlags is array of transaction validation codes. It is used when commiter validates block.
type TxValidationFlags []uint8

// Flag returns validation code at specified transaction
func (obj TxValidationFlags) Flag(txIndex int) pb.TxValidationCode {
	return pb.TxValidationCode(obj[txIndex])
}

// ProviderType indicates the type of an identity provider
type ProviderType int

// The ProviderType of a member relative to the member API
const (
	FABRIC ProviderType = iota // MSP is of FABRIC type
	OTHER                      // MSP is of OTHER TYPE
)

//-------- ChaincodeData is stored on the LSCC -------

//ChaincodeData defines the datastructure for chaincodes to be serialized by proto
//Type provides an additional check by directing to use a specific package after instantiation
//Data is Type specifc (see CDSPackage and SignedCDSPackage)
type ChaincodeData struct {
	//Name of the chaincode
	Name string `protobuf:"bytes,1,opt,name=name"`

	//Version of the chaincode
	Version string `protobuf:"bytes,2,opt,name=version"`

	//Escc for the chaincode instance
	Escc string `protobuf:"bytes,3,opt,name=escc"`

	//Vscc for the chaincode instance
	Vscc string `protobuf:"bytes,4,opt,name=vscc"`

	//Policy endorsement policy for the chaincode instance
	Policy []byte `protobuf:"bytes,5,opt,name=policy,proto3"`

	//Data data specific to the package
	Data []byte `protobuf:"bytes,6,opt,name=data,proto3"`

	//Id of the chaincode that's the unique fingerprint for the CC
	//This is not currently used anywhere but serves as a good
	//eyecatcher
	Id []byte `protobuf:"bytes,7,opt,name=id,proto3"`

	//InstantiationPolicy for the chaincode
	InstantiationPolicy []byte `protobuf:"bytes,8,opt,name=instantiation_policy,proto3"`
}

//----- CDSData ------

//CDSData is data stored in the LSCC on instantiation of a CC
//for CDSPackage.  This needs to be serialized for ChaincodeData
//hence the protobuf format
type CDSData struct {
	//CodeHash hash of CodePackage from ChaincodeDeploymentSpec
	CodeHash []byte `protobuf:"bytes,1,opt,name=codehash,proto3"`

	//MetaDataHash hash of Name and Version from ChaincodeDeploymentSpec
	MetaDataHash []byte `protobuf:"bytes,2,opt,name=metadatahash,proto3"`
}

//Reset resets
func (data *CDSData) Reset() { *data = CDSData{} }

//String converts to string
func (data *CDSData) String() string { return proto.CompactTextString(data) }

//ProtoMessage just exists to make proto happy
func (*CDSData) ProtoMessage() {}

// TxRwSet acts as a proxy of 'rwset.TxReadWriteSet' proto message and helps constructing Read-write set specifically for KV data model
type TxRwSet struct {
	NsRwSets []*NsRwSet
}

// NsRwSet encapsulates 'kvrwset.KVRWSet' proto message for a specific name space (chaincode)
type NsRwSet struct {
	NameSpace string
	KvRwSet   *kvrwset.KVRWSet
}

// FromProtoBytes deserializes protobytes into TxReadWriteSet proto message and populates 'TxRwSet'
func (txRwSet *TxRwSet) FromProtoBytes(protoBytes []byte) error {
	protoTxRwSet := &rwset.TxReadWriteSet{}
	if err := proto.Unmarshal(protoBytes, protoTxRwSet); err != nil {
		return err
	}
	protoNsRwSets := protoTxRwSet.GetNsRwset()
	var nsRwSet *NsRwSet
	for _, protoNsRwSet := range protoNsRwSets {
		nsRwSet = &NsRwSet{}
		nsRwSet.NameSpace = protoNsRwSet.Namespace
		protoRwSetBytes := protoNsRwSet.Rwset

		protoKvRwSet := &kvrwset.KVRWSet{}
		if err := proto.Unmarshal(protoRwSetBytes, protoKvRwSet); err != nil {
			return err
		}
		nsRwSet.KvRwSet = protoKvRwSet
		txRwSet.NsRwSets = append(txRwSet.NsRwSets, nsRwSet)
	}
	return nil
}
