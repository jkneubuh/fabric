/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package kvledger

import (
	"testing"

	configtxtest "github.com/hyperledger/fabric/common/configtx/test"
	"github.com/hyperledger/fabric/core/ledger/kvledger/msgs"
	"github.com/hyperledger/fabric/core/ledger/mock"
	"github.com/stretchr/testify/require"
)

func TestUnjoinChannel(t *testing.T) {
	conf, cleanup := testConfig(t)
	conf.HistoryDBConfig.Enabled = false
	defer cleanup()

	ledgerID := "ledger_unjoin"

	provider := testutilNewProvider(conf, t, &mock.DeployedChaincodeInfoProvider{})
	activeLedgerIDs, err := provider.List()
	require.NoError(t, err)
	require.Len(t, activeLedgerIDs, 0)

	genesisBlock, err := configtxtest.MakeGenesisBlock(ledgerID)
	require.NoError(t, err)
	_, err = provider.CreateFromGenesisBlock(genesisBlock)
	require.NoError(t, err)

	activeLedgerIDs, err = provider.List()
	require.NoError(t, err)
	require.Len(t, activeLedgerIDs, 1)
	provider.Close()

	// Unjoin the channel from the peer
	err = UnjoinChannel(conf, ledgerID)
	require.NoError(t, err)

	// channel should no longer be present in the channel list
	provider = testutilNewProvider(conf, t, &mock.DeployedChaincodeInfoProvider{})
	activeLedgerIDs, err = provider.List()
	require.NoError(t, err)
	require.Len(t, activeLedgerIDs, 0)

	provider.Close()
}

// Unjoining an unjoined channel is OK.
func TestUnjoinUnjoinedChannel(t *testing.T) {
	conf, cleanup := testConfig(t)
	conf.HistoryDBConfig.Enabled = false
	defer cleanup()

	ledgerID := "ledger_unjoin_unjoined"

	provider := testutilNewProvider(conf, t, &mock.DeployedChaincodeInfoProvider{})
	genesisBlock, err := configtxtest.MakeGenesisBlock(ledgerID)
	require.NoError(t, err)
	_, err = provider.CreateFromGenesisBlock(genesisBlock)
	require.NoError(t, err)
	provider.Close()

	// unjoin the channel
	require.NoError(t, UnjoinChannel(conf, ledgerID))

	provider = testutilNewProvider(conf, t, &mock.DeployedChaincodeInfoProvider{})
	activeLedgerIDs, err := provider.List()
	require.NoError(t, err)
	require.Len(t, activeLedgerIDs, 0)
	provider.Close()

	// Subsequent unjoins will not throw an error
	require.NoError(t, UnjoinChannel(conf, ledgerID))
	require.NoError(t, UnjoinChannel(conf, ledgerID))
	require.NoError(t, UnjoinChannel(conf, ledgerID))
}

func TestUpdateLedgerStatus(t *testing.T) {
	conf, cleanup := testConfig(t)
	defer cleanup()

	provider := testutilNewProvider(conf, t, &mock.DeployedChaincodeInfoProvider{})
	ledgerID := constructTestLedgerID(11)
	genesisBlock, err := configtxtest.MakeGenesisBlock(ledgerID)
	require.NoError(t, err)
	_, err = provider.CreateFromGenesisBlock(genesisBlock)
	require.NoError(t, err)

	meta, err := provider.idStore.getLedgerMetadata(ledgerID)
	require.NotNil(t, meta)
	require.NoError(t, err)
	require.Equal(t, msgs.Status_ACTIVE, meta.GetStatus())
	provider.Close()

	// Update the ledger status
	err = updateLedgerStatus(conf, ledgerID, msgs.Status_UNDER_DELETION)
	require.NoError(t, err)

	provider = testutilNewProvider(conf, t, &mock.DeployedChaincodeInfoProvider{})
	meta, err = provider.idStore.getLedgerMetadata(ledgerID)
	require.NotNil(t, meta)
	require.NoError(t, err)
	require.Equal(t, msgs.Status_UNDER_DELETION, meta.GetStatus())
	provider.Close()
}

func TestUnjoinWithRunningPeerErrors(t *testing.T) {
	conf, cleanup := testConfig(t)
	conf.HistoryDBConfig.Enabled = false
	defer cleanup()

	provider := testutilNewProvider(conf, t, &mock.DeployedChaincodeInfoProvider{})
	defer provider.Close()

	ledgerID := constructTestLedgerID(1)
	genesisBlock, _ := configtxtest.MakeGenesisBlock(ledgerID)
	_, err := provider.CreateFromGenesisBlock(genesisBlock)
	require.NoError(t, err)

	// Fail when provider is open (e.g. peer is running)
	require.Error(t, UnjoinChannel(conf, ledgerID),
		"as another peer node command is executing, wait for that command to complete its execution or terminate it before retrying")
}

func TestUnjoinWithMissingChannelErrors(t *testing.T) {
	conf, cleanup := testConfig(t)
	conf.HistoryDBConfig.Enabled = false
	defer cleanup()

	// fail if channel does not exist
	require.EqualError(t, UnjoinChannel(conf, "__invalid_channel"),
		"Unjoin channel [__invalid_channel]: cannot update ledger status, ledger [__invalid_channel] does not exist")
}

func TestUnjoinChannelWithInvalidMetadataErrors(t *testing.T) {
	conf, cleanup := testConfig(t)
	conf.HistoryDBConfig.Enabled = false
	defer cleanup()

	provider := testutilNewProvider(conf, t, &mock.DeployedChaincodeInfoProvider{})

	ledgerID := constructTestLedgerID(99)
	genesisBlock, _ := configtxtest.MakeGenesisBlock(ledgerID)
	_, err := provider.CreateFromGenesisBlock(genesisBlock)
	require.NoError(t, err)

	// purposely set an invalid metatdata
	require.NoError(t, provider.idStore.db.Put(metadataKey(ledgerID), []byte("invalid"), true))
	provider.Close()

	// fail if metadata can not be unmarshaled
	require.EqualError(t, UnjoinChannel(conf, ledgerID),
		"Unjoin channel [ledger_000099]: error unmarshalling ledger metadata: unexpected EOF")
}
