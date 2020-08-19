package distributionx

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/SoftWorxDevelopments/mypc-sdk/modules/distributionx/types"
)

func HandleMsgDonateToCommunityPool(ctx sdk.Context, k Keeper, msg types.MsgDonateToCommunityPool) sdk.Result {
	return handleMsgDonateToCommunityPool(ctx, k, msg)
}
