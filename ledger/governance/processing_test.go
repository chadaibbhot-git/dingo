// Copyright 2026 Blink Labs Software
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package governance

import (
	"io"
	"log/slog"
	"math/big"
	"os"
	"testing"

	"github.com/blinklabs-io/dingo/database"
	"github.com/blinklabs-io/gouroboros/cbor"
	gledger "github.com/blinklabs-io/gouroboros/ledger"
	lcommon "github.com/blinklabs-io/gouroboros/ledger/common"
	"github.com/blinklabs-io/gouroboros/ledger/conway"
	ocommon "github.com/blinklabs-io/gouroboros/protocol/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/utxorpc/go-codegen/utxorpc/v1alpha/cardano"

	"github.com/blinklabs-io/dingo/database/models"
)

type mockGovernanceTransaction struct {
	hash               lcommon.Blake2b256
	txType             int
	votingProcedures   lcommon.VotingProcedures
	proposalProcedures []lcommon.ProposalProcedure
	isValid            bool
}

func (m *mockGovernanceTransaction) Hash() lcommon.Blake2b256 {
	return m.hash
}

func (m *mockGovernanceTransaction) Id() lcommon.Blake2b256 {
	return m.hash
}

func (m *mockGovernanceTransaction) Type() int {
	return m.txType
}

func (m *mockGovernanceTransaction) Fee() *big.Int {
	return big.NewInt(0)
}

func (m *mockGovernanceTransaction) TTL() uint64 {
	return 0
}

func (m *mockGovernanceTransaction) IsValid() bool {
	return m.isValid
}

func (m *mockGovernanceTransaction) Metadata() lcommon.TransactionMetadatum {
	return nil
}

func (m *mockGovernanceTransaction) AuxiliaryData() lcommon.AuxiliaryData {
	return nil
}

func (m *mockGovernanceTransaction) RawAuxiliaryData() []byte {
	return nil
}

func (m *mockGovernanceTransaction) CollateralReturn() lcommon.TransactionOutput {
	return nil
}

func (m *mockGovernanceTransaction) Produced() []lcommon.Utxo {
	return nil
}

func (m *mockGovernanceTransaction) Outputs() []lcommon.TransactionOutput {
	return nil
}

func (m *mockGovernanceTransaction) Inputs() []lcommon.TransactionInput {
	return nil
}

func (m *mockGovernanceTransaction) Collateral() []lcommon.TransactionInput {
	return nil
}

func (m *mockGovernanceTransaction) Certificates() []lcommon.Certificate {
	return nil
}

func (m *mockGovernanceTransaction) ProtocolParameterUpdates() (
	uint64,
	map[lcommon.Blake2b224]lcommon.ProtocolParameterUpdate,
) {
	return 0, nil
}

func (m *mockGovernanceTransaction) AssetMint() *lcommon.MultiAsset[lcommon.MultiAssetTypeMint] {
	return nil
}

func (m *mockGovernanceTransaction) AuxDataHash() *lcommon.Blake2b256 {
	return nil
}

func (m *mockGovernanceTransaction) Cbor() []byte {
	return nil
}

func (m *mockGovernanceTransaction) Consumed() []lcommon.TransactionInput {
	return nil
}

func (m *mockGovernanceTransaction) Witnesses() lcommon.TransactionWitnessSet {
	return nil
}

func (m *mockGovernanceTransaction) ValidityIntervalStart() uint64 {
	return 0
}

func (m *mockGovernanceTransaction) ReferenceInputs() []lcommon.TransactionInput {
	return nil
}

func (m *mockGovernanceTransaction) TotalCollateral() *big.Int {
	return nil
}

func (m *mockGovernanceTransaction) Withdrawals() map[*lcommon.Address]*big.Int {
	return nil
}

func (m *mockGovernanceTransaction) RequiredSigners() []lcommon.Blake2b224 {
	return nil
}

func (m *mockGovernanceTransaction) ScriptDataHash() *lcommon.Blake2b256 {
	return nil
}

func (m *mockGovernanceTransaction) VotingProcedures() lcommon.VotingProcedures {
	return m.votingProcedures
}

func (m *mockGovernanceTransaction) ProposalProcedures() []lcommon.ProposalProcedure {
	return m.proposalProcedures
}

func (m *mockGovernanceTransaction) CurrentTreasuryValue() *big.Int {
	return nil
}

func (m *mockGovernanceTransaction) Donation() *big.Int {
	return nil
}

func (m *mockGovernanceTransaction) Utxorpc() (*cardano.Tx, error) {
	return nil, nil
}

func (m *mockGovernanceTransaction) LeiosHash() lcommon.Blake2b256 {
	return lcommon.Blake2b256{}
}

func testHash32(seed string) []byte {
	hash := make([]byte, 32)
	copy(hash, []byte(seed))
	return hash
}

func testHash28(seed string) []byte {
	hash := make([]byte, 28)
	copy(hash, []byte(seed))
	return hash
}

func testConwayProtocolParameters() *conway.ConwayProtocolParameters {
	return &conway.ConwayProtocolParameters{
		MinFeeA:            44,
		MinFeeB:            155381,
		MaxBlockBodySize:   65536,
		MaxTxSize:          16384,
		MaxBlockHeaderSize: 1100,
		KeyDeposit:         2000000,
		PoolDeposit:        500000000,
		MaxEpoch:           18,
		NOpt:               500,
		A0:                 &cbor.Rat{Rat: big.NewRat(3, 10)},
		Rho:                &cbor.Rat{Rat: big.NewRat(3, 1000)},
		Tau:                &cbor.Rat{Rat: big.NewRat(1, 5)},
		ProtocolVersion: lcommon.ProtocolParametersProtocolVersion{
			Major: 10,
			Minor: 0,
		},
		MinPoolCost:    340000000,
		AdaPerUtxoByte: 4310,
		CostModels:     map[uint][]int64{1: {0}},
		ExecutionCosts: lcommon.ExUnitPrice{
			MemPrice:  &cbor.Rat{Rat: big.NewRat(577, 10000)},
			StepPrice: &cbor.Rat{Rat: big.NewRat(721, 10000000)},
		},
		MaxTxExUnits: lcommon.ExUnits{
			Memory: 10000000,
			Steps:  10000000000,
		},
		MaxBlockExUnits: lcommon.ExUnits{
			Memory: 50000000,
			Steps:  40000000000,
		},
		MaxValueSize:         5000,
		CollateralPercentage: 150,
		MaxCollateralInputs:  3,
		PoolVotingThresholds: conway.PoolVotingThresholds{
			MotionNoConfidence:    cbor.Rat{Rat: big.NewRat(1, 2)},
			CommitteeNormal:       cbor.Rat{Rat: big.NewRat(1, 2)},
			CommitteeNoConfidence: cbor.Rat{Rat: big.NewRat(1, 2)},
			HardForkInitiation:    cbor.Rat{Rat: big.NewRat(1, 2)},
			PpSecurityGroup:       cbor.Rat{Rat: big.NewRat(1, 2)},
		},
		DRepVotingThresholds: conway.DRepVotingThresholds{
			MotionNoConfidence:    cbor.Rat{Rat: big.NewRat(1, 2)},
			CommitteeNormal:       cbor.Rat{Rat: big.NewRat(1, 2)},
			CommitteeNoConfidence: cbor.Rat{Rat: big.NewRat(1, 2)},
			UpdateToConstitution:  cbor.Rat{Rat: big.NewRat(1, 2)},
			HardForkInitiation:    cbor.Rat{Rat: big.NewRat(1, 2)},
			PpNetworkGroup:        cbor.Rat{Rat: big.NewRat(1, 2)},
			PpEconomicGroup:       cbor.Rat{Rat: big.NewRat(1, 2)},
			PpTechnicalGroup:      cbor.Rat{Rat: big.NewRat(1, 2)},
			PpGovGroup:            cbor.Rat{Rat: big.NewRat(1, 2)},
			TreasuryWithdrawal:    cbor.Rat{Rat: big.NewRat(1, 2)},
		},
		MinCommitteeSize:        5,
		CommitteeTermLimit:      146,
		GovActionValidityPeriod: 20,
		GovActionDeposit:        100000000000,
		DRepDeposit:             500000000,
		DRepInactivityPeriod:    20,
		MinFeeRefScriptCostPerByte: &cbor.Rat{
			Rat: big.NewRat(1, 1),
		},
	}
}

func TestMapVoterType(t *testing.T) {
	tests := []struct {
		name     string
		input    uint8
		expected uint8
	}{
		{
			name:     "CC hot key hash",
			input:    lcommon.VoterTypeConstitutionalCommitteeHotKeyHash,
			expected: models.VoterTypeCC,
		},
		{
			name:     "CC hot script hash",
			input:    lcommon.VoterTypeConstitutionalCommitteeHotScriptHash,
			expected: models.VoterTypeCC,
		},
		{
			name:     "DRep key hash",
			input:    lcommon.VoterTypeDRepKeyHash,
			expected: models.VoterTypeDRep,
		},
		{
			name:     "DRep script hash",
			input:    lcommon.VoterTypeDRepScriptHash,
			expected: models.VoterTypeDRep,
		},
		{
			name:     "SPO key hash",
			input:    lcommon.VoterTypeStakingPoolKeyHash,
			expected: models.VoterTypeSPO,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := mapVoterType(tt.input)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestExtractGovActionInfo_ParameterChange(t *testing.T) {
	parentId := &lcommon.GovActionId{
		TransactionId: [32]byte{1, 2, 3},
		GovActionIdx:  5,
	}
	action := &conway.ConwayParameterChangeGovAction{
		ActionId:   parentId,
		PolicyHash: []byte{0xAB, 0xCD},
	}

	actionType, parentTxHash, parentActionIdx, policyHash, err := extractGovActionInfo(
		action,
	)
	assert.NoError(t, err)

	assert.Equal(
		t,
		uint8(lcommon.GovActionTypeParameterChange),
		actionType,
	)
	assert.Equal(t, parentId.TransactionId[:], parentTxHash)
	assert.NotNil(t, parentActionIdx)
	assert.Equal(t, uint32(5), *parentActionIdx)
	assert.Equal(t, []byte{0xAB, 0xCD}, policyHash)
}

func TestExtractGovActionInfo_ParameterChangeNoParent(t *testing.T) {
	action := &conway.ConwayParameterChangeGovAction{}

	actionType, parentTxHash, parentActionIdx, policyHash, err := extractGovActionInfo(
		action,
	)
	assert.NoError(t, err)

	assert.Equal(
		t,
		uint8(lcommon.GovActionTypeParameterChange),
		actionType,
	)
	assert.Nil(t, parentTxHash)
	assert.Nil(t, parentActionIdx)
	assert.Nil(t, policyHash)
}

func TestExtractGovActionInfo_HardForkInitiation(t *testing.T) {
	parentId := &lcommon.GovActionId{
		TransactionId: [32]byte{4, 5, 6},
		GovActionIdx:  2,
	}
	action := &lcommon.HardForkInitiationGovAction{
		ActionId: parentId,
	}

	actionType, parentTxHash, parentActionIdx, _, err := extractGovActionInfo(
		action,
	)
	assert.NoError(t, err)

	assert.Equal(
		t,
		uint8(lcommon.GovActionTypeHardForkInitiation),
		actionType,
	)
	assert.Equal(t, parentId.TransactionId[:], parentTxHash)
	assert.NotNil(t, parentActionIdx)
	assert.Equal(t, uint32(2), *parentActionIdx)
}

func TestExtractGovActionInfo_TreasuryWithdrawal(t *testing.T) {
	action := &lcommon.TreasuryWithdrawalGovAction{
		PolicyHash: []byte{0x01, 0x02, 0x03},
	}

	actionType, parentTxHash, parentActionIdx, policyHash, err := extractGovActionInfo(
		action,
	)
	assert.NoError(t, err)

	assert.Equal(
		t,
		uint8(lcommon.GovActionTypeTreasuryWithdrawal),
		actionType,
	)
	assert.Nil(t, parentTxHash)
	assert.Nil(t, parentActionIdx)
	assert.Equal(t, []byte{0x01, 0x02, 0x03}, policyHash)
}

func TestExtractGovActionInfo_NoConfidence(t *testing.T) {
	parentId := &lcommon.GovActionId{
		TransactionId: [32]byte{7, 8, 9},
		GovActionIdx:  0,
	}
	action := &lcommon.NoConfidenceGovAction{
		ActionId: parentId,
	}

	actionType, parentTxHash, parentActionIdx, _, err := extractGovActionInfo(
		action,
	)
	assert.NoError(t, err)

	assert.Equal(
		t,
		uint8(lcommon.GovActionTypeNoConfidence),
		actionType,
	)
	assert.Equal(t, parentId.TransactionId[:], parentTxHash)
	assert.NotNil(t, parentActionIdx)
	assert.Equal(t, uint32(0), *parentActionIdx)
}

func TestExtractGovActionInfo_UpdateCommittee(t *testing.T) {
	action := &lcommon.UpdateCommitteeGovAction{}

	actionType, parentTxHash, parentActionIdx, policyHash, err := extractGovActionInfo(
		action,
	)
	assert.NoError(t, err)

	assert.Equal(
		t,
		uint8(lcommon.GovActionTypeUpdateCommittee),
		actionType,
	)
	assert.Nil(t, parentTxHash)
	assert.Nil(t, parentActionIdx)
	assert.Nil(t, policyHash)
}

func TestExtractGovActionInfo_NewConstitution(t *testing.T) {
	parentId := &lcommon.GovActionId{
		TransactionId: [32]byte{10, 11, 12},
		GovActionIdx:  3,
	}
	action := &lcommon.NewConstitutionGovAction{
		ActionId: parentId,
	}

	actionType, parentTxHash, parentActionIdx, _, err := extractGovActionInfo(
		action,
	)
	assert.NoError(t, err)

	assert.Equal(
		t,
		uint8(lcommon.GovActionTypeNewConstitution),
		actionType,
	)
	assert.Equal(t, parentId.TransactionId[:], parentTxHash)
	assert.NotNil(t, parentActionIdx)
	assert.Equal(t, uint32(3), *parentActionIdx)
}

func TestExtractGovActionInfo_Info(t *testing.T) {
	action := &lcommon.InfoGovAction{}

	actionType, parentTxHash, parentActionIdx, policyHash, err := extractGovActionInfo(
		action,
	)
	assert.NoError(t, err)

	assert.Equal(t, uint8(lcommon.GovActionTypeInfo), actionType)
	assert.Nil(t, parentTxHash)
	assert.Nil(t, parentActionIdx)
	assert.Nil(t, policyHash)
}

func TestProcessVotesRepairsMissingDRepRow(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "process_votes_missing_drep")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	db, err := database.New(&database.Config{
		DataDir:        tmpDir,
		Logger:         slog.New(slog.NewTextHandler(io.Discard, nil)),
		BlobPlugin:     "badger",
		MetadataPlugin: "sqlite",
	})
	require.NoError(t, err)
	defer db.Close()

	proposalTxHash := testHash32("proposal")
	returnAddress := append([]byte{0xE1}, testHash28("reward-account")...)
	require.NoError(
		t,
		db.SetGovernanceProposal(&models.GovernanceProposal{
			TxHash:        proposalTxHash,
			ActionIndex:   0,
			ActionType:    uint8(lcommon.GovActionTypeInfo),
			ProposedEpoch: 100,
			ExpiresEpoch:  120,
			AnchorURL:     "https://example.com/proposal",
			AnchorHash:    testHash32("proposal-anchor"),
			Deposit:       1,
			ReturnAddress: returnAddress,
			AddedSlot:     1000,
		}, nil),
	)
	proposal, err := db.GetGovernanceProposal(proposalTxHash, 0, nil)
	require.NoError(t, err)

	drepCred := testHash28("drep-voter")
	var voterHash [28]byte
	copy(voterHash[:], drepCred)
	var actionTxHash [32]byte
	copy(actionTxHash[:], proposalTxHash)
	voter := &lcommon.Voter{
		Type: lcommon.VoterTypeDRepKeyHash,
		Hash: voterHash,
	}
	actionID := &lcommon.GovActionId{
		TransactionId: actionTxHash,
		GovActionIdx:  0,
	}
	var txHash lcommon.Blake2b256
	copy(txHash[:], testHash32("vote-tx"))
	tx := &mockGovernanceTransaction{
		hash:    txHash,
		isValid: true,
		votingProcedures: lcommon.VotingProcedures{
			voter: {
				actionID: {Vote: models.VoteYes},
			},
		},
	}
	point := ocommon.Point{
		Slot: 2000,
		Hash: testHash32("vote-block"),
	}

	txn := db.Transaction(true)
	require.NoError(
		t,
		txn.Do(func(txn *database.Txn) error {
			return ProcessVotes(tx, point, 100, 20, db, txn)
		}),
	)
	txn.Release()

	drep, err := db.GetDrep(drepCred, true, nil)
	require.NoError(t, err)
	require.NotNil(t, drep)
	assert.True(t, drep.Active)
	assert.Equal(t, point.Slot, drep.AddedSlot)
	assert.Equal(t, uint64(100), drep.LastActivityEpoch)
	assert.Equal(t, uint64(120), drep.ExpiryEpoch)

	votes, err := db.GetGovernanceVotes(proposal.ID, nil)
	require.NoError(t, err)
	require.Len(t, votes, 1)
	assert.Equal(t, uint8(models.VoterTypeDRep), votes[0].VoterType)
	assert.Equal(t, drepCred, votes[0].VoterCredential)
	assert.Equal(t, uint8(models.VoteYes), votes[0].Vote)
	assert.Equal(t, point.Slot, votes[0].AddedSlot)
}

func TestProcessVotesRepairsMissingGovernanceProposal(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "process_votes_missing_proposal")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	db, err := database.New(&database.Config{
		DataDir:        tmpDir,
		Logger:         slog.New(slog.NewTextHandler(io.Discard, nil)),
		BlobPlugin:     "badger",
		MetadataPlugin: "sqlite",
	})
	require.NoError(t, err)
	defer db.Close()

	const (
		proposalEpoch     = uint64(100)
		proposalSlot      = uint64(1000)
		voteEpoch         = uint64(101)
		govActionLifetime = uint64(20)
	)
	require.NoError(
		t,
		db.SetEpoch(
			proposalSlot,
			proposalEpoch,
			nil,
			nil,
			nil,
			nil,
			conway.EraIdConway,
			1,
			432000,
			nil,
		),
	)
	pparams := testConwayProtocolParameters()
	pparams.GovActionValidityPeriod = govActionLifetime
	pparamsCbor, err := cbor.Encode(pparams)
	require.NoError(t, err)
	require.NoError(
		t,
		db.SetPParams(
			pparamsCbor,
			proposalSlot,
			proposalEpoch,
			conway.EraIdConway,
			nil,
		),
	)

	rewardAddress, err := lcommon.NewAddressFromBytes(
		append([]byte{0xE1}, testHash28("proposal-reward")...),
	)
	require.NoError(t, err)
	proposalProcedure := conway.ConwayProposalProcedure{
		PPDeposit:       42,
		PPRewardAccount: rewardAddress,
		PPGovAction: conway.ConwayGovAction{
			Type: uint(lcommon.GovActionTypeInfo),
			Action: &lcommon.InfoGovAction{
				Type: uint(lcommon.GovActionTypeInfo),
			},
		},
		PPAnchor: lcommon.GovAnchor{
			Url:      "https://example.com/proposal",
			DataHash: [32]byte(testHash32("proposal-anchor")),
		},
	}
	proposalBodyCbor, err := cbor.Encode(
		&conway.ConwayTransactionBody{
			TxProposalProcedures: []conway.ConwayProposalProcedure{
				proposalProcedure,
			},
		},
	)
	require.NoError(t, err)
	proposalTxHash := lcommon.Blake2b256Hash(proposalBodyCbor)
	proposalTx := &mockGovernanceTransaction{
		hash:               proposalTxHash,
		txType:             gledger.TxTypeConway,
		proposalProcedures: []lcommon.ProposalProcedure{proposalProcedure},
		isValid:            true,
	}
	proposalPoint := ocommon.Point{
		Slot: proposalSlot,
		Hash: testHash32("proposal-block"),
	}
	writeTxn := db.Transaction(true)
	require.NoError(
		t,
		writeTxn.Do(func(txn *database.Txn) error {
			if err := db.Blob().SetTx(
				txn.Blob(),
				proposalTxHash.Bytes(),
				proposalBodyCbor,
			); err != nil {
				return err
			}
			return db.Metadata().SetGapBlockTransaction(
				proposalTx,
				proposalPoint,
				0,
				txn.Metadata(),
			)
		}),
	)
	writeTxn.Release()

	_, err = db.GetGovernanceProposal(proposalTxHash.Bytes(), 0, nil)
	require.ErrorIs(t, err, models.ErrGovernanceProposalNotFound)

	ccCred := testHash28("committee-voter")
	var voterHash [28]byte
	copy(voterHash[:], ccCred)
	var actionTxHash [32]byte
	copy(actionTxHash[:], proposalTxHash.Bytes())
	voter := &lcommon.Voter{
		Type: lcommon.VoterTypeConstitutionalCommitteeHotKeyHash,
		Hash: voterHash,
	}
	actionID := &lcommon.GovActionId{
		TransactionId: actionTxHash,
		GovActionIdx:  0,
	}
	var voteTxHash lcommon.Blake2b256
	copy(voteTxHash[:], testHash32("vote-tx"))
	voteTx := &mockGovernanceTransaction{
		hash:    voteTxHash,
		txType:  gledger.TxTypeConway,
		isValid: true,
		votingProcedures: lcommon.VotingProcedures{
			voter: {
				actionID: {Vote: models.VoteYes},
			},
		},
	}
	votePoint := ocommon.Point{
		Slot: proposalSlot + 50,
		Hash: testHash32("vote-block"),
	}

	processTxn := db.Transaction(true)
	require.NoError(
		t,
		processTxn.Do(func(txn *database.Txn) error {
			return ProcessVotes(voteTx, votePoint, voteEpoch, 20, db, txn)
		}),
	)
	processTxn.Release()

	proposal, err := db.GetGovernanceProposal(
		proposalTxHash.Bytes(),
		0,
		nil,
	)
	require.NoError(t, err)
	require.NotNil(t, proposal)
	assert.Equal(t, proposalTxHash.Bytes(), proposal.TxHash)
	assert.Equal(t, uint32(0), proposal.ActionIndex)
	assert.Equal(
		t,
		uint8(lcommon.GovActionTypeInfo),
		proposal.ActionType,
	)
	assert.Equal(t, proposalEpoch, proposal.ProposedEpoch)
	assert.Equal(
		t,
		proposalEpoch+govActionLifetime,
		proposal.ExpiresEpoch,
	)
	assert.Equal(t, proposalProcedure.PPAnchor.Url, proposal.AnchorURL)
	assert.Equal(t, proposalProcedure.PPAnchor.DataHash[:], proposal.AnchorHash)
	assert.Equal(t, proposalProcedure.PPDeposit, proposal.Deposit)
	returnAddress, err := rewardAddress.Bytes()
	require.NoError(t, err)
	assert.Equal(t, returnAddress, proposal.ReturnAddress)
	assert.Equal(t, proposalSlot, proposal.AddedSlot)
	assert.NotEmpty(t, proposal.GovActionCbor)

	votes, err := db.GetGovernanceVotes(proposal.ID, nil)
	require.NoError(t, err)
	require.Len(t, votes, 1)
	assert.Equal(t, uint8(models.VoterTypeCC), votes[0].VoterType)
	assert.Equal(t, ccCred, votes[0].VoterCredential)
	assert.Equal(t, uint8(models.VoteYes), votes[0].Vote)
	assert.Equal(t, votePoint.Slot, votes[0].AddedSlot)
}
