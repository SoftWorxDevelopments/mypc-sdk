package rest

import (
	"net/http"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"

	"github.com/SoftWorxDevelopments/mypc-sdk/modules/market/internal/types"
)

func TestCmd(t *testing.T) {
	createMarket := createMarketReq{
		Stock:          "etc",
		Money:          "mypc",
		PricePrecision: 8,
	}
	addr, _ := sdk.AccAddressFromBech32("mypos1px8alypku5j84qlwzdpynhn4nyrkagaytu5u4a")
	msg, _ := createMarket.GetMsg(nil, addr)
	assert.Equal(t, types.MsgCreateTradingPair{
		Stock:          "etc",
		Money:          "mypc",
		Creator:        addr,
		PricePrecision: 8,
	}, msg)
	//==============
	cancelMarket := cancelMarketReq{
		TradingPair: "etc/mypc",
		Time:        12345678,
	}
	msg, _ = cancelMarket.GetMsg(nil, addr)
	assert.Equal(t, types.MsgCancelTradingPair{
		Sender:        addr,
		TradingPair:   "etc/mypc",
		EffectiveTime: 12345678,
	}, msg)
	//==============
	req := modifyPricePrecision{
		TradingPair:    "etc/mypc",
		PricePrecision: 9,
	}
	msg, _ = req.GetMsg(nil, addr)
	assert.Equal(t, types.MsgModifyPricePrecision{
		Sender:         addr,
		TradingPair:    "etc/mypc",
		PricePrecision: 9,
	}, msg)
	//==============
	createOrder := createOrderReq{
		OrderType:      types.LIMIT,
		TradingPair:    "etc/mypc",
		Identify:       0,
		PricePrecision: 8,
		Price:          12345678,
		Quantity:       123,
		Side:           types.SELL,
		ExistBlocks:    25000,
		TimeInForce:    types.GTE,
	}
	httpReq, _ := http.NewRequest("POST", "http://example.com/market/gte-orders", nil)
	msg, _ = createOrder.GetMsg(httpReq, addr)
	assert.Equal(t, types.MsgCreateOrder{
		Sender:         addr,
		Identify:       0,
		TradingPair:    "etc/mypc",
		OrderType:      types.LIMIT,
		PricePrecision: 8,
		Price:          12345678,
		Quantity:       123,
		Side:           types.SELL,
		TimeInForce:    types.GTE,
		ExistBlocks:    25000,
	}, msg)
	httpReq, _ = http.NewRequest("POST", "http://example.com/market/ioc-orders", nil)
	msg, _ = createOrder.GetMsg(httpReq, addr)
	assert.Equal(t, types.MsgCreateOrder{
		Sender:         addr,
		Identify:       0,
		TradingPair:    "etc/mypc",
		OrderType:      types.LIMIT,
		PricePrecision: 8,
		Price:          12345678,
		Quantity:       123,
		Side:           types.SELL,
		TimeInForce:    types.IOC,
		ExistBlocks:    25000,
	}, msg)
	//==============
	cancelOrder := cancelOrderReq{
		OrderID: "mypos1px8alypku5j84qlwzdpynhn4nyrkagaytu5u4a-1025",
	}
	msg, _ = cancelOrder.GetMsg(nil, addr)
	assert.Equal(t, &types.MsgCancelOrder{
		Sender:  addr,
		OrderID: "mypos1px8alypku5j84qlwzdpynhn4nyrkagaytu5u4a-1025",
	}, msg)
}
