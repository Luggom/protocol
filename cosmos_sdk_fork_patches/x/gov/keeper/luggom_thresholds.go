package keeper

import (
	"encoding/binary"

	"cosmossdk.io/math"

	v1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
)

// Luggom fork: limiares por assunto.
//
// O x/gov tem um threshold global (params), mas duas matérias exigem
// supermaioria própria:
//   - whitelist de mineração (MsgAddTaxEntity/MsgRemoveTaxEntity): 67%;
//   - custódia do BTC do OTC LOCKED (MsgSpendFromReserve com reserve="otc"): 87%.
//
// A classificação é feita pelo type URL e, no caso do reserve, pelo campo 2 do
// wire format. Não importa os tipos do módulo luggom (evita ciclo de import).
const (
	LuggomWhitelistThreshold  = "0.67"
	LuggomOtcCustodyThreshold = "0.87"

	msgTypeAddTaxEntity       = "/luggom.luggom.v1.MsgAddTaxEntity"
	msgTypeRemoveTaxEntity    = "/luggom.luggom.v1.MsgRemoveTaxEntity"
	msgTypeSpendFromReserve   = "/luggom.luggom.v1.MsgSpendFromReserve"
	spendReserveOtc           = "otc"
	spendFromReserveFieldName = 2 // MsgSpendFromReserve.reserve
)

// LuggomProposalExtraThreshold devolve o maior limiar extra exigido pelas
// mensagens da proposta. Zero quando nenhuma matéria especial está presente.
func LuggomProposalExtraThreshold(proposal v1.Proposal) math.LegacyDec {
	extra := math.LegacyZeroDec()
	raise := func(v math.LegacyDec) {
		if v.GT(extra) {
			extra = v
		}
	}

	for _, msg := range proposal.Messages {
		if msg == nil {
			continue
		}
		switch msg.TypeUrl {
		case msgTypeAddTaxEntity, msgTypeRemoveTaxEntity:
			raise(math.LegacyMustNewDecFromStr(LuggomWhitelistThreshold))
		case msgTypeSpendFromReserve:
			if reserve, ok := luggomStringField(msg.Value, spendFromReserveFieldName); ok && reserve == spendReserveOtc {
				raise(math.LegacyMustNewDecFromStr(LuggomOtcCustodyThreshold))
			}
		}
	}
	return extra
}

// luggomStringField lê um campo string (wire type 2) de uma mensagem protobuf
// sem precisar do tipo concreto. Devolve false se a mensagem estiver malformada
// ou o campo não existir.
func luggomStringField(buf []byte, fieldNum uint64) (string, bool) {
	i := 0
	for i < len(buf) {
		tag, n := binary.Uvarint(buf[i:])
		if n <= 0 {
			return "", false
		}
		i += n
		field := tag >> 3
		wire := tag & 7

		switch wire {
		case 0: // varint
			_, n := binary.Uvarint(buf[i:])
			if n <= 0 {
				return "", false
			}
			i += n
		case 1: // fixed64
			if i+8 > len(buf) {
				return "", false
			}
			i += 8
		case 2: // length-delimited
			l, n := binary.Uvarint(buf[i:])
			if n <= 0 || l > uint64(len(buf)-i-n) {
				return "", false
			}
			i += n
			if field == fieldNum {
				return string(buf[i : i+int(l)]), true
			}
			i += int(l)
		case 5: // fixed32
			if i+4 > len(buf) {
				return "", false
			}
			i += 4
		default:
			return "", false
		}
	}
	return "", false
}
