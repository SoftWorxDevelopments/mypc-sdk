package alias

import (
	"github.com/SoftWorxDevelopments/mypc-sdk/modules/alias/internal/keepers"
	"github.com/SoftWorxDevelopments/mypc-sdk/modules/alias/internal/types"
)

const (
	StoreKey   = types.StoreKey
	ModuleName = types.ModuleName
)

var (
	ModuleCdc     = types.ModuleCdc
	NewBaseKeeper = keepers.NewKeeper
	DefaultParams = types.DefaultParams
)

type (
	Keeper         = keepers.Keeper
	MsgAliasUpdate = types.MsgAliasUpdate
)
