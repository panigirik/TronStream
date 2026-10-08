package tron_node

// Структуры разбора ответа GET /v1/accounts/{address}/transactions.

const transferContractType = "TransferContract"

const contractResultSuccess = "SUCCESS"

type transactionsResponse struct {
	Data    []transactionDTO `json:"data"`
	Success bool             `json:"success"`
	Error   string           `json:"error"`
	Meta    metaDTO          `json:"meta"`
}

type metaDTO struct {
	Fingerprint string `json:"fingerprint"`
}

type transactionDTO struct {
	TxID           string     `json:"txID"`
	BlockNumber    int64      `json:"blockNumber"`
	BlockTimestamp int64      `json:"block_timestamp"`
	Ret            []retDTO   `json:"ret"`
	RawData        rawDataDTO `json:"raw_data"`
}

type retDTO struct {
	ContractRet string `json:"contractRet"`
}

type rawDataDTO struct {
	Contract []contractDTO `json:"contract"`
}

type contractDTO struct {
	Type      string       `json:"type"`
	Parameter parameterDTO `json:"parameter"`
}

type parameterDTO struct {
	Value valueDTO `json:"value"`
}

type valueDTO struct {
	Amount       int64  `json:"amount"`
	OwnerAddress string `json:"owner_address"`
	ToAddress    string `json:"to_address"`
}

// succeeded сообщает, что транзакция реально исполнилась.
// Неуспешные транзакции TronGrid тоже возвращает в общей выдаче.
func (t transactionDTO) succeeded() bool {
	if len(t.Ret) == 0 {
		return false
	}

	return t.Ret[0].ContractRet == contractResultSuccess
}

type BroadcastResponse struct {
	Result  bool   `json:"result"`
	TxID    string `json:"txid,omitempty"`
	Code    string `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}
