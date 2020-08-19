package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestErrMsg(t *testing.T) {
	err := ErrMemoMissing()
	require.Equal(t, CodeMemoMissing, err.Code())
	err = ErrInsufficientMYPCForActivatingFee()
	require.Equal(t, CodeInsufficientMYPCForActivationFee, err.Code())
	err = ErrUnlockTime("")
	require.Equal(t, CodeInvalidUnlockTime, err.Code())
	err = ErrTokenForbiddenByOwner()
	require.Equal(t, CodeTokenForbiddenByOwner, err.Code())
}
