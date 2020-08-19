package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	E8 = 100000000

	MYPC              = "mypc"
	DefaultBondDenom = MYPC // default bond denomination
)

func NewMypcCoin(amount int64) sdk.Coin {
	return sdk.NewCoin(MYPC, sdk.NewInt(amount))
}
func NewMypcCoinE8(amount int64) sdk.Coin {
	return sdk.NewCoin(MYPC, sdk.NewInt(amount*E8))
}

func NewMypcCoins(amount int64) sdk.Coins {
	return sdk.NewCoins(NewMypcCoin(amount))
}
func NewMypcCoinsE8(amount int64) sdk.Coins {
	return sdk.NewCoins(NewMypcCoin(amount * E8))
}
func NewCoins(denom string, amount int64) sdk.Coins {
	return sdk.NewCoins(sdk.NewCoin(denom, sdk.NewInt(amount)))
}
func IsMYPC(coin sdk.Coin) bool {
	return coin.Denom == MYPC
}
