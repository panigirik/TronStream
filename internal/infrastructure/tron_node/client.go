package tron_node

import (
	"TronStream/internal/infrastructure/wallet"
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const (
	defaultPageLimit = 50
	maxPageLimit     = 200
	errorBodyLimit   = 4 << 10
)

type Client struct {
	baseURL string
	apiKey  string
	http    *http.Client
}

type Option func(*Client)

func WithAPIKey(key string) Option {
	return func(c *Client) { c.apiKey = key }
}

func NewClient(baseURL string, timeout time.Duration, opts ...Option) *Client {
	client := &Client{
		baseURL: trimTrailingSlash(baseURL),
		http:    &http.Client{Timeout: timeout},
	}

	for _, opt := range opts {
		opt(client)
	}

	return client
}

type TransfersQuery struct {
	AddressBase58 string
	MinTimestamp  int64
	Limit         int
	Fingerprint   string
	OnlyConfirmed bool
}

type Transfer struct {
	TxID           string
	BlockNumber    int64
	BlockTimestamp int64
	AmountSun      int64
	FromHex        string
	ToHex          string
}

type TransfersPage struct {
	Transfers   []Transfer
	RawCount    int
	Fingerprint string
}

func (c *Client) GetIncomingTransfers(ctx context.Context, q TransfersQuery) (TransfersPage, error) {
	requestURL := c.buildURL(q)

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, nil)
	if err != nil {
		return TransfersPage{}, fmt.Errorf("создание запроса: %w", err)
	}
	request.Header.Set("Accept", "application/json")
	if c.apiKey != "" {
		request.Header.Set("TRON-PRO-API-KEY", c.apiKey)
	}

	response, err := c.http.Do(request)
	if err != nil {
		return TransfersPage{}, fmt.Errorf("запрос к TronGrid: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(response.Body, errorBodyLimit))
		return TransfersPage{}, &APIError{StatusCode: response.StatusCode, Message: string(body)}
	}

	var parsed transactionsResponse
	if err := json.NewDecoder(response.Body).Decode(&parsed); err != nil {
		return TransfersPage{}, fmt.Errorf("разбор ответа TronGrid: %w", err)
	}

	if !parsed.Success {
		return TransfersPage{}, &APIError{StatusCode: response.StatusCode, Message: parsed.Error}
	}

	return TransfersPage{
		Transfers:   c.filterTransfers(parsed.Data, q.AddressBase58),
		RawCount:    len(parsed.Data),
		Fingerprint: parsed.Meta.Fingerprint,
	}, nil
}

func (c *Client) SignTransaction(tx *wallet.TronTransaction, privateKeyHex string) error {
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return err
	}

	txBytes, err := hex.DecodeString(tx.TxId)
	if err != nil {
		return err
	}

	signature, err := crypto.Sign(txBytes, privateKey)
	if err != nil {
		return err
	}

	tx.Signature = append(tx.Signature, hex.EncodeToString(signature))

	return nil
}

func (c *Client) BroadcastTransaction(tx *wallet.TronTransaction) (*BroadcastResponse, error) {
	jsonData, err := json.Marshal(tx)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, c.baseURL+"/wallet/broadcasttransaction", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	if c.apiKey != "" {
		req.Header.Set("TRON-PRO-API-KEY", c.apiKey)
	}
	response, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var broadcastResponse BroadcastResponse
	if err := json.Unmarshal(body, &broadcastResponse); err != nil {
		return nil, err
	}

	if !broadcastResponse.Result {
		return nil, &APIError{StatusCode: response.StatusCode, Message: broadcastResponse.Message}
	}

	return &broadcastResponse, nil
}

func (c *Client) CreateTransaction(ctx context.Context, fromAddress string, toAddress string, amountSun int64) (*wallet.TronTransaction, error) {
	requestBody := map[string]interface{}{
		"owner_address": fromAddress,
		"to_address":    toAddress,
		"amount":        amountSun,
		"visible":       true,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/wallet/createtransaction", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	if c.apiKey != "" {
		req.Header.Set("TRON-PRO-API-KEY", c.apiKey)
	}
	response, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		body, _ := io.ReadAll(io.LimitReader(response.Body, errorBodyLimit))
		return nil, &APIError{StatusCode: response.StatusCode, Message: string(body)}
	}

	var transaction wallet.TronTransaction
	if err := json.NewDecoder(response.Body).Decode(&transaction); err != nil {
		return nil, err
	}
	if transaction.TxId == "" {
		return nil, &APIError{StatusCode: response.StatusCode, Message: "Tron node returned an empty transaction"}
	}

	return &transaction, nil
}

func (c *Client) buildURL(q TransfersQuery) string {
	limit := q.Limit
	if limit <= 0 {
		limit = defaultPageLimit
	}
	if limit > maxPageLimit {
		limit = maxPageLimit
	}

	params := url.Values{}
	params.Set("only_to", "true")
	params.Set("only_confirmed", strconv.FormatBool(q.OnlyConfirmed))
	params.Set("limit", strconv.Itoa(limit))
	params.Set("order_by", "block_timestamp,asc")
	if q.MinTimestamp > 0 {
		params.Set("min_timestamp", strconv.FormatInt(q.MinTimestamp, 10))
	}
	if q.Fingerprint != "" {
		params.Set("fingerprint", q.Fingerprint)
	}

	return fmt.Sprintf("%s/v1/accounts/%s/transactions?%s",
		c.baseURL, url.PathEscape(q.AddressBase58), params.Encode())
}

func (c *Client) filterTransfers(transactions []transactionDTO, addressBase58 string) []Transfer {
	var transfers []Transfer

	for _, transaction := range transactions {
		if !transaction.succeeded() {
			continue
		}

		for _, contract := range transaction.RawData.Contract {
			if contract.Type != transferContractType {
				continue
			}

			value := contract.Parameter.Value
			if value.Amount <= 0 {
				continue
			}

			if !wallet.EqualAddress(value.ToAddress, addressBase58) {
				continue
			}

			transfers = append(transfers, Transfer{
				TxID:           transaction.TxID,
				BlockNumber:    transaction.BlockNumber,
				BlockTimestamp: transaction.BlockTimestamp,
				AmountSun:      value.Amount,
				FromHex:        value.OwnerAddress,
				ToHex:          value.ToAddress,
			})
		}
	}

	return transfers
}

func trimTrailingSlash(s string) string {
	for len(s) > 0 && s[len(s)-1] == '/' {
		s = s[:len(s)-1]
	}

	return s
}
