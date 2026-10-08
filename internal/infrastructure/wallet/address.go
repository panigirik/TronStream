package wallet

import (
	"encoding/hex"
	"errors"
	"strings"

	"github.com/btcsuite/btcutil/base58"
)

const TronAddressPrefix byte = 0x41

const addressBodyLen = 20

var ErrInvalidTronAddress = errors.New("wallet: invalid TRON address")

func Base58ToHex(address string) (string, error) {
	decoded, version, err := base58.CheckDecode(address)
	if err != nil {
		return "", ErrInvalidTronAddress
	}
	if version != TronAddressPrefix {
		return "", ErrInvalidTronAddress
	}
	if len(decoded) != addressBodyLen {
		return "", ErrInvalidTronAddress
	}

	return hex.EncodeToString(append([]byte{version}, decoded...)), nil
}

func HexToBase58(hexAddress string) (string, error) {
	raw, err := hex.DecodeString(strings.ToLower(hexAddress))
	if err != nil {
		return "", ErrInvalidTronAddress
	}
	if len(raw) != addressBodyLen+1 {
		return "", ErrInvalidTronAddress
	}
	if raw[0] != TronAddressPrefix {
		return "", ErrInvalidTronAddress
	}

	return base58.CheckEncode(raw[1:], raw[0]), nil
}

func EqualAddress(a, b string) bool {
	normalizedA, err := normalizeAddress(a)
	if err != nil {
		return false
	}
	normalizedB, err := normalizeAddress(b)
	if err != nil {
		return false
	}

	return normalizedA == normalizedB
}

func normalizeAddress(address string) (string, error) {
	if strings.HasPrefix(address, "T") {
		return Base58ToHex(address)
	}

	lowered := strings.ToLower(address)
	if _, err := HexToBase58(lowered); err != nil {
		return "", err
	}

	return lowered, nil
}
