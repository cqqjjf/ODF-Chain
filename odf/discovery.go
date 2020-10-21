// Copyright 2019 The go-odf Authors
// This file is part of the go-odf library.
//
// The go-odf library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-odf library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-odf library. If not, see <http://www.gnu.org/licenses/>.

package odf

import (
	"github.com/odf/go-odf/core"
	"github.com/odf/go-odf/core/forkid"
	"github.com/odf/go-odf/p2p"
	"github.com/odf/go-odf/p2p/dnsdisc"
	"github.com/odf/go-odf/p2p/enode"
	"github.com/odf/go-odf/rlp"
)

// odfEntry is the "odf" ENR entry which advertises odf protocol
// on the discovery network.
type odfEntry struct {
	ForkID forkid.ID // Fork identifier per EIP-2124

	// Ignore additional fields (for forward compatibility).
	Rest []rlp.RawValue `rlp:"tail"`
}

// ENRKey implements enr.Entry.
func (e odfEntry) ENRKey() string {
	return "odf"
}

// startEthEntryUpdate starts the ENR updater loop.
func (odf *Ethereum) startEthEntryUpdate(ln *enode.LocalNode) {
	var newHead = make(chan core.ChainHeadEvent, 10)
	sub := odf.blockchain.SubscribeChainHeadEvent(newHead)

	go func() {
		defer sub.Unsubscribe()
		for {
			select {
			case <-newHead:
				ln.Set(odf.currentEthEntry())
			case <-sub.Err():
				// Would be nice to sync with odf.Stop, but there is no
				// good way to do that.
				return
			}
		}
	}()
}

func (odf *Ethereum) currentEthEntry() *odfEntry {
	return &odfEntry{ForkID: forkid.NewID(odf.blockchain.Config(), odf.blockchain.Genesis().Hash(),
		odf.blockchain.CurrentHeader().Number.Uint64())}
}

// setupDiscovery creates the node discovery source for the odf protocol.
func (odf *Ethereum) setupDiscovery(cfg *p2p.Config) (enode.Iterator, error) {
	if cfg.NoDiscovery || len(odf.config.DiscoveryURLs) == 0 {
		return nil, nil
	}
	client := dnsdisc.NewClient(dnsdisc.Config{})
	return client.NewIterator(odf.config.DiscoveryURLs...)
}
