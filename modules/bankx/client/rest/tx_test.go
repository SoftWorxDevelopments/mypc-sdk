package rest

import (
	"net/http"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/SoftWorxDevelopments/mypc-sdk/modules/bankx/internal/types"
	myposchain "github.com/SoftWorxDevelopments/mypc-sdk/types"
)

func TestCmd(t *testing.T) {
	myposchain.InitSdkConfig()
	sendReq := sendReq{
		Amount:     myposchain.NewMypcCoins(100000000),
		UnlockTime: 0,
	}
	addr, _ := sdk.AccAddressFromBech32("mypos1px8alypku5j84qlwzdpynhn4nyrkagaytu5u4a")
	req := &http.Request{Method: "POST", URL: nil}
	req = mux.SetURLVars(req, map[string]string{"address": "mypos1px8alypku5j84qlwzdpynhn4nyrkagaytu5u4a"})
	msg, _ := sendReq.GetMsg(req, addr)
	assert.Equal(t, types.MsgSend{
		FromAddress: addr,
		ToAddress:   addr,
		Amount:      myposchain.NewMypcCoins(100000000),
		UnlockTime:  0,
	}, msg)

	memoReq := memoReq{
		Required: true,
	}
	msg, _ = memoReq.GetMsg(req, addr)
	assert.Equal(t, types.MsgSetMemoRequired{
		Address:  addr,
		Required: true,
	}, msg)
}
