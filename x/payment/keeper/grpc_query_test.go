package keeper_test

import (
	"testing"
	"time"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"github.com/bnb-chain/greenfield/testutil/sample"
	"github.com/bnb-chain/greenfield/x/payment/types"
)

func TestParamsQuery(t *testing.T) {
	keeper, ctx, _ := makePaymentKeeper(t)
	params := types.DefaultParams()
	err := keeper.SetParams(ctx, params)
	require.NoError(t, err)

	response, err := keeper.Params(ctx, &types.QueryParamsRequest{})
	require.NoError(t, err)
	require.Equal(t, &types.QueryParamsResponse{Params: params}, response)
}

func TestParamsByTimestampQuery(t *testing.T) {
	keeper, ctx, _ := makePaymentKeeper(t)
	params := types.DefaultParams()

	before := time.Now()
	ctx = ctx.WithBlockTime(before)
	err := keeper.SetParams(ctx, params)
	require.NoError(t, err)

	after := time.Unix(before.Unix()+10, 0)
	ctx = ctx.WithBlockTime(after)
	newReserveTime := uint64(1000000000)
	params.VersionedParams.ReserveTime = newReserveTime
	err = keeper.SetParams(ctx, params)
	require.NoError(t, err)

	response, err := keeper.ParamsByTimestamp(ctx, &types.QueryParamsByTimestampRequest{
		Timestamp: before.Unix(),
	})
	require.NoError(t, err)
	require.True(t, newReserveTime != response.Params.VersionedParams.ReserveTime)

	response, err = keeper.ParamsByTimestamp(ctx, &types.QueryParamsByTimestampRequest{
		Timestamp: after.Unix(),
	})
	require.NoError(t, err)
	require.True(t, newReserveTime == response.Params.VersionedParams.ReserveTime)
}

func TestAutoSettleRecordQuery(t *testing.T) {
	keeper, ctx, _ := makePaymentKeeper(t)
	params := types.DefaultParams()
	err := keeper.SetParams(ctx, params)
	require.NoError(t, err)

	record := types.AutoSettleRecord{
		Timestamp: 123,
		Addr:      sample.RandAccAddress().String(),
	}
	keeper.SetAutoSettleRecord(ctx, &record)

	response, err := keeper.AutoSettleRecords(ctx, &types.QueryAutoSettleRecordsRequest{})
	require.NoError(t, err)
	require.Equal(t, record, response.AutoSettleRecords[0])
}

func TestDynamicBalanceQuery(t *testing.T) {
	keeper, ctx, deepKeepers := makePaymentKeeper(t)
	params := types.DefaultParams()
	err := keeper.SetParams(ctx, params)
	require.NoError(t, err)

	deepKeepers.AccountKeeper.EXPECT().HasAccount(gomock.Any(), gomock.Any()).
		Return(true).AnyTimes()
	bankBalance := sdk.NewCoin("BNB", sdkmath.NewInt(1000))

	deepKeepers.BankKeeper.EXPECT().GetBalance(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(bankBalance).AnyTimes()

	record := types.NewStreamRecord(sample.RandAccAddress(), ctx.BlockTime().Unix())
	record.StaticBalance = sdkmath.NewInt(100)
	keeper.SetStreamRecord(ctx, record)

	response, err := keeper.DynamicBalance(ctx, &types.QueryDynamicBalanceRequest{Account: record.Account})
	require.NoError(t, err)
	require.Equal(t, record.StaticBalance.Add(bankBalance.Amount), response.AvailableBalance)
	require.Equal(t, bankBalance.Amount, response.BankBalance)
}

func TestPaymentAccountAllQuery(t *testing.T) {
	keeper, ctx, _ := makePaymentKeeper(t)
	params := types.DefaultParams()
	err := keeper.SetParams(ctx, params)
	require.NoError(t, err)

	owner1 := sample.RandAccAddress()
	record1 := types.PaymentAccount{
		Owner: owner1.String(),
		Addr:  sample.RandAccAddress().String(),
	}
	keeper.SetPaymentAccount(ctx, &record1)

	owner2 := sample.RandAccAddress()
	record2 := types.PaymentAccount{
		Owner: owner2.String(),
		Addr:  sample.RandAccAddress().String(),
	}
	keeper.SetPaymentAccount(ctx, &record2)

	response, err := keeper.PaymentAccounts(ctx, &types.QueryPaymentAccountsRequest{})
	require.NoError(t, err)
	require.Equal(t, 2, len(response.PaymentAccounts))
}

func TestPaymentAccountQuery(t *testing.T) {
	keeper, ctx, _ := makePaymentKeeper(t)
	params := types.DefaultParams()
	err := keeper.SetParams(ctx, params)
	require.NoError(t, err)

	owner1 := sample.RandAccAddress()
	addr1 := sample.RandAccAddress().String()
	record1 := types.PaymentAccount{
		Owner: owner1.String(),
		Addr:  addr1,
	}
	keeper.SetPaymentAccount(ctx, &record1)

	owner2 := sample.RandAccAddress()
	addr2 := sample.RandAccAddress().String()
	record2 := types.PaymentAccount{
		Owner: owner2.String(),
		Addr:  addr2,
	}
	keeper.SetPaymentAccount(ctx, &record2)

	response, err := keeper.PaymentAccount(ctx, &types.QueryPaymentAccountRequest{
		Addr: addr1,
	})
	require.NoError(t, err)
	require.Equal(t, owner1.String(), response.PaymentAccount.Owner)

	response, err = keeper.PaymentAccount(ctx, &types.QueryPaymentAccountRequest{
		Addr: addr2,
	})
	require.NoError(t, err)
	require.Equal(t, owner2.String(), response.PaymentAccount.Owner)
}

func TestPaymentAccountCountAllQuery(t *testing.T) {
	keeper, ctx, _ := makePaymentKeeper(t)
	params := types.DefaultParams()
	err := keeper.SetParams(ctx, params)
	require.NoError(t, err)

	owner1 := sample.RandAccAddress()
	record1 := types.PaymentAccountCount{
		Owner: owner1.String(),
		Count: 10,
	}
	keeper.SetPaymentAccountCount(ctx, &record1)

	owner2 := sample.RandAccAddress()
	record2 := types.PaymentAccountCount{
		Owner: owner2.String(),
		Count: 2,
	}
	keeper.SetPaymentAccountCount(ctx, &record2)

	response, err := keeper.PaymentAccountCounts(ctx, &types.QueryPaymentAccountCountsRequest{})
	require.NoError(t, err)
	require.Equal(t, 2, len(response.PaymentAccountCounts))
}

func TestPaymentAccountCountQuery(t *testing.T) {
	keeper, ctx, _ := makePaymentKeeper(t)
	params := types.DefaultParams()
	err := keeper.SetParams(ctx, params)
	require.NoError(t, err)

	owner1 := sample.RandAccAddress()
	record1 := types.PaymentAccountCount{
		Owner: owner1.String(),
		Count: 10,
	}
	keeper.SetPaymentAccountCount(ctx, &record1)

	owner2 := sample.RandAccAddress()
	record2 := types.PaymentAccountCount{
		Owner: owner2.String(),
		Count: 2,
	}
	keeper.SetPaymentAccountCount(ctx, &record2)

	response, err := keeper.PaymentAccountCount(ctx, &types.QueryPaymentAccountCountRequest{
		Owner: owner1.String(),
	})
	require.NoError(t, err)
	require.Equal(t, record1.Count, response.PaymentAccountCount.Count)

	response, err = keeper.PaymentAccountCount(ctx, &types.QueryPaymentAccountCountRequest{
		Owner: owner2.String(),
	})
	require.NoError(t, err)
	require.Equal(t, record2.Count, response.PaymentAccountCount.Count)
}

func TestPaymentAccountsByOwnerQuery(t *testing.T) {
	keeper, ctx, _ := makePaymentKeeper(t)
	params := types.DefaultParams()
	err := keeper.SetParams(ctx, params)
	require.NoError(t, err)

	owner1 := sample.RandAccAddress()
	record1 := types.PaymentAccountCount{
		Owner: owner1.String(),
		Count: 2,
	}
	keeper.SetPaymentAccountCount(ctx, &record1)

	response, err := keeper.PaymentAccountsByOwner(ctx, &types.QueryPaymentAccountsByOwnerRequest{
		Owner: owner1.String(),
	})
	require.NoError(t, err)
	require.Equal(t, int(record1.Count), len(response.PaymentAccounts))
}

func TestStreamRecordAllQuery(t *testing.T) {
	keeper, ctx, _ := makePaymentKeeper(t)
	params := types.DefaultParams()
	err := keeper.SetParams(ctx, params)
	require.NoError(t, err)

	owner1 := sample.RandAccAddress()
	record1 := types.NewStreamRecord(owner1, ctx.BlockTime().Unix())
	keeper.SetStreamRecord(ctx, record1)

	owner2 := sample.RandAccAddress()
	record2 := types.NewStreamRecord(owner2, ctx.BlockTime().Unix())
	keeper.SetStreamRecord(ctx, record2)

	response, err := keeper.StreamRecords(ctx, &types.QueryStreamRecordsRequest{})
	require.NoError(t, err)
	require.Equal(t, 2, len(response.StreamRecords))
}

func TestStreamRecordQuery(t *testing.T) {
	keeper, ctx, _ := makePaymentKeeper(t)
	params := types.DefaultParams()
	err := keeper.SetParams(ctx, params)
	require.NoError(t, err)

	owner1 := sample.RandAccAddress()
	record1 := types.NewStreamRecord(owner1, ctx.BlockTime().Unix())
	keeper.SetStreamRecord(ctx, record1)

	owner2 := sample.RandAccAddress()
	record2 := types.NewStreamRecord(owner2, ctx.BlockTime().Unix())
	keeper.SetStreamRecord(ctx, record2)

	response, err := keeper.StreamRecord(ctx, &types.QueryGetStreamRecordRequest{
		Account: owner1.String(),
	})
	require.NoError(t, err)
	require.Equal(t, owner1.String(), response.StreamRecord.Account)

	response, err = keeper.StreamRecord(ctx, &types.QueryGetStreamRecordRequest{
		Account: owner2.String(),
	})
	require.NoError(t, err)
	require.Equal(t, owner2.String(), response.StreamRecord.Account)
}

func TestOutFlowQuery(t *testing.T) {
	keeper, ctx, _ := makePaymentKeeper(t)
	params := types.DefaultParams()
	err := keeper.SetParams(ctx, params)
	require.NoError(t, err)

	owner := sample.RandAccAddress()
	record1 := types.OutFlow{
		ToAddress: sample.RandAccAddress().String(),
		Rate:      sdkmath.Int{},
		Status:    types.OUT_FLOW_STATUS_FROZEN,
	}
	keeper.SetOutFlow(ctx, owner, &record1)

	record2 := types.OutFlow{
		ToAddress: sample.RandAccAddress().String(),
		Rate:      sdkmath.Int{},
		Status:    types.OUT_FLOW_STATUS_ACTIVE,
	}
	keeper.SetOutFlow(ctx, owner, &record2)

	response, err := keeper.OutFlows(ctx, &types.QueryOutFlowsRequest{
		Account: owner.String(),
	})
	require.NoError(t, err)
	require.Equal(t, 2, len(response.OutFlows))
}

func TestDelayedWithdrawalQuery(t *testing.T) {
	keeper, ctx, _ := makePaymentKeeper(t)
	params := types.DefaultParams()
	err := keeper.SetParams(ctx, params)
	require.NoError(t, err)

	owner := sample.RandAccAddress()
	from := sample.RandAccAddress()
	_, err = keeper.DelayedWithdrawal(ctx, &types.QueryDelayedWithdrawalRequest{
		Account: owner.String(),
	})
	require.Error(t, err)

	delayedWithdrawal := &types.DelayedWithdrawalRecord{
		Addr:            owner.String(),
		Amount:          sdkmath.NewInt(100),
		From:            from.String(),
		UnlockTimestamp: time.Now().Unix(),
	}
	keeper.SetDelayedWithdrawalRecord(ctx, delayedWithdrawal)

	response, err := keeper.DelayedWithdrawal(ctx, &types.QueryDelayedWithdrawalRequest{
		Account: owner.String(),
	})
	require.NoError(t, err)
	require.True(t, response.DelayedWithdrawal.Addr == delayedWithdrawal.Addr)
	require.True(t, response.DelayedWithdrawal.From == from.String())
	require.True(t, response.DelayedWithdrawal.Amount.Equal(delayedWithdrawal.Amount))
	require.True(t, response.DelayedWithdrawal.UnlockTimestamp == delayedWithdrawal.UnlockTimestamp)
}
