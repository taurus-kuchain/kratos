package types

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"bytes"
	"github.com/KuChainNetwork/kuchain/chain/types/coin"
	stakingexport "github.com/KuChainNetwork/kuchain/x/staking/exported"
	"github.com/KuChainNetwork/kuchain/x/staking/external"
	"github.com/cosmos/cosmos-sdk/codec"
	yaml "gopkg.in/yaml.v2"
)

// Staking params default values
const (
	// DefaultUnbondingTime reflects three weeks in seconds as the default
	// unbonding time.
	// TODO: Justify our choice of default here.
	DefaultUnbondingTime time.Duration = time.Hour * 24 * 7 * 2

	// Default maximum number of bonded validators
	DefaultMaxValidators uint32 = 33

	// Default maximum entries in a UBD/RED pair
	DefaultMaxEntries uint32 = 7

	// DefaultHistorical entries is 0 since it must only be non-zero for
	// IBC connected chains
	DefaultHistoricalEntries uint32 = 0
)

// nolint - Keys for parameter access
var (
	KeyUnbondingTime     = []byte("UnbondingTime")
	KeyMaxValidators     = []byte("MaxValidators")
	KeyMaxEntries        = []byte("KeyMaxEntries")
	KeyBondDenom         = []byte("BondDenom")
	KeyHistoricalEntries = []byte("HistoricalEntries")
)

var _ external.ParamsSet = (*Params)(nil)

// Params defines the parameters for the staking module.
type Params struct {
	UnbondingTime     time.Duration `json:"unbonding_time" yaml:"unbonding_time"`
	MaxValidators     uint32        `json:"max_validators,omitempty" yaml:"max_validators"`
	MaxEntries        uint32        `json:"max_entries,omitempty" yaml:"max_entries"`
	HistoricalEntries uint32        `json:"historical_entries,omitempty" yaml:"historical_entries"`
	BondDenom         string        `json:"bond_denom,omitempty" yaml:"bond_denom"`
}

// NewParams creates a new Params instance
func NewParams(
	unbondingTime time.Duration, maxValidators, maxEntries, historicalEntries uint32, bondDenom string,
) Params {

	return Params{
		UnbondingTime:     unbondingTime,
		MaxValidators:     maxValidators,
		MaxEntries:        maxEntries,
		HistoricalEntries: historicalEntries,
		BondDenom:         bondDenom,
	}
}

// Implements params.ParamSet
func (p *Params) ParamSetPairs() external.ParamsSetPairs {
	return external.ParamsSetPairs{
		external.NewParamSetPair(KeyUnbondingTime, &p.UnbondingTime, validateUnbondingTime),
		external.NewParamSetPair(KeyMaxValidators, &p.MaxValidators, validateMaxValidators),
		external.NewParamSetPair(KeyMaxEntries, &p.MaxEntries, validateMaxEntries),
		external.NewParamSetPair(KeyHistoricalEntries, &p.HistoricalEntries, validateHistoricalEntries),
		external.NewParamSetPair(KeyBondDenom, &p.BondDenom, validateBondDenom),
	}
}

// DefaultParams returns a default set of parameters.
func DefaultParams() Params {
	return NewParams(
		DefaultUnbondingTime,
		DefaultMaxValidators,
		DefaultMaxEntries,
		DefaultHistoricalEntries,
		stakingexport.DefaultBondDenom,
	)
}

// String returns a human readable string representation of the parameters.
func (p Params) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}

// unmarshal the current staking params value from store key or panic
func MustUnmarshalParams(cdc *codec.Codec, value []byte) Params {
	params, err := UnmarshalParams(cdc, value)
	if err != nil {
		panic(err)
	}
	return params
}

// unmarshal the current staking params value from store key
func UnmarshalParams(cdc *codec.Codec, value []byte) (params Params, err error) {
	err = cdc.UnmarshalBinaryBare(value, &params)
	if err != nil {
		return
	}
	return
}

// validate a set of params
func (p Params) Validate() error {
	if err := validateUnbondingTime(p.UnbondingTime); err != nil {
		return err
	}
	if err := validateMaxValidators(p.MaxValidators); err != nil {
		return err
	}
	if err := validateMaxEntries(p.MaxEntries); err != nil {
		return err
	}
	if err := validateBondDenom(p.BondDenom); err != nil {
		return err
	}

	return nil
}

func validateUnbondingTime(i interface{}) error {
	v, ok := i.(time.Duration)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v <= 0 {
		return fmt.Errorf("unbonding time must be positive: %d", v)
	}

	return nil
}

func validateMaxValidators(i interface{}) error {
	v, ok := i.(uint32)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v == 0 {
		return fmt.Errorf("max validators must be positive: %d", v)
	}

	return nil
}

func validateMaxEntries(i interface{}) error {
	v, ok := i.(uint32)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v == 0 {
		return fmt.Errorf("max entries must be positive: %d", v)
	}

	return nil
}

func validateHistoricalEntries(i interface{}) error {
	_, ok := i.(uint32)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	return nil
}

func validateBondDenom(i interface{}) error {
	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if strings.TrimSpace(v) == "" {
		return errors.New("bond denom cannot be blank")
	}
	if err := coin.ValidateDenom(v); err != nil {
		return err
	}

	return nil
}

// Equal returns a boolean determining if two Param types are identical.
// TODO: This is slower than comparing struct fields directly
func (p Params) Equal(p2 Params) bool {
	bz1 := ModuleCdc.MustMarshalBinaryLengthPrefixed(&p)
	bz2 := ModuleCdc.MustMarshalBinaryLengthPrefixed(&p2)
	return bytes.Equal(bz1, bz2)
}
