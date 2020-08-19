package types

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestToString(t *testing.T) {
	sdk.GetConfig().SetBech32PrefixForAccount("mypos", "mypospub")
	fromAddr, _ := sdk.AccAddressFromBech32("mypos1px8alypku5j84qlwzdpynhn4nyrkagaytu5u4a")
	supervisor, _ := sdk.AccAddressFromBech32("mypos15fvnexrvsm9ryw3nn4mcrnqyhvhazkkrd4aqvd")
	lockedCoin := NewSupervisedLockedCoin("mypc", sdk.NewInt(100), 12345, fromAddr, supervisor, 1)
	require.Equal(t,
		"coin: 100mypc, unlocked_time: 12345, from: mypos1px8alypku5j84qlwzdpynhn4nyrkagaytu5u4a, supervisor: mypos15fvnexrvsm9ryw3nn4mcrnqyhvhazkkrd4aqvd, reward: 1\n",
		lockedCoin.String())

	lockedCoin2 := NewLockedCoin("mypc", sdk.NewInt(100), 12345)
	require.Equal(t,
		"coin: 100mypc, unlocked_time: 12345\n",
		lockedCoin2.String())

	lockedCoins := LockedCoins{lockedCoin, lockedCoin2}
	require.Equal(t,
		"coin: 100mypc, unlocked_time: 12345, from: mypos1px8alypku5j84qlwzdpynhn4nyrkagaytu5u4a, supervisor: mypos15fvnexrvsm9ryw3nn4mcrnqyhvhazkkrd4aqvd, reward: 1\n"+
			"coin: 100mypc, unlocked_time: 12345",
		lockedCoins.String())
}
