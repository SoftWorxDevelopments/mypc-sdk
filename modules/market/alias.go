package market

import (
	"github.com/SoftWorxDevelopments/mypc-sdk/modules/market/internal/keepers"
	"github.com/SoftWorxDevelopments/mypc-sdk/modules/market/internal/types"
	myposchain "github.com/SoftWorxDevelopments/mypc-sdk/types"
)

const (
	StoreKey   = types.StoreKey
	ModuleName = types.ModuleName
)

const (
	IntegrationNetSubString = types.IntegrationNetSubString
	OrderIDPartsNum         = types.OrderIDPartsNum
	SymbolSeparator         = types.SymbolSeparator
	LimitOrder              = types.LimitOrder
	GTE                     = types.GTE
	BID                     = types.BID
	ASK                     = types.ASK
	BUY                     = types.BUY
	SELL                    = types.SELL
)

var (
	NewBaseKeeper       = keepers.NewKeeper
	DefaultParams       = types.DefaultParams
	DecToBigEndianBytes = types.DecToBigEndianBytes
	ValidateOrderID     = types.ValidateOrderID
	IsValidTradingPair  = types.IsValidTradingPair
	ModuleCdc           = types.ModuleCdc
	GetSymbol           = myposchain.GetSymbol
	SplitSymbol         = myposchain.SplitSymbol
)

type (
	Keeper                  = keepers.Keeper
	Order                   = types.Order
	MarketInfo              = types.MarketInfo
	Params                  = types.Params
	MsgCreateOrder          = types.MsgCreateOrder
	MsgCreateTradingPair    = types.MsgCreateTradingPair
	MsgCancelOrder          = types.MsgCancelOrder
	MsgCancelTradingPair    = types.MsgCancelTradingPair
	MsgModifyPricePrecision = types.MsgModifyPricePrecision
	CreateOrderInfo         = types.CreateOrderInfo
	FillOrderInfo           = types.FillOrderInfo
	CancelOrderInfo         = types.CancelOrderInfo
)
