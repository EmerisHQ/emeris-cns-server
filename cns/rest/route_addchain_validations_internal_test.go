package rest

import (
	"testing"

	"github.com/allinbits/demeris-backend-models/cns"
	"github.com/stretchr/testify/require"
)

// Test the validations for Fees
func TestValidateFees(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		dataStruct cns.Chain
		success    bool
	}{
		{
			"Validate Fees - No fee tokens",
			cns.Chain{
				Denoms: cns.DenomList{
					cns.Denom{
						Name:     "test",
						FeeToken: false,
					},
				},
			},
			false,
		},
		{
			"Validate Fees - No gas price levels",
			cns.Chain{
				Denoms: cns.DenomList{
					cns.Denom{
						Name:           "test",
						FeeToken:       true,
						GasPriceLevels: cns.GasPrice{},
					},
				},
			},
			false,
		},
		{
			"Validate Fees - Valid",
			cns.Chain{
				Denoms: cns.DenomList{
					cns.Denom{
						Name:     "test",
						FeeToken: true,
						GasPriceLevels: cns.GasPrice{
							Low:     0.3,
							Average: 0.4,
							High:    0.5,
						},
					},
				},
			},
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// act
			err := validateFees(tt.dataStruct)

			// assert
			if !tt.success {
				require.Error(t, err, "Expecting a failed test case")
			} else {
				require.NoError(t, err)
			}
		})
	}
}

// Test the validations for Denoms
func TestValidateDenom(t *testing.T) {
	t.Parallel()

	helperInt64 := int64(24)

	tests := []struct {
		name       string
		dataStruct cns.Chain
		success    bool
	}{
		{
			"Validate Denoms - Multiple relayer denoms",
			cns.Chain{
				Denoms: cns.DenomList{
					cns.Denom{
						Name:         "test",
						FeeToken:     false,
						RelayerDenom: true,
					},
					cns.Denom{
						Name:         "test",
						FeeToken:     false,
						RelayerDenom: true,
					},
				},
			},
			false,
		},
		{
			"Validate Fees - No MinimumThreshRelayerBalance",
			cns.Chain{
				Denoms: cns.DenomList{
					cns.Denom{
						Name:         "test",
						FeeToken:     true,
						RelayerDenom: true,
					},
				},
			},
			false,
		},
		{
			"Validate Fees - Valid",
			cns.Chain{
				Denoms: cns.DenomList{
					cns.Denom{
						Name:                        "test",
						FeeToken:                    true,
						RelayerDenom:                true,
						MinimumThreshRelayerBalance: &helperInt64,
					},
				},
			},
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// act
			err := validateDenoms(tt.dataStruct)

			// assert
			if !tt.success {
				require.Error(t, err, "Expecting a failed test case")
			} else {
				require.NoError(t, err)
			}
		})
	}
}
