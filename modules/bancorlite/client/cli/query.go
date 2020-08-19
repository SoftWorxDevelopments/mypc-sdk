package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/SoftWorxDevelopments/mypc-sdk/modules/bancorlite/internal/keepers"
	"github.com/SoftWorxDevelopments/mypc-sdk/modules/bancorlite/internal/types"
	myposchain "github.com/SoftWorxDevelopments/mypc-sdk/types"
	"github.com/coinexchain/cosmos-utils/client/cliutil"
)

func QueryParamsCmd(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "params",
		Args:  cobra.NoArgs,
		Short: "Query bancorlite params",
		RunE: func(cmd *cobra.Command, args []string) error {
			route := fmt.Sprintf("custom/%s/%s", types.StoreKey, keepers.QueryParameters)
			return cliutil.CliQuery(cdc, route, nil)
		},
	}
}

func QueryBancorInfoCmd(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "info [stock] [money]",
		Short: "query the banor pool's information about a symbol pair",
		Long: `query the banor pool's information about a symbol pair. 

Example : 
	cetcli query bancorlite info stock money --trust-node=true --chain-id=myposmyposchain`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			query := fmt.Sprintf("custom/%s/%s", types.StoreKey, keepers.QueryBancorInfo)
			symbol := myposchain.GetSymbol(args[0], args[1])
			param := &keepers.QueryBancorInfoParam{Symbol: symbol}
			return cliutil.CliQuery(cdc, query, param)
		},
	}
}

func QueryBancorListCmd(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "infos",
		Short: "query all bancor infos in blockchain",
		Long: `query all bancor infos in blockchain.

Example :
	cetcli query bancorlite infos \
	--trust-node=true --chain-id=myposmyposchain`,
		RunE: func(cmd *cobra.Command, args []string) error {
			query := fmt.Sprintf("custom/%s/%s", types.StoreKey, keepers.QueryBancors)
			return cliutil.CliQuery(cdc, query, nil)
		},
	}
}
