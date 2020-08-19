package types

import (
	"strings"

	myposchain "github.com/SoftWorxDevelopments/mypc-sdk/types"
)

func IsOnlyForMypos(alias string) bool {
	if strings.HasPrefix(alias, "mypos") ||
		strings.HasSuffix(alias, "mypos") ||
		strings.HasSuffix(alias, "mypos.org") ||
		strings.HasSuffix(alias, "mypos.com") ||
		strings.HasSuffix(alias, "mypos.net") {
		return true
	}

	return alias == myposchain.MYPC || alias == "viabtc" || alias == "mypcdac"
}

func IsValidAlias(alias string) bool {
	if len(alias) < 2 || len(alias) > 45 {
		return false
	}
	for _, c := range alias {
		if !isValidChar(c) {
			return false
		}
	}
	return true
}

func isValidChar(c rune) bool {
	if '0' <= c && c <= '9' {
		return true
	}
	if 'a' <= c && c <= 'z' {
		return true
	}
	if c == '-' || c == '_' || c == '.' || c == '@' {
		return true
	}
	return false
}
