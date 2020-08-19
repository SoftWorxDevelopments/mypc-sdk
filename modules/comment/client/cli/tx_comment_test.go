package cli

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"

	"github.com/SoftWorxDevelopments/mypc-sdk/modules/comment/internal/types"
	"github.com/coinexchain/cosmos-utils/client/cliutil"
)

var ResultMsg *types.MsgCommentToken

func CliRunCommandForTest(cdc *codec.Codec, msg cliutil.MsgWithAccAddress) error {
	cliCtx := context.NewCLIContext().WithCodec(cdc)
	senderAddr := cliCtx.GetFromAddress()
	msg.SetAccAddress(senderAddr)
	if err := msg.ValidateBasic(); err != nil {
		return err
	}
	ResultMsg = msg.(*types.MsgCommentToken)
	return nil
}

func Test1(t *testing.T) {
	cliutil.CliRunCommand = CliRunCommandForTest

	sdk.GetConfig().SetBech32PrefixForAccount("mypos", "mypospub")
	cmd := GetTxCmd(nil)

	addr, _ := sdk.AccAddressFromHex("01234567890123456789012345678901234abcde")
	addrStr := addr.String()

	args := []string{
		"new-thread",
		"--token=mypc",
		"--donation=2000000",
		`--title=I love mypc.`,
		`--content=MYPC to da moon!!!`,
		"--content-type=UTF8Text",
		"--from=" + addrStr,
		"--generate-only",
	}
	cmd.SetArgs(args)
	cliutil.SetViperWithArgs(args)
	err := cmd.Execute()
	assert.Equal(t, nil, err)
	correct, _ := json.Marshal(&types.MsgCommentToken{
		Sender:      addr,
		Token:       "mypc",
		Donation:    2000000,
		Title:       "I love mypc.",
		Content:     []byte("MYPC to da moon!!!"),
		ContentType: types.UTF8Text,
	})
	msgStr, _ := json.Marshal(ResultMsg)
	assert.Equal(t, string(correct), string(msgStr))

	args = []string{
		"follow-up",
		"--token=mypc",
		"--donation=0",
		`--title=I love mypc too.`,
		`--content=MYPC to da mars!!!`,
		`--follow=10001;mypos1px8alypku5j84qlwzdpynhn4nyrkagaytu5u4a;2;mypc;like,favorite`,
		"--content-type=UTF8Text",
		"--from=" + addrStr,
		"--generate-only",
	}
	cmd.SetArgs(args)
	cliutil.SetViperWithArgs(args)
	err = cmd.Execute()
	assert.Equal(t, nil, err)
	addr1, _ := sdk.AccAddressFromBech32("mypos1px8alypku5j84qlwzdpynhn4nyrkagaytu5u4a")
	correct, _ = json.Marshal(&types.MsgCommentToken{
		Sender:      addr,
		Token:       "mypc",
		Donation:    0,
		Title:       "I love mypc too.",
		Content:     []byte("MYPC to da mars!!!"),
		ContentType: types.UTF8Text,
		References: []types.CommentRef{
			{
				ID:           10001,
				RewardTarget: addr1,
				RewardToken:  "mypc",
				RewardAmount: 2,
				Attitudes:    []int32{types.Like, types.Favorite},
			},
		},
	})
	msgStr, _ = json.Marshal(ResultMsg)
	assert.Equal(t, string(correct), string(msgStr))

	args = []string{
		"reward-comments",
		"--token=mypc",
		`--reward-to=10001;mypos1px8alypku5j84qlwzdpynhn4nyrkagaytu5u4a;2;mypc;like,favorite`,
		`--reward-to=20021;mypos1jv65s3grqf6v6jl3dp4t6c9t9rk99cd8vc4efa;1;mypc;like`,
		"--from=" + addrStr,
		"--generate-only",
	}
	cmd.SetArgs(args)
	cliutil.SetViperWithArgs(args)
	err = cmd.Execute()
	assert.Equal(t, nil, err)
	addr2, _ := sdk.AccAddressFromBech32("mypos1jv65s3grqf6v6jl3dp4t6c9t9rk99cd8vc4efa")
	correct, _ = json.Marshal(&types.MsgCommentToken{
		Sender:      addr,
		Token:       "mypc",
		Donation:    0,
		Title:       "",
		Content:     []byte("No-Content"),
		ContentType: types.UTF8Text,
		References: []types.CommentRef{
			{
				ID:           10001,
				RewardTarget: addr1,
				RewardToken:  "mypc",
				RewardAmount: 2,
				Attitudes:    []int32{types.Like, types.Favorite},
			},
			{
				ID:           20021,
				RewardTarget: addr2,
				RewardToken:  "mypc",
				RewardAmount: 1,
				Attitudes:    []int32{types.Like},
			},
		},
	})
	msgStr, _ = json.Marshal(ResultMsg)
	assert.Equal(t, string(correct), string(msgStr))

	args = []string{
		"reward-comments",
		fmt.Sprintf("--token=%s", "mypc"),
		`--reward-to=10001;mypos1px8alypku5j84qlwzdpynhn4nyrkagaytu5u4a;2;mypc;like,favorite`,
		"--from=" + addrStr,
		"--generate-only",
	}
	cmd.SetArgs(args)
	cliutil.SetViperWithArgs(args)
	err = cmd.Execute()
	assert.Equal(t, nil, err)
	correct, _ = json.Marshal(&types.MsgCommentToken{
		Sender:      addr,
		Token:       "mypc",
		Donation:    0,
		Title:       "reward-comments",
		Content:     []byte("No-Content"),
		ContentType: types.UTF8Text,
		References: []types.CommentRef{
			{
				ID:           10001,
				RewardTarget: addr1,
				RewardToken:  "mypc",
				RewardAmount: 2,
				Attitudes:    []int32{types.Like, types.Favorite},
			},
		},
	})
	msgStr, _ = json.Marshal(ResultMsg)
	assert.Equal(t, string(correct), string(msgStr))

	args = []string{
		"new-thread",
		"--token=mypc",
		"--donation=2000000",
		`--title=I love mypc.`,
		`--content=MYPC to da moon!!!`,
		"--content-type=Haha",
		"--from=" + addrStr,
		"--generate-only",
	}
	cmd.SetArgs(args)
	cliutil.SetViperWithArgs(args)
	err = cmd.Execute()
	errStr := "tx flag is error (Haha is not a valid content type.), please see help : $ cetcli tx comment new-thread -h"
	assert.Equal(t, errStr, err.Error())

	args = []string{
		"reward-comments",
		fmt.Sprintf("--token=%s", "mypc"),
		`--reward-to=10001;2;mypc;like,favorite`,
		"--from=" + addrStr,
		"--generate-only",
	}
	cmd.SetArgs(args)
	cliutil.SetViperWithArgs(args)
	err = cmd.Execute()
	errStr = "tx flag is error (invalid format: 10001;2;mypc;like,favorite), please see help : $ cetcli tx comment reward-comments -h"
	assert.Equal(t, errStr, err.Error())

	args = []string{
		"reward-comments",
		fmt.Sprintf("--token=%s", "mypc"),
		`--reward-to=1a0001;mypos1px8alypku5j84qlwzdpynhn4nyrkagaytu5u4a;2;mypc;like,favorite`,
		"--from=" + addrStr,
		"--generate-only",
	}
	cmd.SetArgs(args)
	cliutil.SetViperWithArgs(args)
	err = cmd.Execute()
	errStr = "tx flag is error (Not a valid comment id: 1a0001), please see help : $ cetcli tx comment reward-comments -h"
	assert.Equal(t, errStr, err.Error())

	args = []string{
		"reward-comments",
		fmt.Sprintf("--token=%s", "mypc"),
		`--reward-to=10001;mypos1px8alypku5j84qlwzdp;2;mypc;like,favorite`,
		"--from=" + addrStr,
		"--generate-only",
	}
	cmd.SetArgs(args)
	cliutil.SetViperWithArgs(args)
	err = cmd.Execute()
	errStr = "tx flag is error (Not a valid address: mypos1px8alypku5j84qlwzdp), please see help : $ cetcli tx comment reward-comments -h"
	assert.Equal(t, errStr, err.Error())

	args = []string{
		"reward-comments",
		fmt.Sprintf("--token=%s", "mypc"),
		`--reward-to=10001;mypos1px8alypku5j84qlwzdpynhn4nyrkagaytu5u4a;2a;mypc;like,favorite`,
		"--from=" + addrStr,
		"--generate-only",
	}
	cmd.SetArgs(args)
	cliutil.SetViperWithArgs(args)
	err = cmd.Execute()
	errStr = "tx flag is error (Not a valid amount: 2a), please see help : $ cetcli tx comment reward-comments -h"
	assert.Equal(t, errStr, err.Error())

	args = []string{
		"reward-comments",
		fmt.Sprintf("--token=%s", "mypc"),
		`--reward-to=10001;mypos1px8alypku5j84qlwzdpynhn4nyrkagaytu5u4a;2;mypc;like,fuck`,
		"--from=" + addrStr,
		"--generate-only",
	}
	cmd.SetArgs(args)
	cliutil.SetViperWithArgs(args)
	err = cmd.Execute()
	errStr = "tx flag is error (invalid attitude: fuck), please see help : $ cetcli tx comment reward-comments -h"
	assert.Equal(t, errStr, err.Error())

	args = []string{
		"follow-up",
		"--token=mypc",
		"--donation=0",
		`--title=I love mypc too.`,
		`--content=MYPC to da mars!!!`,
		`--follow=10001;mypos1px8alypku5j8;2;mypc;like,favorite`,
		"--content-type=UTF8Text",
		"--from=" + addrStr,
		"--generate-only",
	}
	cmd.SetArgs(args)
	cliutil.SetViperWithArgs(args)
	err = cmd.Execute()
	errStr = "tx flag is error (Not a valid address: mypos1px8alypku5j8), please see help : $ cetcli tx comment follow-up -h"
	assert.Equal(t, errStr, err.Error())

	//fmt.Printf("|%s\n", err.Error())
}
