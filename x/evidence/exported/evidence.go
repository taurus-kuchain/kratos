package exported

import (
	"github.com/KuChain-io/kuchain/chain/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	tmbytes "github.com/tendermint/tendermint/libs/bytes"
)

// Evidence defines the contract which concrete evidence types of misbehavior
// must implement.
type Evidence interface {
	Route() string
	Type() string
	String() string
	Hash() tmbytes.HexBytes
	ValidateBasic() error

	// The consensus address of the malicious validator at time of infraction
	GetConsensusAddress() sdk.ConsAddress

	// Height at which the infraction occurred
	GetHeight() int64

	// The total power of the malicious validator at time of infraction
	GetValidatorPower() int64

	// The total validator set power at time of infraction
	GetTotalPower() int64
}

// MsgSubmitEvidence defines the specific interface a concrete message must
// implement in order to process submitted evidence. The concrete MsgSubmitEvidence
// must be defined at the application-level.
type MsgSubmitEvidence interface {
	sdk.Msg

	GetEvidence() Evidence
	GetSubmitter() types.AccountID
}
