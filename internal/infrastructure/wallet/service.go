package wallet

import (
	"TronStream/internal/database"
	"TronStream/internal/infrastructure/tron_node"
	"context"
	"encoding/hex"
	"github.com/btcsuite/btcutil/base58"
	"github.com/ethereum/go-ethereum/crypto"

	"TronStream/internal/config"
)

type Service struct {
	encryptionKeys   []string
	walletRepository database.WalletRepository
	config           *config.Config
	client           *tron_node.Client
}

type GeneratedWallet struct {
	PublicKey  string
	PrivateKey string
	Address    string
}

func NewWalletService(encryptionKeys []string) *Service {
	return &Service{encryptionKeys: encryptionKeys}
}

func (w *Service) GenerateKey() (GeneratedWallet, error) {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		return GeneratedWallet{}, err
	}

	privateKeyBytes := crypto.FromECDSA(privateKey)
	publicKeyBytes := crypto.FromECDSAPub(&privateKey.PublicKey)
	publicKeyBytes = publicKeyBytes[1:]

	hash := crypto.Keccak256(publicKeyBytes)

	addressBytes := hash[len(hash)-20:]

	address := base58.CheckEncode(addressBytes, 0x41)

	return GeneratedWallet{
		PublicKey:  hex.EncodeToString(publicKeyBytes),
		PrivateKey: hex.EncodeToString(privateKeyBytes),
		Address:    address,
	}, nil
}

func (w *Service) TransferCrypto(ctx context.Context, user_id int64, amount float64, toAddress string) error {
	wallet, err := w.walletRepository.GetWalletByUserId(ctx, user_id)
	if err != nil {
		return err
	}

	rawTx := TronTransaction{
		TxId:       "cbd66c888d3e64c51475c7b39868e4de8a7199bc4f3d2fddde491efc77d2ef1c", // Пример TxID (хеш raw_data)
		RawDataHEx: "0a022408...",                                                      // Hex-представление сырых данных
		RawData: RawData{
			RefBlockBytes: "2408",
			RefBlockHash:  "63a2...",
			Expiration:    1712345678000,
			Timestamp:     1712342078000,
			Contract: []Contract{
				{
					Type: "TransferContract",
					Parameter: Parameter{
						TypeUrl: "://googleapis.com",
						Value: map[string]interface{}{
							"amount":        amount / 1_000_000,
							"owner_address": wallet.Address,
							"to_address":    toAddress,
						},
					},
				},
			},
		},
	}
	//TODO доправить до конца
	w.client.SignTransaction(&rawTx, wallet.PrivateKeyEncrypted)
	w.client.BroadcastTransaction(&rawTx)
	return nil
}
