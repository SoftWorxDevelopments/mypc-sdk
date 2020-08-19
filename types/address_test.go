package types

import (
	"encoding/hex"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestAddressPrefixes(t *testing.T) {
	InitSdkConfig()
	config := sdk.GetConfig()

	require.Equal(t, "mypos", config.GetBech32AccountAddrPrefix())
	require.Equal(t, "mypospub", config.GetBech32AccountPubPrefix())
	require.Equal(t, "myposvaloper", config.GetBech32ValidatorAddrPrefix())
	require.Equal(t, "myposvaloperpub", config.GetBech32ValidatorPubPrefix())
	require.Equal(t, "myposvalcons", config.GetBech32ConsensusAddrPrefix())
	require.Equal(t, "myposvalconspub", config.GetBech32ConsensusPubPrefix())
}

func randomHex(n int) (string, error) {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func TestAccountAddressLengthIs45(t *testing.T) {
	rand.Seed(time.Now().UTC().UnixNano())
	randomAddr, _ := randomHex(20)
	addr, _ := sdk.AccAddressFromHex(randomAddr)
	require.Equal(t, 45, len(addr.String()))
}
