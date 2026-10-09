package wallet

import (
	"encoding/hex"

	"github.com/btcsuite/btcutil/base58"
	"github.com/ethereum/go-ethereum/crypto"
)

type Service struct {
	encryptionKeys []string
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
