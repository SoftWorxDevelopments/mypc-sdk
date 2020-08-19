package cli

import (
	"fmt"
	"testing"

	"github.com/spf13/cobra"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/SoftWorxDevelopments/mypc-sdk/modules/distributionx/types"
	myposchain "github.com/SoftWorxDevelopments/mypc-sdk/types"
	"github.com/coinexchain/cosmos-utils/client/cliutil"
)

var testAddrBech32 = "mypos12kcupm2x8fw0gglgcz8850kw0k2kx0ff8sr3rn"

func TestDonateTxCmd(t *testing.T) {
	cmdFactory := func() *cobra.Command {
		return DonateTxCmd(nil)
	}

	testAddr, _ := sdk.AccAddressFromBech32(testAddrBech32)
	args := fmt.Sprintf("1000mypc --from=%s", testAddr)
	msg := types.NewMsgDonateToCommunityPool(testAddr, myposchain.NewMypcCoins(1000))

	cliutil.TestTxCmd(t, cmdFactory, args, &msg)
}
