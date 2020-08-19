package keeper_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/supply"

	"github.com/SoftWorxDevelopments/mypc-sdk/modules/asset"
	"github.com/SoftWorxDevelopments/mypc-sdk/modules/authx"
	"github.com/SoftWorxDevelopments/mypc-sdk/modules/bankx/internal/keeper"
	"github.com/SoftWorxDevelopments/mypc-sdk/modules/bankx/internal/types"
	"github.com/SoftWorxDevelopments/mypc-sdk/testapp"
	"github.com/SoftWorxDevelopments/mypc-sdk/testutil"
	myposchain "github.com/SoftWorxDevelopments/mypc-sdk/types"
)

var myaddr = testutil.ToAccAddress("myaddr")
var ownerAddr = testutil.ToAccAddress("owneraddr")

func defaultContext() (keeper.Keeper, sdk.Context) {
	app := testapp.NewTestApp()
	ctx := sdk.NewContext(app.Cms, abci.Header{}, false, log.NewNopLogger())
	app.AccountKeeper.SetAccount(ctx, supply.NewEmptyModuleAccount(authx.ModuleName))
	app.AccountKeeper.SetAccount(ctx, supply.NewEmptyModuleAccount(asset.ModuleName, supply.Minter))
	app.SupplyKeeper.SetSupply(ctx, supply.NewSupply(sdk.Coins{sdk.Coin{Denom: "abc", Amount: sdk.NewInt(10e10)}}))
	app.SupplyKeeper.SetSupply(ctx, supply.NewSupply(sdk.Coins{sdk.Coin{Denom: "mypc", Amount: sdk.NewInt(10e10)}}))

	_ = app.AssetKeeper.IssueToken(ctx, "abc", "abc", sdk.NewInt(100000000000), ownerAddr,
		false, false, false, false,
		"", "", "abc")
	_ = app.AssetKeeper.IssueToken(ctx, "mypc", "mypc", sdk.NewInt(100000000000), ownerAddr,
		false, false, false, false,
		"", "", "mypc")
	return app.BankxKeeper, ctx
}

func givenAccountWith(ctx sdk.Context, keeper keeper.Keeper, addr sdk.AccAddress, coinsString string) error {
	coins, _ := sdk.ParseCoins(coinsString)
	if err := keeper.AddCoins(ctx, addr, coins); err != nil {
		return err
	}
	return nil
}

func coinsOf(ctx sdk.Context, keeper keeper.Keeper, addr sdk.AccAddress) string {
	return keeper.GetCoins(ctx, addr).String()
}

func frozenCoinsOf(ctx sdk.Context, keeper keeper.Keeper, addr sdk.AccAddress) string {
	return keeper.GetFrozenCoins(ctx, addr).String()
}

func TestFreezeMultiCoins(t *testing.T) {
	bkx, ctx := defaultContext()
	err := givenAccountWith(ctx, bkx, myaddr, "1000000000mypc,100abc")
	require.NoError(t, err)

	freezeCoins, _ := sdk.ParseCoins("300000000mypc, 20abc")
	err = bkx.FreezeCoins(ctx, myaddr, freezeCoins)

	require.Nil(t, err)
	require.Equal(t, "80abc,700000000mypc", coinsOf(ctx, bkx, myaddr))
	require.Equal(t, "20abc,300000000mypc", frozenCoinsOf(ctx, bkx, myaddr))

	err = bkx.UnFreezeCoins(ctx, myaddr, freezeCoins)

	require.Nil(t, err)
	require.Equal(t, "100abc,1000000000mypc", coinsOf(ctx, bkx, myaddr))
	require.Equal(t, "", frozenCoinsOf(ctx, bkx, myaddr))
}

func TestFreezeUnFreezeOK(t *testing.T) {
	bkx, ctx := defaultContext()
	err := givenAccountWith(ctx, bkx, myaddr, "1000000000mypc")
	require.NoError(t, err)

	freezeCoins := myposchain.NewMypcCoins(300000000)
	err = bkx.FreezeCoins(ctx, myaddr, freezeCoins)

	require.Nil(t, err)
	require.Equal(t, "700000000mypc", coinsOf(ctx, bkx, myaddr))
	require.Equal(t, "300000000mypc", frozenCoinsOf(ctx, bkx, myaddr))

	err = bkx.UnFreezeCoins(ctx, myaddr, freezeCoins)

	require.Nil(t, err)
	require.Equal(t, "1000000000mypc", coinsOf(ctx, bkx, myaddr))
	require.Equal(t, "", frozenCoinsOf(ctx, bkx, myaddr))
}

func TestFreezeUnFreezeInvalidAccount(t *testing.T) {
	bkx, ctx := defaultContext()

	freezeCoins := myposchain.NewMypcCoins(500000000)
	err := bkx.FreezeCoins(ctx, myaddr, freezeCoins)
	require.Equal(t, sdk.ErrInsufficientCoins("insufficient account funds;  < 500000000mypc"), err)

	err = bkx.UnFreezeCoins(ctx, myaddr, freezeCoins)
	require.Equal(t, sdk.ErrUnknownAddress(fmt.Sprintf("account %s does not exist", myaddr)), err)
}

func TestFreezeUnFreezeInsufficientCoins(t *testing.T) {
	bkx, ctx := defaultContext()

	err := givenAccountWith(ctx, bkx, myaddr, "10mypc")
	require.NoError(t, err)

	InvalidFreezeCoins := myposchain.NewMypcCoins(50)
	err = bkx.FreezeCoins(ctx, myaddr, InvalidFreezeCoins)
	require.Equal(t, sdk.ErrInsufficientCoins("insufficient account funds; 10mypc < 50mypc"), err)

	freezeCoins := myposchain.NewMypcCoins(5)
	err = bkx.FreezeCoins(ctx, myaddr, freezeCoins)
	require.Nil(t, err)

	err = bkx.UnFreezeCoins(ctx, myaddr, InvalidFreezeCoins)
	require.Equal(t, sdk.ErrInsufficientCoins("account has insufficient coins to unfreeze"), err)
}

func TestGetTotalCoins(t *testing.T) {
	bkx, ctx := defaultContext()
	err := givenAccountWith(ctx, bkx, myaddr, "100mypc, 20bch, 30btc")
	require.NoError(t, err)

	lockedCoins := authx.LockedCoins{
		authx.NewLockedCoin("bch", sdk.NewInt(20), 1000),
		authx.NewLockedCoin("eth", sdk.NewInt(30), 2000),
	}

	frozenCoins := sdk.NewCoins(
		sdk.NewCoin("btc", sdk.NewInt(50)),
		sdk.NewCoin("eth", sdk.NewInt(10)),
	)

	bkx.MockAddLockedCoins(ctx, myaddr, lockedCoins)
	bkx.MockAddFrozenCoins(ctx, myaddr, frozenCoins)
	expected := sdk.NewCoins(
		sdk.NewCoin("bch", sdk.NewInt(40)),
		sdk.NewCoin("btc", sdk.NewInt(80)),
		sdk.NewCoin("mypc", sdk.NewInt(100)),
		sdk.NewCoin("eth", sdk.NewInt(40)),
	)
	expected = expected.Sort()
	coins := bkx.GetTotalCoins(ctx, myaddr)

	require.Equal(t, expected, coins)
}

func TestKeeper_TotalAmountOfCoin(t *testing.T) {

	bkx, ctx := defaultContext()
	amount := bkx.TotalAmountOfCoin(ctx, "mypc")
	require.Equal(t, int64(100000000000), amount.Int64())

	err := givenAccountWith(ctx, bkx, myaddr, "100mypc")
	require.NoError(t, err)

	lockedCoins := authx.LockedCoins{
		authx.NewLockedCoin("mypc", sdk.NewInt(100), 1000),
	}
	frozenCoins := sdk.NewCoins(sdk.NewCoin("mypc", sdk.NewInt(100)))

	bkx.MockAddLockedCoins(ctx, myaddr, lockedCoins)
	bkx.MockAddFrozenCoins(ctx, myaddr, frozenCoins)

	amount = bkx.TotalAmountOfCoin(ctx, "mypc")
	require.Equal(t, int64(100000000300), amount.Int64())
}

func TestKeeper_AddCoins(t *testing.T) {
	bkx, ctx := defaultContext()
	coins := sdk.NewCoins(
		sdk.NewCoin("aaa", sdk.NewInt(10)),
		sdk.NewCoin("bbb", sdk.NewInt(20)),
	)

	coins2 := sdk.NewCoins(
		sdk.NewCoin("aaa", sdk.NewInt(5)),
		sdk.NewCoin("bbb", sdk.NewInt(10)),
	)

	err := bkx.AddCoins(ctx, myaddr, coins)
	require.Equal(t, nil, err)
	err = bkx.SubtractCoins(ctx, myaddr, coins2)
	require.Equal(t, nil, err)
	cs := bkx.GetTotalCoins(ctx, myaddr)
	require.Equal(t, coins2, cs)

	coins3 := sdk.NewCoins(
		sdk.NewCoin("aaa", sdk.NewInt(15)),
		sdk.NewCoin("bbb", sdk.NewInt(10)),
	)
	err = bkx.SubtractCoins(ctx, myaddr, coins3)
	require.Error(t, err)
}

func TestParamGetSet(t *testing.T) {
	bkx, ctx := defaultContext()

	//expect DefaultActivationFees=1
	defaultParam := types.DefaultParams()
	require.Equal(t, int64(100000000), defaultParam.ActivationFee)

	//expect SetParam don't panic
	require.NotPanics(t, func() { bkx.SetParams(ctx, defaultParam) }, "bankxKeeper SetParam panics")

	//expect GetParam equals defaultParam
	require.Equal(t, defaultParam, bkx.GetParams(ctx))
}

func TestKeeper_SendCoins(t *testing.T) {
	bkx, ctx := defaultContext()
	coins := sdk.NewCoins(
		sdk.NewCoin("abc", sdk.NewInt(10)),
	)
	addr2 := testutil.ToAccAddress("addr2")
	_ = bkx.AddCoins(ctx, myaddr, coins)
	exist := bkx.HasCoins(ctx, myaddr, coins)
	assert.True(t, exist)
	err := bkx.SendCoins(ctx, myaddr, addr2, coins)
	require.Equal(t, nil, err)
	cs := bkx.GetTotalCoins(ctx, addr2)
	require.Equal(t, coins, cs)
}
