/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package node

import (
	"os"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

func TestUnjoinWithMissingChannelID(t *testing.T) {
	cmd := unjoinCmd()
	args := []string{}
	cmd.SetArgs(args)

	err := cmd.Execute()
	require.EqualError(t, err, "Must supply channel ID")
}

func TestUnjoinWithInvalidChannelID(t *testing.T) {
	testPath := "/tmp/hyperledger/test"
	os.RemoveAll(testPath)
	viper.Set("peer.fileSystemPath", testPath)
	defer os.RemoveAll(testPath)

	cmd := unjoinCmd()
	args := []string{"-c", "ch_xyz"}
	cmd.SetArgs(args)
	err := cmd.Execute()
	require.EqualError(t, err, "Unjoin channel [ch_xyz]: cannot update ledger status, ledger [ch_xyz] does not exist")
}
