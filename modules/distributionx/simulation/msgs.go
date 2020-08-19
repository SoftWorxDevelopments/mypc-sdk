package simulation

import (
	"fmt"
	"math/rand"

	"github.com/cosmos/cosmos-sdk/x/distribution"

	"github.com/cosmos/cosmos-sdk/x/supply"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"github.com/SoftWorxDevelopments/mypc-sdk/modules/distributionx"
	myposchainsim "github.com/SoftWorxDevelopments/mypc-sdk/simulation"
	myposchain "github.com/SoftWorxDevelopments/mypc-sdk/types"
)

func SimulateMsgDonateToCommunityPool(ak auth.AccountKeeper, dxk distributionx.Keeper) simulation.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simulation.Account) (
		opMsg simulation.OperationMsg, fOps []simulation.FutureOperation, err error) {

		acc := simulation.RandomAcc(r, accs)
		dbAcc := ak.GetAccount(ctx, acc.Address)
		mypcAmt := dbAcc.GetCoins().AmountOf(myposchain.MYPC)
		rdmAmt := simulation.RandomAmount(r, mypcAmt)

		if rdmAmt.LT(sdk.OneInt()) {
			return simulation.NoOpMsg(distributionx.ModuleName), nil, nil
		}

		msg := distributionx.MsgDonateToCommunityPool{
			FromAddr: acc.Address,
			Amount:   sdk.NewCoins(sdk.NewCoin(myposchain.MYPC, rdmAmt)),
		}

		oldCoins := getCommunityPoolCoins(ctx, ak)
		ok := myposchainsim.SimulateHandleMsg(msg, distributionx.NewHandler(dxk), ctx)
		if !ok {
			return simulation.NewOperationMsg(msg, ok, ""), nil, nil
		}

		ok = verifyDonateToCommunityPool(ctx, ak, oldCoins, msg)
		if !ok {
			return simulation.NewOperationMsg(msg, ok, ""), nil, fmt.Errorf("donation to community pool failed")
		}
		return simulation.NewOperationMsg(msg, ok, ""), nil, nil
	}
}

func getCommunityPoolCoins(ctx sdk.Context, ak auth.AccountKeeper) sdk.Coins {
	if acc := ak.GetAccount(ctx, supply.NewModuleAddress(distribution.ModuleName)); acc != nil {
		return acc.GetCoins()
	}
	return nil
}
func verifyDonateToCommunityPool(ctx sdk.Context, ak auth.AccountKeeper, oldCoins sdk.Coins, msg distributionx.MsgDonateToCommunityPool) bool {
	return getCommunityPoolCoins(ctx, ak).Sub(oldCoins).IsEqual(msg.Amount)
}
