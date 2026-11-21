package types

import (
	"strings"
	"testing"
)

// §13.170 A–E: os params de comissão/staleness/flags foram apagados ou
// congelados em código (frozen.go). Este teste cobre o que RESTA.

func TestDefaultParams_GovWeights(t *testing.T) {
	params := DefaultParams()

	if params.GovValidatorMinForWeightShift != 100 {
		t.Fatalf("expected default gov_validator_min_for_weight_shift=100, got %d", params.GovValidatorMinForWeightShift)
	}
	// DECISÃO DO FUNDADOR (2026-09-18): pré-pool o Fundador decide sozinho.
	if params.GovLlcVoteWeightBps != 0 || params.GovFounderVoteWeightBps != 10000 || params.GovValidatorsVoteWeightBps != 0 {
		t.Fatalf("unexpected default gov BPS: %d/%d/%d", params.GovLlcVoteWeightBps, params.GovFounderVoteWeightBps, params.GovValidatorsVoteWeightBps)
	}
	if params.UsdOracleUpdaterAddress != "" {
		t.Fatalf("expected default usd_oracle_updater_address empty")
	}
}

func TestParamsValidate_GovBpsSum(t *testing.T) {
	params := DefaultParams()
	params.GovValidatorsVoteWeightBps = 1
	if err := params.Validate(); err == nil {
		t.Fatalf("expected validation error when gov BPS sum != 10000")
	}
}

func TestParamsValidate_PoolAssets(t *testing.T) {
	validIBC := "ibc/" + strings.Repeat("A", 64)

	cases := []struct {
		name    string
		assets  []string
		wantErr bool
	}{
		{"default", DefaultPoolAssets(), false},
		{"empty", nil, true},
		{"duplicate", []string{"btc", "btc"}, true},
		{"base denom", []string{"uluggom"}, true},
		{"invalid denom", []string{"a"}, true},
		{"malformed ibc", []string{"ibc/abc"}, true},
		{"valid ibc", []string{"btc", validIBC}, false},
		{"whitespace", []string{" btc"}, true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			params := DefaultParams()
			params.PoolAssets = tc.assets
			err := params.Validate()
			if tc.wantErr && err == nil {
				t.Fatalf("expected error, got nil")
			}
			if !tc.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestParamsValidate_TaxEntityWhitelist(t *testing.T) {
	cases := []struct {
		name    string
		list    []string
		wantErr bool
	}{
		{"default", DefaultGlobalTaxEntities(), false},
		{"empty", nil, true},
		{"out of BCA", []string{"de:finanzamt"}, true},
		{"bare prefix", []string{"br:"}, true},
		{"uppercase", []string{"BR:IPTU_X"}, true},
		{"whitespace", []string{" br:iptu_x"}, true},
		{"duplicate", []string{"br:iptu_x", "br:iptu_x"}, true},
		{"path in slug", []string{"br:iptu_x/../y"}, true},
		{"valid mixed", []string{"br:iptu_x", "us:county:cook_il", "ca:city:toronto", "irs"}, false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			params := DefaultParams()
			params.TaxEntityWhitelist = tc.list
			err := params.Validate()
			if tc.wantErr && err == nil {
				t.Fatalf("expected error, got nil")
			}
			if !tc.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestIsBCASlug(t *testing.T) {
	cases := map[string]bool{
		"irs":                  true,
		"ca_ftb":               true,
		"ny_dtf":               true,
		"br:iptu_sao_paulo":    true,
		"us:county:cook_il":    true,
		"ca:city:toronto":      true,
		"":                     false,
		"br:":                  false,
		"us:":                  false,
		"de:finanzamt":         false,
		"br:iptu_x/../y":       false,
		"br:iptu x":            false,
		"br:iptu_x?redirect=1": false,
	}

	for slug, want := range cases {
		if got := IsBCASlug(slug); got != want {
			t.Errorf("IsBCASlug(%q) = %v, want %v", slug, got, want)
		}
	}
}

func TestFrozenConstants(t *testing.T) {
	if TlsnotaryPhaseFrozen != 2 {
		t.Fatalf("tlsnotary_phase deve estar congelado em 2 (não há fase 1)")
	}
	if !DynamicUsdOracleEnabledFrozen {
		t.Fatalf("o oráculo dinâmico está congelado ligado")
	}
	if UsdOracleMaxStaleSecondsFrozen != 129_600 {
		t.Fatalf("staleness congelada em 36 h")
	}
	if TaxPaymentMaxAgeDaysFrozen != 365 {
		t.Fatalf("janela do comprovante congelada em 365 dias")
	}
}
