package types

import (
	"reflect"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestNewMypcCoin(t *testing.T) {
	coin := NewMypcCoin(1)
	if coin.Amount.Int64() != 1 {
		t.Error("coin is not 1")
	}

	coin = NewMypcCoin(0)
	if coin.Amount.Int64() != 0 {
		t.Error("coin is not 0")
	}
}

func TestNewMypcCoins(t *testing.T) {
	coins := NewMypcCoins(1)
	if coins[0].Amount.Int64() != 1 {
		t.Error("coin is not 1")
	}
}

func TestNewMypcCoinE8(t *testing.T) {
	type args struct {
		amount int64
	}
	tests := []struct {
		name string
		args args
		want sdk.Coin
	}{
		{name: "mypc", args: args{1}, want: sdk.NewInt64Coin("mypc", E8)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewMypcCoinE8(tt.args.amount); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewMypcCoinE8() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewMypcCoinsE8(t *testing.T) {
	type args struct {
		amount int64
	}
	tests := []struct {
		name string
		args args
		want sdk.Coins
	}{
		{name: "mypc", args: args{1}, want: []sdk.Coin{sdk.NewInt64Coin("mypc", E8)}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewMypcCoinsE8(tt.args.amount); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewMypcCoinsE8() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsMYPC(t *testing.T) {
	type args struct {
		coin sdk.Coin
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{name: "mypc", args: args{sdk.NewInt64Coin("mypc", 1)}, want: true},
		{name: "btc", args: args{sdk.NewInt64Coin("btc", 1)}, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsMYPC(tt.args.coin); got != tt.want {
				t.Errorf("IsMYPC() = %v, want %v", got, tt.want)
			}
		})
	}
}
