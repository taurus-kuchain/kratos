package simulation

// DONTCOVER

import (
	"fmt"
	"math/rand"

	"github.com/KuChain-io/kuchain/x/staking/types"
	"github.com/cosmos/cosmos-sdk/x/simulation"
)

const (
	keyMaxValidators = "MaxValidators"
	keyUnbondingTime = "UnbondingTime"
)

// ParamChanges defines the parameters that can be modified by param change proposals
// on the simulation
func ParamChanges(r *rand.Rand) []simulation.ParamChange {
	return []simulation.ParamChange{
		simulation.NewSimParamChange(types.ModuleName, keyMaxValidators,
			func(r *rand.Rand) string {
				return fmt.Sprintf("%d", GenMaxValidators(r))
			},
		),
		simulation.NewSimParamChange(types.ModuleName, keyUnbondingTime,
			func(r *rand.Rand) string {
				return fmt.Sprintf("\"%d\"", GenUnbondingTime(r))
			},
		),
	}
}
