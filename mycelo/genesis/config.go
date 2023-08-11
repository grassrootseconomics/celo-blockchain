package genesis

import (
	"math/big"

	"github.com/celo-org/celo-blockchain/common"
	"github.com/celo-org/celo-blockchain/common/decimal/bigintstr"
	"github.com/celo-org/celo-blockchain/common/decimal/fixed"
	"github.com/celo-org/celo-blockchain/mycelo/internal/utils"
	"github.com/celo-org/celo-blockchain/params"
	"github.com/shopspring/decimal"
)

// durations in seconds
const (
	Second = 1
	Minute = 60 * Second
	Hour   = 60 * Minute
	Day    = 24 * Hour
	Week   = 7 * Day
	Year   = 365 * Day
)

// Config represent all celo-blockchain configuration options for the genesis block
type Config struct {
	ChainID          *big.Int              `json:"chainId"` // chainId identifies the current chain and is used for replay protection
	Istanbul         params.IstanbulConfig `json:"istanbul"`
	Hardforks        HardforkConfig        `json:"hardforks"`
	GenesisTimestamp uint64                `json:"genesisTimestamp"`

	SortedOracles              SortedOraclesParameters
	GasPriceMinimum            GasPriceMinimumParameters
	Reserve                    ReserveParameters
	StableToken                StableTokenParameters
	StableTokenEUR             StableTokenParameters
	StableTokenBRL             StableTokenParameters
	Exchange                   ExchangeParameters
	ExchangeEUR                ExchangeParameters
	ExchangeBRL                ExchangeParameters
	LockedGold                 LockedGoldParameters
	GoldToken                  GoldTokenParameters
	Validators                 ValidatorsParameters
	Election                   ElectionParameters
	EpochRewards               EpochRewardsParameters
	Blockchain                 BlockchainParameters
	Random                     RandomParameters
	Attestations               AttestationsParameters
	ReserveSpenderMultiSig     MultiSigParameters
	GovernanceApproverMultiSig MultiSigParameters
	DoubleSigningSlasher       DoubleSigningSlasherParameters
	DowntimeSlasher            DowntimeSlasherParameters
	Governance                 GovernanceParameters
	GrandaMento                GrandaMentoParameters
	FeeHandler                 FeeHandlerParameters
}

// Save will write config into a json file
func (cfg *Config) Save(filepath string) error {
	return utils.WriteJson(cfg, filepath)
}

// LoadConfig will read config from a json file
func LoadConfig(filepath string) (*Config, error) {
	var cfg Config
	if err := utils.ReadJson(&cfg, filepath); err != nil {
		return nil, err
	}
	return &cfg, nil
}

// ChainConfig returns the chain config objt for the blockchain
func (cfg *Config) ChainConfig() *params.ChainConfig {
	return &params.ChainConfig{
		ChainID:             cfg.ChainID,
		HomesteadBlock:      common.Big0,
		EIP150Block:         common.Big0,
		EIP150Hash:          common.Hash{},
		EIP155Block:         common.Big0,
		EIP158Block:         common.Big0,
		ByzantiumBlock:      common.Big0,
		ConstantinopleBlock: common.Big0,
		PetersburgBlock:     common.Big0,
		IstanbulBlock:       common.Big0,

		ChurritoBlock:      cfg.Hardforks.ChurritoBlock,
		DonutBlock:         cfg.Hardforks.DonutBlock,
		EspressoBlock:      cfg.Hardforks.EspressoBlock,
		GingerbreadBlock:   cfg.Hardforks.GingerbreadBlock,
		GingerbreadP2Block: cfg.Hardforks.GingerbreadP2Block,

		Istanbul: &params.IstanbulConfig{
			Epoch:          cfg.Istanbul.Epoch,
			ProposerPolicy: cfg.Istanbul.ProposerPolicy,
			LookbackWindow: cfg.Istanbul.LookbackWindow,
			BlockPeriod:    cfg.Istanbul.BlockPeriod,
			RequestTimeout: cfg.Istanbul.RequestTimeout,
		},
	}
}

// HardforkConfig contains celo hardforks activation blocks
type HardforkConfig struct {
	ChurritoBlock      *big.Int `json:"churritoBlock"`
	DonutBlock         *big.Int `json:"donutBlock"`
	EspressoBlock      *big.Int `json:"espressoBlock"`
	GingerbreadBlock   *big.Int `json:"gingerbreadBlock"`
	GingerbreadP2Block *big.Int `json:"gingerbreadP2Block"`
}

// MultiSigParameters are the initial configuration parameters for a MultiSig contract
type MultiSigParameters struct {
	Signatories                      []common.Address `json:"signatories"`
	NumRequiredConfirmations         uint64           `json:"numRequiredConfirmations"`
	NumInternalRequiredConfirmations uint64           `json:"numInternalRequiredConfirmations"`
}

//go:generate gencodec -type LockedGoldRequirements -field-override LockedgoldRequirementsMarshaling -out gen_locked_gold_requirements_json.go

// LockedGoldRequirements represents value/duration requirments on locked gold
type LockedGoldRequirements struct {
	Value    *big.Int `json:"value"`
	Duration uint64   `json:"duration"`
}

type LockedgoldRequirementsMarshaling struct {
	Value *bigintstr.BigIntStr `json:"value"`
}

//go:generate gencodec -type ElectionParameters -field-override ElectionParametersMarshaling -out gen_election_parameters_json.go

// ElectionParameters are the initial configuration parameters for Elections
type ElectionParameters struct {
	MinElectableValidators uint64       `json:"minElectableValidators"`
	MaxElectableValidators uint64       `json:"maxElectableValidators"`
	MaxVotesPerAccount     *big.Int     `json:"maxVotesPerAccount"`
	ElectabilityThreshold  *fixed.Fixed `json:"electabilityThreshold"`
}

type ElectionParametersMarshaling struct {
	MaxVotesPerAccount *bigintstr.BigIntStr `json:"maxVotesPerAccount"`
}

// BlockchainParameters are the initial configuration parameters for Blockchain
type BlockchainParameters struct {
	GasForNonGoldCurrencies uint64 `json:"gasForNonGoldCurrencies"`
	BlockGasLimit           uint64 `json:"blockGasLimit"`
}

//go:generate gencodec -type DoubleSigningSlasherParameters -field-override DoubleSigningSlasherParametersMarshaling -out gen_double_signing_slasher_parameters_json.go

// DoubleSigningSlasherParameters are the initial configuration parameters for DoubleSigningSlasher
type DoubleSigningSlasherParameters struct {
	Penalty *big.Int `json:"penalty"`
	Reward  *big.Int `json:"reward"`
}

type DoubleSigningSlasherParametersMarshaling struct {
	Penalty *bigintstr.BigIntStr `json:"penalty"`
	Reward  *bigintstr.BigIntStr `json:"reward"`
}

//go:generate gencodec -type DowntimeSlasherParameters -field-override DowntimeSlasherParametersMarshaling -out gen_downtime_slasher_parameters_json.go

// DowntimeSlasherParameters are the initial configuration parameters for DowntimeSlasher
type DowntimeSlasherParameters struct {
	Penalty           *big.Int `json:"penalty"`
	Reward            *big.Int `json:"reward"`
	SlashableDowntime uint64   `json:"slashableDowntime"`
}

type DowntimeSlasherParametersMarshaling struct {
	Penalty *bigintstr.BigIntStr `json:"penalty"`
	Reward  *bigintstr.BigIntStr `json:"reward"`
}

//go:generate gencodec -type GovernanceParameters -field-override GovernanceParametersMarshaling -out gen_governance_parameters_json.go

// GovernanceParameters are the initial configuration parameters for Governance
type GovernanceParameters struct {
	UseMultiSig             bool         `json:"useMultiSig"` // whether the approver should be the multisig (otherwise it's the admin)
	ConcurrentProposals     uint64       `json:"concurrentProposals"`
	MinDeposit              *big.Int     `json:"minDeposit"`
	QueueExpiry             uint64       `json:"queueExpiry"`
	DequeueFrequency        uint64       `json:"dequeueFrequency"`
	ReferendumStageDuration uint64       `json:"referendumStageDuration"`
	ExecutionStageDuration  uint64       `json:"executionStageDuration"`
	ParticipationBaseline   *fixed.Fixed `json:"participationBaseline"`
	ParticipationFloor      *fixed.Fixed `json:"participationFloor"`
	BaselineUpdateFactor    *fixed.Fixed `json:"baselineUpdateFactor"`
	BaselineQuorumFactor    *fixed.Fixed `json:"baselineQuorumFactor"`
}

type GovernanceParametersMarshaling struct {
	MinDeposit *bigintstr.BigIntStr `json:"minDeposit"`
}

// ValidatorsParameters are the initial configuration parameters for Validators
type ValidatorsParameters struct {
	GroupLockedGoldRequirements     LockedGoldRequirements `json:"groupLockedGoldRequirements"`
	ValidatorLockedGoldRequirements LockedGoldRequirements `json:"validatorLockedGoldRequirements"`
	ValidatorScoreExponent          uint64                 `json:"validatorScoreExponent"`
	ValidatorScoreAdjustmentSpeed   *fixed.Fixed           `json:"validatorScoreAdjustmentSpeed"`
	MembershipHistoryLength         uint64                 `json:"membershipHistoryLength"`
	SlashingPenaltyResetPeriod      uint64                 `json:"slashingPenaltyResetPeriod"`
	MaxGroupSize                    uint64                 `json:"maxGroupSize"`
	CommissionUpdateDelay           uint64                 `json:"commissionUpdateDelay"`
	DowntimeGracePeriod             uint64                 `json:"downtimeGracePeriod"`

	Commission *fixed.Fixed `json:"commission"` // commission for genesis registered validator groups
}

//go:generate gencodec -type EpochRewardsParameters -field-override EpochRewardsParametersMarshaling -out gen_epoch_rewards_parameters_json.go

// EpochRewardsParameters are the initial configuration parameters for EpochRewards
type EpochRewardsParameters struct {
	TargetVotingYieldInitial                     *fixed.Fixed   `json:"targetVotingYieldInitial"`
	TargetVotingYieldMax                         *fixed.Fixed   `json:"targetVotingYieldMax"`
	TargetVotingYieldAdjustmentFactor            *fixed.Fixed   `json:"targetVotingYieldAdjustmentFactor"`
	RewardsMultiplierMax                         *fixed.Fixed   `json:"rewardsMultiplierMax"`
	RewardsMultiplierAdjustmentFactorsUnderspend *fixed.Fixed   `json:"rewardsMultiplierAdjustmentFactorsUnderspend"`
	RewardsMultiplierAdjustmentFactorsOverspend  *fixed.Fixed   `json:"rewardsMultiplierAdjustmentFactorsOverspend"`
	TargetVotingGoldFraction                     *fixed.Fixed   `json:"targetVotingGoldFraction"`
	MaxValidatorEpochPayment                     *big.Int       `json:"maxValidatorEpochPayment"`
	CommunityRewardFraction                      *fixed.Fixed   `json:"communityRewardFraction"`
	CarbonOffsettingPartner                      common.Address `json:"carbonOffsettingPartner"`
	CarbonOffsettingFraction                     *fixed.Fixed   `json:"carbonOffsettingFraction"`
	Frozen                                       bool           `json:"frozen"`
}

type EpochRewardsParametersMarshaling struct {
	MaxValidatorEpochPayment *bigintstr.BigIntStr `json:"maxValidatorEpochPayment"`
}

// GoldTokenParameters are the initial configuration parameters for GoldToken
type GoldTokenParameters struct {
	Frozen          bool        `json:"frozen"`
	InitialBalances BalanceList `json:"initialBalances"`
}

// RandomParameters are the initial configuration parameters for Random
type RandomParameters struct {
	RandomnessBlockRetentionWindow uint64 `json:"randomnessBlockRetentionWindow"`
}

// AttestationsParameters are the initial configuration parameters for Attestations
type AttestationsParameters struct {
	AttestationExpiryBlocks        uint64          `json:"attestationExpiryBlocks"`
	SelectIssuersWaitBlocks        uint64          `json:"selectIssuersWaitBlocks"`
	MaxAttestations                uint64          `json:"maxAttestations"`
	AttestationRequestFeeInDollars decimal.Decimal `json:"AttestationRequestFeeInDollars"`
}

// SortedOraclesParameters are the initial configuration parameters for SortedOracles
type SortedOraclesParameters struct {
	ReportExpirySeconds uint64 `json:"reportExpirySeconds"`
}

//go:generate gencodec -type GasPriceMinimumParameters -field-override GasPriceMinimumParametersMarshaling -out gen_gas_price_minimum_parameters_json.go

// GasPriceMinimumParameters are the initial configuration parameters for GasPriceMinimum
type GasPriceMinimumParameters struct {
	MinimumFloor                 *big.Int     `json:"minimumFloor"`
	TargetDensity                *fixed.Fixed `json:"targetDensity"`
	AdjustmentSpeed              *fixed.Fixed `json:"adjustmentSpeed"`
	BaseFeeOpCodeActivationBlock *big.Int     `json:"baseFeeOpCodeActivationBlock"`
}

type GasPriceMinimumParametersMarshaling struct {
	MinimumFloor                 *bigintstr.BigIntStr `json:"minimumFloor"`
	BaseFeeOpCodeActivationBlock *bigintstr.BigIntStr `json:"baseFeeOpCodeActivationBlock"`
}

// GrandaMentoParameters are the initial configuration parameters for GrandaMento
type GrandaMentoParameters struct {
	Approver                      common.Address                `json:"approver"`
	MaxApprovalExchangeRateChange *fixed.Fixed                  `json:"maxApprovalExchangeRateChange"`
	Spread                        *fixed.Fixed                  `json:"spread"`
	VetoPeriodSeconds             uint64                        `json:"vetoPeriodSeconds"`
	StableTokenExchangeLimits     StableTokenExchangeLimitsList `json:"stableTokenExchangeLimits"`
}

// FeeHandlerParameters are the initial configuration parameters for FeeHandler
type FeeHandlerParameters struct {
	NewFeeBeneficiary common.Address   `json:"newFeeBeneficiary"`
	NewBurnFraction   *fixed.Fixed     `json:"newBurnFraction"`
	Tokens            []common.Address `json:"tokens"`
	Handlers          []common.Address `json:"handlers"`
	NewLimits         []*big.Int       `json:"newLimits"`
	NewMaxSlippages   []*big.Int       `json:"newMaxSlippages"`
}

//go:generate gencodec -type StableTokenExchangeLimit -field-override StableTokenExchangeLimitsMarshaling -out gen_stable_token_exchange_limit_json.go

// StableTokenExchangeLimit represents the granda mento's exchange limit for a specific stable
type StableTokenExchangeLimit struct {
	StableToken       string   `json:"stableToken"`
	MinExchangeAmount *big.Int `json:"minExchangeAmount"`
	MaxExchangeAmount *big.Int `json:"maxExchangeAmount"`
}

type StableTokenExchangeLimitsMarshaling struct {
	MinExchangeAmount *bigintstr.BigIntStr `json:"minExchangeAmount"`
	MaxExchangeAmount *bigintstr.BigIntStr `json:"maxExchangeAmount"`
}

// BalanceList list of balances
type StableTokenExchangeLimitsList []StableTokenExchangeLimit

//go:generate gencodec -type ReserveParameters -field-override ReserveParametersMarshaling -out gen_reserve_parameters_json.go

// ReserveParameters are the initial configuration parameters for Reserve
type ReserveParameters struct {
	TobinTaxStalenessThreshold uint64              `json:"tobinTaxStalenessThreshold"`
	DailySpendingRatio         *fixed.Fixed        `json:"dailySpendingRatio"`
	AssetAllocations           AssetAllocationList `json:"assetAllocations"`
	TobinTax                   *fixed.Fixed        `json:"tobinTax"`
	TobinTaxReserveRatio       *fixed.Fixed        `json:"tobinTaxReserveRatio"`

	// Other parameters
	Spenders                 []common.Address `json:"spenders"`
	OtherAddresses           []common.Address `json:"otherAddresses"`
	InitialBalance           *big.Int         `json:"initialBalance"`
	FrozenAssetsStartBalance *big.Int         `json:"frozenAssetsStartBalance"`
	FrozenAssetsDays         uint64           `json:"frozenAssetsDays"`
}

type ReserveParametersMarshaling struct {
	InitialBalance           *bigintstr.BigIntStr `json:"initialBalance"`
	FrozenAssetsStartBalance *bigintstr.BigIntStr `json:"frozenAssetsStartBalance"`
}

// StableTokenParameters are the initial configuration parameters for StableToken
type StableTokenParameters struct {
	Name                        string           `json:"name"`
	Symbol                      string           `json:"symbol"`
	Decimals                    uint8            `json:"decimals"`
	Rate                        *fixed.Fixed     `json:"rate"`
	InflationFactorUpdatePeriod uint64           `json:"inflationFactorUpdatePeriod"` // How often the inflation factor is updated.
	InitialBalances             BalanceList      `json:"initialBalances"`
	Frozen                      bool             `json:"frozen"`
	Oracles                     []common.Address `json:"oracles"`
	GoldPrice                   *fixed.Fixed     `json:"goldPrice"`
	ExchangeIdentifier          string           `json:"exchangeIdentifier"`
}

// ExchangeParameters are the initial configuration parameters for Exchange
type ExchangeParameters struct {
	Frozen          bool         `json:"frozen"`
	Spread          *fixed.Fixed `json:"spread"`
	ReserveFraction *fixed.Fixed `json:"reserveFraction"`
	UpdateFrequency uint64       `json:"updateFrequency"`
	MinimumReports  uint64       `json:"minimumReports"`
}

// LockedGoldParameters are the initial configuration parameters for LockedGold
type LockedGoldParameters struct {
	UnlockingPeriod uint64 `json:"unlockingPeriod"`
}

//go:generate gencodec -type Balance -field-override BalanceMarshaling -out gen_balance_json.go

// Balance represents an account and it's initial balance in wei
type Balance struct {
	Account common.Address `json:"account"`
	Amount  *big.Int       `json:"amount"`
}

type BalanceMarshaling struct {
	Amount *bigintstr.BigIntStr `json:"amount"`
}

// BalanceList list of balances
type BalanceList []Balance

// Accounts returns all the addresses
func (bl BalanceList) Accounts() []common.Address {
	res := make([]common.Address, len(bl))
	for i, x := range bl {
		res[i] = x.Account
	}
	return res
}

// Amounts returns all the amounts
func (bl BalanceList) Amounts() []*big.Int {
	res := make([]*big.Int, len(bl))
	for i, x := range bl {
		res[i] = x.Amount
	}
	return res
}

// AssetAllocation config for Reserve
type AssetAllocation struct {
	Symbol string       `json:"symbol"`
	Weight *fixed.Fixed `json:"weight"`
}

// AssetAllocationList list of AssetAllocation
type AssetAllocationList []AssetAllocation

// SymbolsABI returns symbols in ABI format for assets in list
func (aa AssetAllocationList) SymbolsABI() []common.Hash {
	res := make([]common.Hash, len(aa))
	for i, x := range aa {

		res[i] = common.BytesToHash(common.RightPadBytes([]byte(x.Symbol), 32))
	}
	return res
}

// Weights returns weights for assets in list
func (aa AssetAllocationList) Weights() []*big.Int {
	res := make([]*big.Int, len(aa))
	for i, x := range aa {
		res[i] = x.Weight.BigInt()
	}
	return res
}
