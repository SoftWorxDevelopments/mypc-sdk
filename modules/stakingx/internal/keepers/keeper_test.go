package keepers_test

import (
	"math"
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/cosmos/cosmos-sdk/x/staking"

	"github.com/SoftWorxDevelopments/mypc-sdk/modules/stakingx"
	"github.com/SoftWorxDevelopments/mypc-sdk/modules/stakingx/internal/keepers"
	"github.com/SoftWorxDevelopments/mypc-sdk/testapp"
	"github.com/SoftWorxDevelopments/mypc-sdk/testutil"
	myposchain "github.com/SoftWorxDevelopments/mypc-sdk/types"
)

func setUpInput() (keepers.MockKeeper, sdk.Context, auth.AccountKeeper) {
	testApp := testapp.NewTestApp()
	ctx := testApp.NewCtx()
	testApp.BankKeeper.SetSendEnabled(ctx, true)

	keeper := keepers.InitStates(ctx, testApp.StakingXKeeper, testApp.AccountKeeper, testApp.SupplyKeeper)

	return keeper, ctx, testApp.AccountKeeper
}

func TestInitExportGenesis(t *testing.T) {
	sxk, ctx, _ := setUpInput()

	genesisState := stakingx.GenesisState{
		Params: stakingx.Params{
			MinSelfDelegation:          stakingx.DefaultMinSelfDelegation,
			MinMandatoryCommissionRate: stakingx.DefaultMinMandatoryCommissionRate,
		},
	}

	stakingx.InitGenesis(ctx, sxk.Keeper, genesisState)
	exportGenesis := stakingx.ExportGenesis(ctx, sxk.Keeper)
	require.Equal(t, genesisState, exportGenesis)
}

func TestCalcBondPoolStatus(t *testing.T) {
	//initialize keeper & params &state
	sxk, ctx, _ := setUpInput()

	_, _, addr := testutil.KeyPubAddr()
	testParam := stakingx.Params{
		MinSelfDelegation: 0,
	}
	acc := auth.BaseAccount{
		Address: addr,
		Coins:   myposchain.NewMypcCoins(1e8),
	}
	vacc := auth.NewDelayedVestingAccount(&acc, math.MaxInt64)
	sxk.Ak.SetAccount(ctx, vacc)
	stakingx.InitGenesis(ctx, sxk.Keeper, stakingx.GenesisState{Params: testParam})

	feePool := types.FeePool{
		CommunityPool: sdk.NewDecCoins(myposchain.NewMypcCoins(1000)),
	}
	sxk.Dk.SetFeePool(ctx, feePool)

	bondedAcc := sxk.SupplyKeeper.GetModuleAccount(ctx, staking.BondedPoolName)
	bondedAcc.SetCoins(myposchain.NewMypcCoins(1000))
	sxk.Ak.SetAccount(ctx, bondedAcc)

	bondedAcc = sxk.SupplyKeeper.GetModuleAccount(ctx, staking.BondedPoolName)
	notBondedAcc := sxk.SupplyKeeper.GetModuleAccount(ctx, staking.NotBondedPoolName)
	expectedNonBondableTokens := feePool.CommunityPool.AmountOf("mypc").Add(acc.Coins.AmountOf("mypc").ToDec())
	expectedTotalSupply := sxk.SupplyKeeper.GetSupply(ctx).GetTotal().AmountOf("mypc")
	expectedBondRatio := bondedAcc.GetCoins().AmountOf("mypc").ToDec().QuoInt(expectedTotalSupply.Sub(expectedNonBondableTokens.RoundInt()))

	//test
	bondPool := sxk.CalcBondPoolStatus(ctx)

	require.Equal(t, expectedNonBondableTokens, bondPool.NonBondableTokens.ToDec())
	require.Equal(t, expectedTotalSupply, bondPool.TotalSupply)
	require.Equal(t, expectedBondRatio, bondPool.BondRatio)
	require.Equal(t, bondedAcc.GetCoins().AmountOf("mypc"), bondPool.BondedTokens)
	require.Equal(t, notBondedAcc.GetCoins().AmountOf("mypc"), bondPool.NotBondedTokens)
}

func TestCalcBondedRatio(t *testing.T) {
	bondPool := keepers.BondPool{
		BondedTokens:      sdk.NewInt(10e8),
		NotBondedTokens:   sdk.NewInt(500e8),
		NonBondableTokens: sdk.NewInt(10000),
		TotalSupply:       sdk.NewInt(510e8),
	}
	expectedBondRatio := bondPool.BondedTokens.ToDec().QuoInt(bondPool.TotalSupply.Sub(bondPool.NonBondableTokens))
	require.Equal(t, expectedBondRatio, keepers.CalcBondedRatio(&bondPool))
}

func TestCalcBondedRatioNegative(t *testing.T) {
	bondPool := keepers.BondPool{
		BondedTokens:      sdk.NewInt(-10e8),
		NotBondedTokens:   sdk.NewInt(500e8),
		NonBondableTokens: sdk.NewInt(10000),
		TotalSupply:       sdk.NewInt(510e8),
	}
	require.Equal(t, sdk.ZeroDec(), keepers.CalcBondedRatio(&bondPool))
}
