package keepers

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/staking"

	"github.com/SoftWorxDevelopments/mypc-sdk/modules/stakingx/internal/types"
	myposchain "github.com/SoftWorxDevelopments/mypc-sdk/types"
)

var (
	NonBondableAddressesKey = []byte("0x01")
)

type Keeper struct {
	key sdk.StoreKey
	cdc *codec.Codec

	paramSubspace params.Subspace

	assetViewKeeper AssetViewKeeper

	sk *staking.Keeper

	dk DistributionKeeper

	ak auth.AccountKeeper

	bk ExpectBankxKeeper

	supplyKeeper ExpectSupplyKeeper

	feeCollectorName string
}

func NewKeeper(key sdk.StoreKey, cdc *codec.Codec,
	paramSubspace params.Subspace, assetViewKeeper AssetViewKeeper, sk *staking.Keeper,
	dk DistributionKeeper, ak auth.AccountKeeper, bk ExpectBankxKeeper,
	supplyKeeper ExpectSupplyKeeper, feeCollectorName string) Keeper {

	return Keeper{
		key:              key,
		cdc:              cdc,
		paramSubspace:    paramSubspace.WithKeyTable(types.ParamKeyTable()),
		assetViewKeeper:  assetViewKeeper,
		sk:               sk,
		dk:               dk,
		ak:               ak,
		bk:               bk,
		supplyKeeper:     supplyKeeper,
		feeCollectorName: feeCollectorName,
	}
}

// -----------------------------------------------------------------------------
// Params

// SetParams sets the asset module's parameters.
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramSubspace.SetParamSet(ctx, &params)
}

// GetParams gets the asset module's parameters.
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	k.paramSubspace.GetParamSet(ctx, &params)
	return
}

func (k Keeper) GetMinMandatoryCommissionRate(ctx sdk.Context) (rate sdk.Dec) {
	k.paramSubspace.Get(ctx, types.KeyMinMandatoryCommissionRate, &rate)
	return
}

// -----------------------------------------------------------------------------
// BondPoolStatus

func (k Keeper) CalcBondPoolStatus(ctx sdk.Context) BondPool {
	total := k.supplyKeeper.GetSupply(ctx).GetTotal().AmountOf(myposchain.MYPC)
	var bondPool BondPool

	bondPool.TotalSupply = total
	bondPool.BondedTokens = k.supplyKeeper.GetModuleAccount(ctx, staking.BondedPoolName).GetCoins().AmountOf(myposchain.MYPC)
	bondPool.NotBondedTokens = k.supplyKeeper.GetModuleAccount(ctx, staking.NotBondedPoolName).GetCoins().AmountOf(myposchain.MYPC)
	bondPool.NonBondableTokens = calcNonBondableTokens(ctx, &k)

	bondPool.BondRatio = calcBondedRatio(&bondPool)

	return bondPool
}

func calcBondedRatio(p *BondPool) sdk.Dec {
	if p.BondedTokens.IsNegative() || p.NonBondableTokens.IsNegative() {
		return sdk.ZeroDec()
	}

	bondableTokens := p.TotalSupply.Sub(p.NonBondableTokens)
	if !bondableTokens.IsPositive() {
		return sdk.ZeroDec()
	}

	return p.BondedTokens.ToDec().QuoInt(bondableTokens)
}

func calcNonBondableTokens(ctx sdk.Context, k *Keeper) sdk.Int {
	ret := sdk.ZeroInt()
	addrs := k.getNonBondableAddresses(ctx)

	for _, addr := range addrs {
		if acc := k.ak.GetAccount(ctx, addr); acc != nil {
			if amt := acc.GetCoins().AmountOf(myposchain.MYPC); amt.IsPositive() {
				ret = ret.Add(amt)
			}
		}
	}

	communityPoolAmt := k.dk.GetFeePoolCommunityCoins(ctx).AmountOf(myposchain.MYPC)
	ret = ret.Add(communityPoolAmt.TruncateInt())

	return ret
}

// -----------------------------------------------------------------------------
// Non-bondable addresses

func (k Keeper) GetMypcOwnerAddress(ctx sdk.Context) sdk.AccAddress {
	mypc := k.assetViewKeeper.GetToken(ctx, myposchain.MYPC)
	if mypc == nil {
		return nil
	}
	return mypc.GetOwner()
}

func (k Keeper) GetAllVestingAccountAddresses(ctx sdk.Context) []sdk.AccAddress {
	addresses := make([]sdk.AccAddress, 0, 8)
	k.ak.IterateAccounts(ctx, func(acc auth.Account) bool {
		if vacc, ok := acc.(auth.VestingAccount); ok {
			addresses = append(addresses, vacc.GetAddress())
		}
		return false
	})
	return addresses
}

func (k Keeper) SetNonBondableAddresses(ctx sdk.Context, addresses []sdk.AccAddress) {
	store := ctx.KVStore(k.key)
	bz, err := k.cdc.MarshalBinaryBare(addresses)
	if err != nil {
		panic(err)
	}
	store.Set(NonBondableAddressesKey, bz)
}

func (k Keeper) getNonBondableAddresses(ctx sdk.Context) (addresses []sdk.AccAddress) {
	store := ctx.KVStore(k.key)
	bz := store.Get(NonBondableAddressesKey)
	if bz == nil {
		return
	}

	err := k.cdc.UnmarshalBinaryBare(bz, &addresses)
	if err != nil {
		panic(err) // TODO
	}
	return
}
