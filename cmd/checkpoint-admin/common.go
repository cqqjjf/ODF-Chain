// Copyright 2019 The go-odf Authors
// This file is part of go-odf.
//
// go-odf is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// go-odf is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with go-odf. If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"strconv"

	"github.com/odf/go-odf/accounts"
	"github.com/odf/go-odf/accounts/abi/bind"
	"github.com/odf/go-odf/accounts/external"
	"github.com/odf/go-odf/cmd/utils"
	"github.com/odf/go-odf/common"
	"github.com/odf/go-odf/contracts/checkpointoracle"
	"github.com/odf/go-odf/odfclient"
	"github.com/odf/go-odf/params"
	"github.com/odf/go-odf/rpc"
	"gopkg.in/urfave/cli.v1"
)

// newClient creates a client with specified remote URL.
func newClient(ctx *cli.Context) *odfclient.Client {
	client, err := odfclient.Dial(ctx.GlobalString(nodeURLFlag.Name))
	if err != nil {
		utils.Fatalf("Failed to connect to Ethereum node: %v", err)
	}
	return client
}

// newRPCClient creates a rpc client with specified node URL.
func newRPCClient(url string) *rpc.Client {
	client, err := rpc.Dial(url)
	if err != nil {
		utils.Fatalf("Failed to connect to Ethereum node: %v", err)
	}
	return client
}

// getContractAddr retrieves the register contract address through
// rpc request.
func getContractAddr(client *rpc.Client) common.Address {
	var addr string
	if err := client.Call(&addr, "les_getCheckpointContractAddress"); err != nil {
		utils.Fatalf("Failed to fetch checkpoint oracle address: %v", err)
	}
	return common.HexToAddress(addr)
}

// getCheckpoint retrieves the specified checkpoint or the latest one
// through rpc request.
func getCheckpoint(ctx *cli.Context, client *rpc.Client) *params.TrustedCheckpoint {
	var checkpoint *params.TrustedCheckpoint

	if ctx.GlobalIsSet(indexFlag.Name) {
		var result [3]string
		index := uint64(ctx.GlobalInt64(indexFlag.Name))
		if err := client.Call(&result, "les_getCheckpoint", index); err != nil {
			utils.Fatalf("Failed to get local checkpoint %v, please ensure the les API is exposed", err)
		}
		checkpoint = &params.TrustedCheckpoint{
			SectionIndex: index,
			SectionHead:  common.HexToHash(result[0]),
			CHTRoot:      common.HexToHash(result[1]),
			BloomRoot:    common.HexToHash(result[2]),
		}
	} else {
		var result [4]string
		err := client.Call(&result, "les_latestCheckpoint")
		if err != nil {
			utils.Fatalf("Failed to get local checkpoint %v, please ensure the les API is exposed", err)
		}
		index, err := strconv.ParseUint(result[0], 0, 64)
		if err != nil {
			utils.Fatalf("Failed to parse checkpoint index %v", err)
		}
		checkpoint = &params.TrustedCheckpoint{
			SectionIndex: index,
			SectionHead:  common.HexToHash(result[1]),
			CHTRoot:      common.HexToHash(result[2]),
			BloomRoot:    common.HexToHash(result[3]),
		}
	}
	return checkpoint
}

// newContract creates a registrar contract instance with specified
// contract address or the default contracts for mainnet or testnet.
func newContract(client *rpc.Client) (common.Address, *checkpointoracle.CheckpointOracle) {
	addr := getContractAddr(client)
	if addr == (common.Address{}) {
		utils.Fatalf("No specified registrar contract address")
	}
	contract, err := checkpointoracle.NewCheckpointOracle(addr, odfclient.NewClient(client))
	if err != nil {
		utils.Fatalf("Failed to setup registrar contract %s: %v", addr, err)
	}
	return addr, contract
}

// newClefSigner sets up a clef backend and returns a clef transaction signer.
func newClefSigner(ctx *cli.Context) *bind.TransactOpts {
	clef, err := external.NewExternalSigner(ctx.String(clefURLFlag.Name))
	if err != nil {
		utils.Fatalf("Failed to create clef signer %v", err)
	}
	return bind.NewClefTransactor(clef, accounts.Account{Address: common.HexToAddress(ctx.String(signerFlag.Name))})
}
