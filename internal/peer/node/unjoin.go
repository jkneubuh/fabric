/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package node

import (
	"github.com/hyperledger/fabric/core/ledger/kvledger"
	"github.com/hyperledger/fabric/internal/peer/common"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func unjoinCmd() *cobra.Command {
	var channelID string

	cmd := &cobra.Command{
		Use:   "unjoin",
		Short: "Unjoin the peer from a channel.",
		Long:  "Unjoin the peer from a channel.  When the command is executed, the peer must be offline.",
		RunE: func(cmd *cobra.Command, args []string) error {
			if channelID == common.UndefinedParamValue {
				return errors.New("Must supply channel ID")
			}

			config := ledgerConfig()
			return kvledger.UnjoinChannel(config, channelID)
		},
	}
	flags := cmd.Flags()
	flags.StringVarP(&channelID, "channelID", "c", common.UndefinedParamValue, "Channel to unjoin.")

	return cmd
}
