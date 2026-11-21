package govtally_test

import (
	"testing"

	"cosmossdk.io/math"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	govkeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
	v1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	"github.com/stretchr/testify/require"
)

func anyMsg(typeURL string, value []byte) *codectypes.Any {
	return &codectypes.Any{TypeUrl: typeURL, Value: value}
}

// reserveField is the wire encoding of MsgSpendFromReserve.reserve (field 2).
func reserveField(v string) []byte {
	b := []byte{0x12, byte(len(v))}
	return append(b, []byte(v)...)
}

func TestLuggomProposalExtraThreshold(t *testing.T) {
	zero := math.LegacyZeroDec()
	whitelist := math.LegacyMustNewDecFromStr(govkeeper.LuggomWhitelistThreshold)
	custody := math.LegacyMustNewDecFromStr(govkeeper.LuggomOtcCustodyThreshold)

	cases := []struct {
		name string
		msgs []*codectypes.Any
		want math.LegacyDec
	}{
		{"none", nil, zero},
		{"other message", []*codectypes.Any{anyMsg("/cosmos.bank.v1beta1.MsgSend", nil)}, zero},
		{"add whitelist", []*codectypes.Any{anyMsg("/luggom.luggom.v1.MsgAddTaxEntity", nil)}, whitelist},
		{"remove whitelist", []*codectypes.Any{anyMsg("/luggom.luggom.v1.MsgRemoveTaxEntity", nil)}, whitelist},
		{"reserve gov", []*codectypes.Any{anyMsg("/luggom.luggom.v1.MsgSpendFromReserve", reserveField("gov"))}, zero},
		{"reserve notary", []*codectypes.Any{anyMsg("/luggom.luggom.v1.MsgSpendFromReserve", reserveField("notary"))}, zero},
		{"reserve otc", []*codectypes.Any{anyMsg("/luggom.luggom.v1.MsgSpendFromReserve", reserveField("otc"))}, custody},
		{
			"highest wins",
			[]*codectypes.Any{
				anyMsg("/luggom.luggom.v1.MsgAddTaxEntity", nil),
				anyMsg("/luggom.luggom.v1.MsgSpendFromReserve", reserveField("otc")),
			},
			custody,
		},
		{"malformed", []*codectypes.Any{anyMsg("/luggom.luggom.v1.MsgSpendFromReserve", []byte{0xff})}, zero},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := govkeeper.LuggomProposalExtraThreshold(v1.Proposal{Messages: tc.msgs})
			require.True(t, got.Equal(tc.want), "got %s, want %s", got, tc.want)
		})
	}
}
