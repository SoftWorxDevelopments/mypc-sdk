package rest

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"

	"github.com/SoftWorxDevelopments/mypc-sdk/modules/alias/internal/types"
)

func TestCmd(t *testing.T) {
	aliasUpdateReq := AliasUpdateReq{
		Alias:     "superboy",
		IsAdd:     true,
		AsDefault: true,
	}
	addr, _ := sdk.AccAddressFromBech32("mypos1px8alypku5j84qlwzdpynhn4nyrkagaytu5u4a")
	msg, _ := aliasUpdateReq.GetMsg(nil, addr)
	msgAliasUpdate, _ := msg.(*types.MsgAliasUpdate)
	assert.Equal(t, &types.MsgAliasUpdate{
		Owner:     addr,
		Alias:     "superboy",
		IsAdd:     true,
		AsDefault: true,
	}, msgAliasUpdate)
}
