// Copyright 2016 The go-odf Authors
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

// Package les implements the Light Ethereum Subprotocol.
package les

import (
	"fmt"
	"time"

	"github.com/odf/go-odf/accounts"
	"github.com/odf/go-odf/common"
	"github.com/odf/go-odf/common/hexutil"
	"github.com/odf/go-odf/common/mclock"
	"github.com/odf/go-odf/consensus"
	"github.com/odf/go-odf/core"
	"github.com/odf/go-odf/core/bloombits"
	"github.com/odf/go-odf/core/rawdb"
	"github.com/odf/go-odf/core/types"
	"github.com/odf/go-odf/odf"
	"github.com/odf/go-odf/odf/downloader"
	"github.com/odf/go-odf/odf/filters"
	"github.com/odf/go-odf/odf/gasprice"
	"github.com/odf/go-odf/event"
	"github.com/odf/go-odf/internal/odfapi"
	lpc "github.com/odf/go-odf/les/lespay/client"
	"github.com/odf/go-odf/light"
	"github.com/odf/go-odf/log"
	"github.com/odf/go-odf/node"
	"github.com/odf/go-odf/p2p"
	"github.com/odf/go-odf/p2p/enode"
	"github.com/odf/go-odf/params"
	"github.com/odf/go-odf/rpc"
)

type LightEthereum struct {
	lesCommons

	peers          *serverPeerSet
	reqDist        *requestDistributor
	retriever      *retrieveManager
	odr            *LesOdr
	relay          *lesTxRelay
	handler        *clientHandler
	txPool         *light.TxPool
	blockchain     *light.LightChain
	serverPool     *serverPool
	valueTracker   *lpc.ValueTracker
	dialCandidates enode.Iterator
	pruner         *pruner

	bloomRequests chan chan *bloombits.Retrieval // Channel receiving bloom data retrieval requests
	bloomIndexer  *core.ChainIndexer             // Bloom indexer operating during block imports

	ApiBackend     *LesApiBackend
	eventMux       *event.TypeMux
	engine         consensus.Engine
	accountManager *accounts.Manager
	netRPCService  *odfapi.PublicNetAPI

	p2pServer *p2p.Server
}

// New creates an instance of the light client.
func New(stack *node.Node, config *odf.Config) (*LightEthereum, error) {
	chainDb, err := stack.OpenDatabase("lightchaindata", config.DatabaseCache, config.DatabaseHandles, "odf/db/chaindata/")
	if err != nil {
		return nil, err
	}
	lespayDb, err := stack.OpenDatabase("lespay", 0, 0, "odf/db/lespay")
	if err != nil {
		return nil, err
	}
	chainConfig, genesisHash, genesisErr := core.SetupGenesisBlock(chainDb, config.Genesis)
	if _, isCompat := genesisErr.(*params.ConfigCompatError); genesisErr != nil && !isCompat {
		return nil, genesisErr
	}
	log.Info("Initialised chain configuration", "config", chainConfig)

	peers := newServerPeerSet()
	lodf := &LightEthereum{
		lesCommons: lesCommons{
			genesis:     genesisHash,
			config:      config,
			chainConfig: chainConfig,
			iConfig:     light.DefaultClientIndexerConfig,
			chainDb:     chainDb,
			closeCh:     make(chan struct{}),
		},
		peers:          peers,
		eventMux:       stack.EventMux(),
		reqDist:        newRequestDistributor(peers, &mclock.System{}),
		accountManager: stack.AccountManager(),
		engine:         odf.CreateConsensusEngine(stack, chainConfig, &config.Ethash, nil, false, chainDb),
		bloomRequests:  make(chan chan *bloombits.Retrieval),
		bloomIndexer:   odf.NewBloomIndexer(chainDb, params.BloomBitsBlocksClient, params.HelperTrieConfirmations),
		valueTracker:   lpc.NewValueTracker(lespayDb, &mclock.System{}, requestList, time.Minute, 1/float64(time.Hour), 1/float64(time.Hour*100), 1/float64(time.Hour*1000)),
		p2pServer:      stack.Server(),
	}
	peers.subscribe((*vtSubscription)(lodf.valueTracker))

	dnsdisc, err := lodf.setupDiscovery(&stack.Config().P2P)
	if err != nil {
		return nil, err
	}
	lodf.serverPool = newServerPool(lespayDb, []byte("serverpool:"), lodf.valueTracker, dnsdisc, time.Second, nil, &mclock.System{}, config.UltraLightServers)
	peers.subscribe(lodf.serverPool)
	lodf.dialCandidates = lodf.serverPool.dialIterator

	lodf.retriever = newRetrieveManager(peers, lodf.reqDist, lodf.serverPool.getTimeout)
	lodf.relay = newLesTxRelay(peers, lodf.retriever)

	lodf.odr = NewLesOdr(chainDb, light.DefaultClientIndexerConfig, lodf.retriever)
	lodf.chtIndexer = light.NewChtIndexer(chainDb, lodf.odr, params.CHTFrequency, params.HelperTrieConfirmations, config.LightNoPrune)
	lodf.bloomTrieIndexer = light.NewBloomTrieIndexer(chainDb, lodf.odr, params.BloomBitsBlocksClient, params.BloomTrieFrequency, config.LightNoPrune)
	lodf.odr.SetIndexers(lodf.chtIndexer, lodf.bloomTrieIndexer, lodf.bloomIndexer)

	checkpoint := config.Checkpoint
	if checkpoint == nil {
		checkpoint = params.TrustedCheckpoints[genesisHash]
	}
	// Note: NewLightChain adds the trusted checkpoint so it needs an ODR with
	// indexers already set but not started yet
	if lodf.blockchain, err = light.NewLightChain(lodf.odr, lodf.chainConfig, lodf.engine, checkpoint); err != nil {
		return nil, err
	}
	lodf.chainReader = lodf.blockchain
	lodf.txPool = light.NewTxPool(lodf.chainConfig, lodf.blockchain, lodf.relay)

	// Set up checkpoint oracle.
	lodf.oracle = lodf.setupOracle(stack, genesisHash, config)

	// Note: AddChildIndexer starts the update process for the child
	lodf.bloomIndexer.AddChildIndexer(lodf.bloomTrieIndexer)
	lodf.chtIndexer.Start(lodf.blockchain)
	lodf.bloomIndexer.Start(lodf.blockchain)

	// Start a light chain pruner to delete useless historical data.
	lodf.pruner = newPruner(chainDb, lodf.chtIndexer, lodf.bloomTrieIndexer)

	// Rewind the chain in case of an incompatible config upgrade.
	if compat, ok := genesisErr.(*params.ConfigCompatError); ok {
		log.Warn("Rewinding chain to upgrade configuration", "err", compat)
		lodf.blockchain.SetHead(compat.RewindTo)
		rawdb.WriteChainConfig(chainDb, genesisHash, chainConfig)
	}

	lodf.ApiBackend = &LesApiBackend{stack.Config().ExtRPCEnabled(), lodf, nil}
	gpoParams := config.GPO
	if gpoParams.Default == nil {
		gpoParams.Default = config.Miner.GasPrice
	}
	lodf.ApiBackend.gpo = gasprice.NewOracle(lodf.ApiBackend, gpoParams)

	lodf.handler = newClientHandler(config.UltraLightServers, config.UltraLightFraction, checkpoint, lodf)
	if lodf.handler.ulc != nil {
		log.Warn("Ultra light client is enabled", "trustedNodes", len(lodf.handler.ulc.keys), "minTrustedFraction", lodf.handler.ulc.fraction)
		lodf.blockchain.DisableCheckFreq()
	}

	lodf.netRPCService = odfapi.NewPublicNetAPI(lodf.p2pServer, lodf.config.NetworkId)

	// Register the backend on the node
	stack.RegisterAPIs(lodf.APIs())
	stack.RegisterProtocols(lodf.Protocols())
	stack.RegisterLifecycle(lodf)

	return lodf, nil
}

// vtSubscription implements serverPeerSubscriber
type vtSubscription lpc.ValueTracker

// registerPeer implements serverPeerSubscriber
func (v *vtSubscription) registerPeer(p *serverPeer) {
	vt := (*lpc.ValueTracker)(v)
	p.setValueTracker(vt, vt.Register(p.ID()))
	p.updateVtParams()
}

// unregisterPeer implements serverPeerSubscriber
func (v *vtSubscription) unregisterPeer(p *serverPeer) {
	vt := (*lpc.ValueTracker)(v)
	vt.Unregister(p.ID())
	p.setValueTracker(nil, nil)
}

type LightDummyAPI struct{}

// Etherbase is the address that mining rewards will be send to
func (s *LightDummyAPI) Etherbase() (common.Address, error) {
	return common.Address{}, fmt.Errorf("mining is not supported in light mode")
}

// Coinbase is the address that mining rewards will be send to (alias for Etherbase)
func (s *LightDummyAPI) Coinbase() (common.Address, error) {
	return common.Address{}, fmt.Errorf("mining is not supported in light mode")
}

// Hashrate returns the POW hashrate
func (s *LightDummyAPI) Hashrate() hexutil.Uint {
	return 0
}

// Mining returns an indication if this node is currently mining.
func (s *LightDummyAPI) Mining() bool {
	return false
}

// APIs returns the collection of RPC services the odf package offers.
// NOTE, some of these services probably need to be moved to somewhere else.
func (s *LightEthereum) APIs() []rpc.API {
	apis := odfapi.GetAPIs(s.ApiBackend)
	apis = append(apis, s.engine.APIs(s.BlockChain().HeaderChain())...)
	return append(apis, []rpc.API{
		{
			Namespace: "odf",
			Version:   "1.0",
			Service:   &LightDummyAPI{},
			Public:    true,
		}, {
			Namespace: "odf",
			Version:   "1.0",
			Service:   downloader.NewPublicDownloaderAPI(s.handler.downloader, s.eventMux),
			Public:    true,
		}, {
			Namespace: "odf",
			Version:   "1.0",
			Service:   filters.NewPublicFilterAPI(s.ApiBackend, true),
			Public:    true,
		}, {
			Namespace: "net",
			Version:   "1.0",
			Service:   s.netRPCService,
			Public:    true,
		}, {
			Namespace: "les",
			Version:   "1.0",
			Service:   NewPrivateLightAPI(&s.lesCommons),
			Public:    false,
		}, {
			Namespace: "lespay",
			Version:   "1.0",
			Service:   lpc.NewPrivateClientAPI(s.valueTracker),
			Public:    false,
		},
	}...)
}

func (s *LightEthereum) ResetWithGenesisBlock(gb *types.Block) {
	s.blockchain.ResetWithGenesisBlock(gb)
}

func (s *LightEthereum) BlockChain() *light.LightChain      { return s.blockchain }
func (s *LightEthereum) TxPool() *light.TxPool              { return s.txPool }
func (s *LightEthereum) Engine() consensus.Engine           { return s.engine }
func (s *LightEthereum) LesVersion() int                    { return int(ClientProtocolVersions[0]) }
func (s *LightEthereum) Downloader() *downloader.Downloader { return s.handler.downloader }
func (s *LightEthereum) EventMux() *event.TypeMux           { return s.eventMux }

// Protocols returns all the currently configured network protocols to start.
func (s *LightEthereum) Protocols() []p2p.Protocol {
	return s.makeProtocols(ClientProtocolVersions, s.handler.runPeer, func(id enode.ID) interface{} {
		if p := s.peers.peer(id.String()); p != nil {
			return p.Info()
		}
		return nil
	}, s.dialCandidates)
}

// Start implements node.Lifecycle, starting all internal goroutines needed by the
// light odf protocol implementation.
func (s *LightEthereum) Start() error {
	log.Warn("Light client mode is an experimental feature")

	s.serverPool.start()
	// Start bloom request workers.
	s.wg.Add(bloomServiceThreads)
	s.startBloomHandlers(params.BloomBitsBlocksClient)
	s.handler.start()

	return nil
}

// Stop implements node.Lifecycle, terminating all internal goroutines used by the
// Ethereum protocol.
func (s *LightEthereum) Stop() error {
	close(s.closeCh)
	s.serverPool.stop()
	s.valueTracker.Stop()
	s.peers.close()
	s.reqDist.close()
	s.odr.Stop()
	s.relay.Stop()
	s.bloomIndexer.Close()
	s.chtIndexer.Close()
	s.blockchain.Stop()
	s.handler.stop()
	s.txPool.Stop()
	s.engine.Close()
	s.pruner.close()
	s.eventMux.Stop()
	s.chainDb.Close()
	s.wg.Wait()
	log.Info("Light odf stopped")
	return nil
}
