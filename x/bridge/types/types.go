package types

import (
	"math/big"

	"cosmossdk.io/errors"
	"github.com/cosmos/cosmos-sdk/bsc/rlp"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	TransferOutChannel = "transferOut"
	TransferInChannel  = "transferIn"

	TransferOutChannelID sdk.ChannelID = 1
	TransferInChannelID  sdk.ChannelID = 2
)

type RefundReason uint32

const (
	RefundReasonInsufficientBalance RefundReason = 1
	RefundReasonFailAck             RefundReason = 2
	Unknown                         RefundReason = 3
)

type TransferOutSynPackage struct {
	Amount        *big.Int
	Recipient     sdk.AccAddress
	RefundAddress sdk.AccAddress
}

func DeserializeTransferOutSynPackage(serializedPackage []byte) (*TransferOutSynPackage, error) {
	var tp TransferOutSynPackage
	err := rlp.DecodeBytes(serializedPackage, &tp)
	if err != nil {
		return nil, errors.Wrapf(ErrInvalidPackage, "deserialize transfer out package failed")
	}
	return &tp, nil
}

type TransferOutRefundPackage struct {
	RefundAmount *big.Int
	RefundAddr   sdk.AccAddress
	RefundReason RefundReason
}

func DeserializeTransferOutRefundPackage(serializedPackage []byte) (*TransferOutRefundPackage, error) {
	var tp TransferOutRefundPackage
	err := rlp.DecodeBytes(serializedPackage, &tp)
	if err != nil {
		return nil, errors.Wrapf(ErrInvalidPackage, "deserialize transfer out refund package failed")
	}
	return &tp, nil
}

type TransferInSynPackage struct {
	Amount          *big.Int
	ReceiverAddress sdk.AccAddress
	RefundAddress   sdk.AccAddress
}

func DeserializeTransferInSynPackage(serializedPackage []byte) (*TransferInSynPackage, error) {
	var tp TransferInSynPackage
	err := rlp.DecodeBytes(serializedPackage, &tp)
	if err != nil {
		return nil, errors.Wrapf(ErrInvalidPackage, "deserialize transfer in package failed")

	}
	return &tp, nil
}

type TransferInRefundPackage struct {
	RefundAmount  *big.Int
	RefundAddress sdk.AccAddress
	RefundReason  RefundReason
}