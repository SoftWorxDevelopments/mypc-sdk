package incentive_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/supply"

	"github.com/SoftWorxDevelopments/mypc-sdk/modules/incentive"
	"github.com/SoftWorxDevelopments/mypc-sdk/modules/incentive/internal/types"
	"github.com/SoftWorxDevelopments/mypc-sdk/testapp"
	myposchain "github.com/SoftWorxDevelopments/mypc-sdk/types"
)

type TestInput struct {
	ctx    sdk.Context
	cdc    *codec.Codec
	keeper incentive.Keeper
	ak     auth.AccountKeeper
	sk     supply.Keeper
}

func SetupTestInput() TestInput {
	app := testapp.NewTestApp()
	ctx := sdk.NewContext(app.Cms, abci.Header{ChainID: "test-chain-id"}, false, log.NewNopLogger())
	return TestInput{ctx: ctx, cdc: app.Cdc, keeper: app.IncentiveKeeper, ak: app.AccountKeeper, sk: app.SupplyKeeper}
}

func TestBeginBlockerInvalidCoin(t *testing.T) {
	input := SetupTestInput()
	_ = input.keeper.SetState(input.ctx, incentive.State{HeightAdjustment: 10})
	input.keeper.SetParams(input.ctx, incentive.DefaultParams())

	feeBalanceBefore := input.sk.GetModuleAccount(input.ctx, auth.FeeCollectorName).GetCoins().AmountOf(myposchain.MYPC).Int64()
	incentive.BeginBlocker(input.ctx, input.keeper)
	feeBalanceAfter := input.sk.GetModuleAccount(input.ctx, auth.FeeCollectorName).GetCoins().AmountOf(myposchain.MYPC).Int64()

	// no coins in pool
	require.Equal(t, int64(0), feeBalanceAfter-feeBalanceBefore)
}

func TestBeginBlockerInPlan(t *testing.T) {
	input := SetupTestInput()
	plans := types.Params{
		DefaultRewardPerBlock: 2e8,
		Plans: []types.Plan{
			{
				StartHeight:    0,
				EndHeight:      10,
				RewardPerBlock: 10e8,
				TotalIncentive: 100e8,
			},
		},
	}
	_ = input.keeper.SetState(input.ctx, incentive.State{HeightAdjustment: 10})
	input.keeper.SetParams(input.ctx, plans)
	acc := input.ak.NewAccountWithAddress(input.ctx, incentive.PoolAddr)
	_ = acc.SetCoins(myposchain.NewMypcCoins(10000 * 1e8))
	input.ak.SetAccount(input.ctx, acc)

	poolBalanceBefore := input.ak.GetAccount(input.ctx, incentive.PoolAddr).GetCoins().AmountOf(myposchain.MYPC).Int64()
	feeBalanceBefore := input.sk.GetModuleAccount(input.ctx, auth.FeeCollectorName).GetCoins().AmountOf(myposchain.MYPC).Int64()

	incentive.BeginBlocker(input.ctx, input.keeper)

	poolBalanceAfter := input.ak.GetAccount(input.ctx, incentive.PoolAddr).GetCoins().AmountOf(myposchain.MYPC).Int64()
	feeBalanceAfter := input.sk.GetModuleAccount(input.ctx, auth.FeeCollectorName).GetCoins().AmountOf(myposchain.MYPC).Int64()

	reward := plans.Plans[0].RewardPerBlock
	require.Equal(t, -reward, poolBalanceAfter-poolBalanceBefore)
	require.Equal(t, reward, feeBalanceAfter-feeBalanceBefore)
}

func TestBeginBlockerNotInPlan(t *testing.T) {
	input := SetupTestInput()
	plans := types.Params{
		DefaultRewardPerBlock: 2e8,
		Plans: []types.Plan{
			{
				StartHeight:    0,
				EndHeight:      10,
				RewardPerBlock: 10e8,
				TotalIncentive: 100e8,
			},
		},
	}
	_ = input.keeper.SetState(input.ctx, incentive.State{HeightAdjustment: 20})
	input.keeper.SetParams(input.ctx, plans)
	acc := input.ak.NewAccountWithAddress(input.ctx, incentive.PoolAddr)
	_ = acc.SetCoins(myposchain.NewMypcCoins(10000 * 1e8))
	input.ak.SetAccount(input.ctx, acc)

	poolBalanceBefore := input.ak.GetAccount(input.ctx, incentive.PoolAddr).GetCoins().AmountOf(myposchain.MYPC).Int64()
	feeBalanceBefore := input.sk.GetModuleAccount(input.ctx, auth.FeeCollectorName).GetCoins().AmountOf(myposchain.MYPC).Int64()

	incentive.BeginBlocker(input.ctx, input.keeper)

	poolBalanceAfter := input.ak.GetAccount(input.ctx, incentive.PoolAddr).GetCoins().AmountOf(myposchain.MYPC).Int64()
	feeBalanceAfter := input.sk.GetModuleAccount(input.ctx, auth.FeeCollectorName).GetCoins().AmountOf(myposchain.MYPC).Int64()

	reward := plans.DefaultRewardPerBlock
	require.Equal(t, -reward, poolBalanceAfter-poolBalanceBefore)
	require.Equal(t, reward, feeBalanceAfter-feeBalanceBefore)
}

func TestIncentiveCoinsAddress(t *testing.T) {
	require.Equal(t, "mypos1gc5t98jap4zyhmhmyq5af5s7pyv57w5694el97", incentive.PoolAddr.String())
}

func TestIncentiveCoinsAddressInTestNet(t *testing.T) {
	config := sdk.GetConfig()
	testnetAddrPrefix := "mypctest"
	config.SetBech32PrefixForAccount(testnetAddrPrefix, testnetAddrPrefix+sdk.PrefixPublic)
	require.Equal(t, "mypctest1gc5t98jap4zyhmhmyq5af5s7pyv57w566ewmx0", incentive.PoolAddr.String())
}

func TestMain(m *testing.M) {
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(myposchain.Bech32MainPrefix, myposchain.Bech32MainPrefix+sdk.PrefixPublic)
	os.Exit(m.Run())
}
