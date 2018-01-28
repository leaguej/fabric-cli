package sdk

type TxSimpleData struct {
	TxID     string `json:"TxId"`
	Request  string `json:"Request"`
	Response struct {
		Message string `json:"Message"`
		Status  string `json:"Status"`
		Payload string `json:"Payload"`
	} `json:"Response"`
}

type TransactionData struct {
	Signature string `json:"Signature"`
	Payload   struct {
		Header struct {
			ChannelHeader struct {
				Type      string `json:"Type"`
				ChannelID string `json:"ChannelId"`
				Epoch     string `json:"Epoch"`
				Extension struct {
					ChaincodeID struct {
						Name    string `json:"Name"`
						Version string `json:"Version"`
						Path    string `json:"Path"`
					} `json:"ChaincodeId"`
					PayloadVisibility string `json:"PayloadVisibility"`
				} `json:"Extension"`
				Timestamp string `json:"Timestamp"`
				TxID      string `json:"TxId"`
				Version   string `json:"Version"`
			} `json:"ChannelHeader"`
			SignatureHeader struct {
				Nonce   string `json:"Nonce"`
				Creator string `json:"Creator"`
			} `json:"SignatureHeader"`
		} `json:"Header"`
		Data struct {
			Type    string `json:"Type"`
			Actions []struct {
				Header struct {
					Nonce   string `json:"Nonce"`
					Creator string `json:"Creator"`
				} `json:"Header"`
				Payload struct {
					Action struct {
						Endorsements []struct {
							Endorser  string `json:"Endorser"`
							Signature string `json:"Signature"`
						} `json:"Endorsements"`
						ProposalResponsePayload struct {
							ProposalHash string `json:"ProposalHash"`
							Extension    struct {
								Response struct {
									Message string `json:"Message"`
									Status  string `json:"Status"`
									Payload string `json:"Payload"`
								} `json:"Response"`
								Results struct {
									NsRWs []struct {
										NameSpace string `json:"NameSpace"`
										KvRwSet   struct {
											Reads []struct {
												Key     string `json:"Key"`
												Version struct {
													BlockNum string `json:"BlockNum"`
													TxNum    string `json:"TxNum"`
												} `json:"Version"`
											} `json:"Reads"`
											Writes []struct {
												Key      string `json:"Key"`
												IsDelete string `json:"IsDelete"`
												Value    string `json:"Value"`
											} `json:"Writes"`
											RangeQueriesInfo []interface{} `json:"RangeQueriesInfo"`
										} `json:"KvRwSet"`
									} `json:"NsRWs"`
								} `json:"Results"`
								Events struct {
								} `json:"Events"`
							} `json:"Extension"`
						} `json:"ProposalResponsePayload"`
					} `json:"Action"`
				} `json:"Payload"`
			} `json:"Actions"`
		} `json:"Data"`
	} `json:"Payload"`
}

type BlockData struct {
	Header struct {
		Number       string `json:"Number"`
		PreviousHash string `json:"PreviousHash"`
		DataHash     string `json:"DataHash"`
	} `json:"Header"`
	Metadata struct {
		Signatures []struct {
			Nonce   string `json:"Nonce"`
			Creator string `json:"Creator"`
		} `json:"Signatures"`
		LastConfigIndex    string   `json:"Last Config Index"`
		TransactionFilters []string `json:"TransactionFilters"`
		OrdererMetadata    string   `json:"Orderer Metadata"`
	} `json:"Metadata"`
	Data struct {
		Data []struct {
			Signature string `json:"Signature"`
			Payload   struct {
				Header struct {
					ChannelHeader struct {
						Type      string `json:"Type"`
						ChannelID string `json:"ChannelId"`
						Epoch     string `json:"Epoch"`
						Extension struct {
							ChaincodeID struct {
								Name    string `json:"Name"`
								Version string `json:"Version"`
								Path    string `json:"Path"`
							} `json:"ChaincodeId"`
							PayloadVisibility string `json:"PayloadVisibility"`
						} `json:"Extension"`
						Timestamp string `json:"Timestamp"`
						TxID      string `json:"TxId"`
						Version   string `json:"Version"`
					} `json:"ChannelHeader"`
					SignatureHeader struct {
						Nonce   string `json:"Nonce"`
						Creator string `json:"Creator"`
					} `json:"SignatureHeader"`
				} `json:"Header"`
				Data struct {
					Type    string `json:"Type"`
					Actions []struct {
						Header struct {
							Nonce   string `json:"Nonce"`
							Creator string `json:"Creator"`
						} `json:"Header"`
						Payload struct {
							Action struct {
								Endorsements []struct {
									Endorser  string `json:"Endorser"`
									Signature string `json:"Signature"`
								} `json:"Endorsements"`
								ProposalResponsePayload struct {
									ProposalHash string `json:"ProposalHash"`
									Extension    struct {
										Response struct {
											Message string `json:"Message"`
											Status  string `json:"Status"`
											Payload string `json:"Payload"`
										} `json:"Response"`
										Results struct {
											NsRWs []struct {
												NameSpace string `json:"NameSpace"`
												KvRwSet   struct {
													Reads []struct {
														Key     string `json:"Key"`
														Version struct {
														} `json:"Version"`
													} `json:"Reads"`
													Writes []struct {
														Key      string `json:"Key"`
														IsDelete string `json:"IsDelete"`
														Value    string `json:"Value"`
													} `json:"Writes"`
													RangeQueriesInfo []interface{} `json:"RangeQueriesInfo"`
												} `json:"KvRwSet"`
											} `json:"NsRWs"`
										} `json:"Results"`
										Events struct {
										} `json:"Events"`
									} `json:"Extension"`
								} `json:"ProposalResponsePayload"`
							} `json:"Action"`
						} `json:"Payload"`
					} `json:"Actions"`
				} `json:"Data"`
			} `json:"Payload"`
		} `json:"Data"`
	} `json:"Data"`
}
