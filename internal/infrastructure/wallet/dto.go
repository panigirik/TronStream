package wallet

type TronTransaction struct {
	TxId       string   `json:"txID"`
	RawData    RawData  `json:"raw_data"`
	RawDataHEx string   `json:"raw_data_hex"`
	Signature  []string `json:"signature,omitempty"`
}

type RawData struct {
	Contract      []Contract `json:"contract"`
	RefBlockBytes string     `json:"ref_block_bytes"`
	RefBlockHash  string     `json:"ref_block_hash"`
	Expiration    int        `json:"expiration"`
	Timestamp     int        `json:"timestamp"`
	FeeLimit      int        `json:"fee_limit,omitempty"`
}

type Contract struct {
	Parameter Parameter `json:"parameter"`
	Type      string    `json:"type"`
}

type Parameter struct {
	Value   map[string]interface{} `json:"value"`
	TypeUrl string                 `json:"type_url"`
}

type BroadCastResponse struct {
	Result  bool   `json:"result"`
	TxID    string `json:"txId"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}
