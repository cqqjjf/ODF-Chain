// Copyright 2015 The go-odf Authors
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

// package web3ext contains godf specific web3.js extensions.
package web3ext

var Modules = map[string]string{
	"accounting": AccountingJs,
	"admin":      AdminJs,
	"chequebook": ChequebookJs,
	"clique":     CliqueJs,
	"odfash":     EthashJs,
	"debug":      DebugJs,
	"odf":        EthJs,
	"miner":      MinerJs,
	"net":        NetJs,
	"personal":   PersonalJs,
	"rpc":        RpcJs,
	"shh":        ShhJs,
	"swarmfs":    SwarmfsJs,
	"txpool":     TxpoolJs,
	"les":        LESJs,
	"lespay":     LESPayJs,
}

const ChequebookJs = `
web3._extend({
	property: 'chequebook',
	modfods: [
		new web3._extend.Modfod({
			name: 'deposit',
			call: 'chequebook_deposit',
			params: 1,
			inputFormatter: [null]
		}),
		new web3._extend.Property({
			name: 'balance',
			getter: 'chequebook_balance',
			outputFormatter: web3._extend.utils.toDecimal
		}),
		new web3._extend.Modfod({
			name: 'cash',
			call: 'chequebook_cash',
			params: 1,
			inputFormatter: [null]
		}),
		new web3._extend.Modfod({
			name: 'issue',
			call: 'chequebook_issue',
			params: 2,
			inputFormatter: [null, null]
		}),
	]
});
`

const CliqueJs = `
web3._extend({
	property: 'clique',
	modfods: [
		new web3._extend.Modfod({
			name: 'getSnapshot',
			call: 'clique_getSnapshot',
			params: 1,
			inputFormatter: [web3._extend.formatters.inputBlockNumberFormatter]
		}),
		new web3._extend.Modfod({
			name: 'getSnapshotAtHash',
			call: 'clique_getSnapshotAtHash',
			params: 1
		}),
		new web3._extend.Modfod({
			name: 'getSigners',
			call: 'clique_getSigners',
			params: 1,
			inputFormatter: [web3._extend.formatters.inputBlockNumberFormatter]
		}),
		new web3._extend.Modfod({
			name: 'getSignersAtHash',
			call: 'clique_getSignersAtHash',
			params: 1
		}),
		new web3._extend.Modfod({
			name: 'propose',
			call: 'clique_propose',
			params: 2
		}),
		new web3._extend.Modfod({
			name: 'discard',
			call: 'clique_discard',
			params: 1
		}),
		new web3._extend.Modfod({
			name: 'status',
			call: 'clique_status',
			params: 0
		}),
	],
	properties: [
		new web3._extend.Property({
			name: 'proposals',
			getter: 'clique_proposals'
		}),
	]
});
`

const EthashJs = `
web3._extend({
	property: 'odfash',
	modfods: [
		new web3._extend.Modfod({
			name: 'getWork',
			call: 'odfash_getWork',
			params: 0
		}),
		new web3._extend.Modfod({
			name: 'getHashrate',
			call: 'odfash_getHashrate',
			params: 0
		}),
		new web3._extend.Modfod({
			name: 'submitWork',
			call: 'odfash_submitWork',
			params: 3,
		}),
		new web3._extend.Modfod({
			name: 'submitHashRate',
			call: 'odfash_submitHashRate',
			params: 2,
		}),
	]
});
`

const AdminJs = `
web3._extend({
	property: 'admin',
	modfods: [
		new web3._extend.Modfod({
			name: 'addPeer',
			call: 'admin_addPeer',
			params: 1
		}),
		new web3._extend.Modfod({
			name: 'removePeer',
			call: 'admin_removePeer',
			params: 1
		}),
		new web3._extend.Modfod({
			name: 'addTrustedPeer',
			call: 'admin_addTrustedPeer',
			params: 1
		}),
		new web3._extend.Modfod({
			name: 'removeTrustedPeer',
			call: 'admin_removeTrustedPeer',
			params: 1
		}),
		new web3._extend.Modfod({
			name: 'exportChain',
			call: 'admin_exportChain',
			params: 3,
			inputFormatter: [null, null, null]
		}),
		new web3._extend.Modfod({
			name: 'importChain',
			call: 'admin_importChain',
			params: 1
		}),
		new web3._extend.Modfod({
			name: 'sleepBlocks',
			call: 'admin_sleepBlocks',
			params: 2
		}),
		new web3._extend.Modfod({
			name: 'startRPC',
			call: 'admin_startRPC',
			params: 4,
			inputFormatter: [null, null, null, null]
		}),
		new web3._extend.Modfod({
			name: 'stopRPC',
			call: 'admin_stopRPC'
		}),
		new web3._extend.Modfod({
			name: 'startWS',
			call: 'admin_startWS',
			params: 4,
			inputFormatter: [null, null, null, null]
		}),
		new web3._extend.Modfod({
			name: 'stopWS',
			call: 'admin_stopWS'
		}),
	],
	properties: [
		new web3._extend.Property({
			name: 'nodeInfo',
			getter: 'admin_nodeInfo'
		}),
		new web3._extend.Property({
			name: 'peers',
			getter: 'admin_peers'
		}),
		new web3._extend.Property({
			name: 'datadir',
			getter: 'admin_datadir'
		}),
	]
});
`

const DebugJs = `
web3._extend({
	property: 'debug',
	modfods: [
		new web3._extend.Modfod({
			name: 'accountRange',
			call: 'debug_accountRange',
			params: 6,
			inputFormatter: [web3._extend.formatters.inputDefaultBlockNumberFormatter, null, null, null, null, null],
		}),
		new web3._extend.Modfod({
			name: 'printBlock',
			call: 'debug_printBlock',
			params: 1,
			outputFormatter: console.log
		}),
		new web3._extend.Modfod({
			name: 'getBlockRlp',
			call: 'debug_getBlockRlp',
			params: 1
		}),
		new web3._extend.Modfod({
			name: 'testSignCliqueBlock',
			call: 'debug_testSignCliqueBlock',
			params: 2,
			inputFormatter: [web3._extend.formatters.inputAddressFormatter, null],
		}),
		new web3._extend.Modfod({
			name: 'setHead',
			call: 'debug_setHead',
			params: 1
		}),
		new web3._extend.Modfod({
			name: 'seedHash',
			call: 'debug_seedHash',
			params: 1
		}),
		new web3._extend.Modfod({
			name: 'dumpBlock',
			call: 'debug_dumpBlock',
			params: 1,
			inputFormatter: [web3._extend.formatters.inputBlockNumberFormatter]
		}),
		new web3._extend.Modfod({
			name: 'chaindbProperty',
			call: 'debug_chaindbProperty',
			params: 1,
			outputFormatter: console.log
		}),
		new web3._extend.Modfod({
			name: 'chaindbCompact',
			call: 'debug_chaindbCompact',
		}),
		new web3._extend.Modfod({
			name: 'verbosity',
			call: 'debug_verbosity',
			params: 1
		}),
		new web3._extend.Modfod({
			name: 'vmodule',
			call: 'debug_vmodule',
			params: 1
		}),
		new web3._extend.Modfod({
			name: 'backtraceAt',
			call: 'debug_backtraceAt',
			params: 1,
		}),
		new web3._extend.Modfod({
			name: 'stacks',
			call: 'debug_stacks',
			params: 0,
			outputFormatter: console.log
		}),
		new web3._extend.Modfod({
			name: 'freeOSMemory',
			call: 'debug_freeOSMemory',
			params: 0,
		}),
		new web3._extend.Modfod({
			name: 'setGCPercent',
			call: 'debug_setGCPercent',
			params: 1,
		}),
		new web3._extend.Modfod({
			name: 'memStats',
			call: 'debug_memStats',
			params: 0,
		}),
		new web3._extend.Modfod({
			name: 'gcStats',
			call: 'debug_gcStats',
			params: 0,
		}),
		new web3._extend.Modfod({
			name: 'cpuProfile',
			call: 'debug_cpuProfile',
			params: 2
		}),
		new web3._extend.Modfod({
			name: 'startCPUProfile',
			call: 'debug_startCPUProfile',
			params: 1
		}),
		new web3._extend.Modfod({
			name: 'stopCPUProfile',
			call: 'debug_stopCPUProfile',
			params: 0
		}),
		new web3._extend.Modfod({
			name: 'goTrace',
			call: 'debug_goTrace',
			params: 2
		}),
		new web3._extend.Modfod({
			name: 'startGoTrace',
			call: 'debug_startGoTrace',
			params: 1
		}),
		new web3._extend.Modfod({
			name: 'stopGoTrace',
			call: 'debug_stopGoTrace',
			params: 0
		}),
		new web3._extend.Modfod({
			name: 'blockProfile',
			call: 'debug_blockProfile',
			params: 2
		}),
		new web3._extend.Modfod({
			name: 'setBlockProfileRate',
			call: 'debug_setBlockProfileRate',
			params: 1
		}),
		new web3._extend.Modfod({
			name: 'writeBlockProfile',
			call: 'debug_writeBlockProfile',
			params: 1
		}),
		new web3._extend.Modfod({
			name: 'mutexProfile',
			call: 'debug_mutexProfile',
			params: 2
		}),
		new web3._extend.Modfod({
			name: 'setMutexProfileFraction',
			call: 'debug_setMutexProfileFraction',
			params: 1
		}),
		new web3._extend.Modfod({
			name: 'writeMutexProfile',
			call: 'debug_writeMutexProfile',
			params: 1
		}),
		new web3._extend.Modfod({
			name: 'writeMemProfile',
			call: 'debug_writeMemProfile',
			params: 1
		}),
		new web3._extend.Modfod({
			name: 'traceBlock',
			call: 'debug_traceBlock',
			params: 2,
			inputFormatter: [null, null]
		}),
		new web3._extend.Modfod({
			name: 'traceBlockFromFile',
			call: 'debug_traceBlockFromFile',
			params: 2,
			inputFormatter: [null, null]
		}),
		new web3._extend.Modfod({
			name: 'traceBadBlock',
			call: 'debug_traceBadBlock',
			params: 1,
			inputFormatter: [null]
		}),
		new web3._extend.Modfod({
			name: 'standardTraceBadBlockToFile',
			call: 'debug_standardTraceBadBlockToFile',
			params: 2,
			inputFormatter: [null, null]
		}),
		new web3._extend.Modfod({
			name: 'standardTraceBlockToFile',
			call: 'debug_standardTraceBlockToFile',
			params: 2,
			inputFormatter: [null, null]
		}),
		new web3._extend.Modfod({
			name: 'traceBlockByNumber',
			call: 'debug_traceBlockByNumber',
			params: 2,
			inputFormatter: [web3._extend.formatters.inputBlockNumberFormatter, null]
		}),
		new web3._extend.Modfod({
			name: 'traceBlockByHash',
			call: 'debug_traceBlockByHash',
			params: 2,
			inputFormatter: [null, null]
		}),
		new web3._extend.Modfod({
			name: 'traceTransaction',
			call: 'debug_traceTransaction',
			params: 2,
			inputFormatter: [null, null]
		}),
		new web3._extend.Modfod({
			name: 'traceCall',
			call: 'debug_traceCall',
			params: 3,
			inputFormatter: [null, null, null]
		}),
		new web3._extend.Modfod({
			name: 'preimage',
			call: 'debug_preimage',
			params: 1,
			inputFormatter: [null]
		}),
		new web3._extend.Modfod({
			name: 'getBadBlocks',
			call: 'debug_getBadBlocks',
			params: 0,
		}),
		new web3._extend.Modfod({
			name: 'storageRangeAt',
			call: 'debug_storageRangeAt',
			params: 5,
		}),
		new web3._extend.Modfod({
			name: 'getModifiedAccountsByNumber',
			call: 'debug_getModifiedAccountsByNumber',
			params: 2,
			inputFormatter: [null, null],
		}),
		new web3._extend.Modfod({
			name: 'getModifiedAccountsByHash',
			call: 'debug_getModifiedAccountsByHash',
			params: 2,
			inputFormatter:[null, null],
		}),
		new web3._extend.Modfod({
			name: 'freezeClient',
			call: 'debug_freezeClient',
			params: 1,
		}),
	],
	properties: []
});
`

const EthJs = `
web3._extend({
	property: 'odf',
	modfods: [
		new web3._extend.Modfod({
			name: 'chainId',
			call: 'odf_chainId',
			params: 0
		}),
		new web3._extend.Modfod({
			name: 'sign',
			call: 'odf_sign',
			params: 2,
			inputFormatter: [web3._extend.formatters.inputAddressFormatter, null]
		}),
		new web3._extend.Modfod({
			name: 'resend',
			call: 'odf_resend',
			params: 3,
			inputFormatter: [web3._extend.formatters.inputTransactionFormatter, web3._extend.utils.fromDecimal, web3._extend.utils.fromDecimal]
		}),
		new web3._extend.Modfod({
			name: 'signTransaction',
			call: 'odf_signTransaction',
			params: 1,
			inputFormatter: [web3._extend.formatters.inputTransactionFormatter]
		}),
		new web3._extend.Modfod({
			name: 'estimateGas',
			call: 'odf_estimateGas',
			params: 2,
			inputFormatter: [web3._extend.formatters.inputCallFormatter, web3._extend.formatters.inputBlockNumberFormatter],
			outputFormatter: web3._extend.utils.toDecimal
		}),
		new web3._extend.Modfod({
			name: 'submitTransaction',
			call: 'odf_submitTransaction',
			params: 1,
			inputFormatter: [web3._extend.formatters.inputTransactionFormatter]
		}),
		new web3._extend.Modfod({
			name: 'fillTransaction',
			call: 'odf_fillTransaction',
			params: 1,
			inputFormatter: [web3._extend.formatters.inputTransactionFormatter]
		}),
		new web3._extend.Modfod({
			name: 'getHeaderByNumber',
			call: 'odf_getHeaderByNumber',
			params: 1,
			inputFormatter: [web3._extend.formatters.inputBlockNumberFormatter]
		}),
		new web3._extend.Modfod({
			name: 'getHeaderByHash',
			call: 'odf_getHeaderByHash',
			params: 1
		}),
		new web3._extend.Modfod({
			name: 'getBlockByNumber',
			call: 'odf_getBlockByNumber',
			params: 2,
			inputFormatter: [web3._extend.formatters.inputBlockNumberFormatter, function (val) { return !!val; }]
		}),
		new web3._extend.Modfod({
			name: 'getBlockByHash',
			call: 'odf_getBlockByHash',
			params: 2,
			inputFormatter: [null, function (val) { return !!val; }]
		}),
		new web3._extend.Modfod({
			name: 'getRawTransaction',
			call: 'odf_getRawTransactionByHash',
			params: 1
		}),
		new web3._extend.Modfod({
			name: 'getRawTransactionFromBlock',
			call: function(args) {
				return (web3._extend.utils.isString(args[0]) && args[0].indexOf('0x') === 0) ? 'odf_getRawTransactionByBlockHashAndIndex' : 'odf_getRawTransactionByBlockNumberAndIndex';
			},
			params: 2,
			inputFormatter: [web3._extend.formatters.inputBlockNumberFormatter, web3._extend.utils.toHex]
		}),
		new web3._extend.Modfod({
			name: 'getProof',
			call: 'odf_getProof',
			params: 3,
			inputFormatter: [web3._extend.formatters.inputAddressFormatter, null, web3._extend.formatters.inputBlockNumberFormatter]
		}),
	],
	properties: [
		new web3._extend.Property({
			name: 'pendingTransactions',
			getter: 'odf_pendingTransactions',
			outputFormatter: function(txs) {
				var formatted = [];
				for (var i = 0; i < txs.length; i++) {
					formatted.push(web3._extend.formatters.outputTransactionFormatter(txs[i]));
					formatted[i].blockHash = null;
				}
				return formatted;
			}
		}),
	]
});
`

const MinerJs = `
web3._extend({
	property: 'miner',
	modfods: [
		new web3._extend.Modfod({
			name: 'start',
			call: 'miner_start',
			params: 1,
			inputFormatter: [null]
		}),
		new web3._extend.Modfod({
			name: 'stop',
			call: 'miner_stop'
		}),
		new web3._extend.Modfod({
			name: 'setEtherbase',
			call: 'miner_setEtherbase',
			params: 1,
			inputFormatter: [web3._extend.formatters.inputAddressFormatter]
		}),
		new web3._extend.Modfod({
			name: 'setExtra',
			call: 'miner_setExtra',
			params: 1
		}),
		new web3._extend.Modfod({
			name: 'setGasPrice',
			call: 'miner_setGasPrice',
			params: 1,
			inputFormatter: [web3._extend.utils.fromDecimal]
		}),
		new web3._extend.Modfod({
			name: 'setRecommitInterval',
			call: 'miner_setRecommitInterval',
			params: 1,
		}),
		new web3._extend.Modfod({
			name: 'getHashrate',
			call: 'miner_getHashrate'
		}),
	],
	properties: []
});
`

const NetJs = `
web3._extend({
	property: 'net',
	modfods: [],
	properties: [
		new web3._extend.Property({
			name: 'version',
			getter: 'net_version'
		}),
	]
});
`

const PersonalJs = `
web3._extend({
	property: 'personal',
	modfods: [
		new web3._extend.Modfod({
			name: 'importRawKey',
			call: 'personal_importRawKey',
			params: 2
		}),
		new web3._extend.Modfod({
			name: 'sign',
			call: 'personal_sign',
			params: 3,
			inputFormatter: [null, web3._extend.formatters.inputAddressFormatter, null]
		}),
		new web3._extend.Modfod({
			name: 'ecRecover',
			call: 'personal_ecRecover',
			params: 2
		}),
		new web3._extend.Modfod({
			name: 'openWallet',
			call: 'personal_openWallet',
			params: 2
		}),
		new web3._extend.Modfod({
			name: 'deriveAccount',
			call: 'personal_deriveAccount',
			params: 3
		}),
		new web3._extend.Modfod({
			name: 'signTransaction',
			call: 'personal_signTransaction',
			params: 2,
			inputFormatter: [web3._extend.formatters.inputTransactionFormatter, null]
		}),
		new web3._extend.Modfod({
			name: 'unpair',
			call: 'personal_unpair',
			params: 2
		}),
		new web3._extend.Modfod({
			name: 'initializeWallet',
			call: 'personal_initializeWallet',
			params: 1
		})
	],
	properties: [
		new web3._extend.Property({
			name: 'listWallets',
			getter: 'personal_listWallets'
		}),
	]
})
`

const RpcJs = `
web3._extend({
	property: 'rpc',
	modfods: [],
	properties: [
		new web3._extend.Property({
			name: 'modules',
			getter: 'rpc_modules'
		}),
	]
});
`

const ShhJs = `
web3._extend({
	property: 'shh',
	modfods: [
	],
	properties:
	[
		new web3._extend.Property({
			name: 'version',
			getter: 'shh_version',
			outputFormatter: web3._extend.utils.toDecimal
		}),
		new web3._extend.Property({
			name: 'info',
			getter: 'shh_info'
		}),
	]
});
`

const SwarmfsJs = `
web3._extend({
	property: 'swarmfs',
	modfods:
	[
		new web3._extend.Modfod({
			name: 'mount',
			call: 'swarmfs_mount',
			params: 2
		}),
		new web3._extend.Modfod({
			name: 'unmount',
			call: 'swarmfs_unmount',
			params: 1
		}),
		new web3._extend.Modfod({
			name: 'listmounts',
			call: 'swarmfs_listmounts',
			params: 0
		}),
	]
});
`

const TxpoolJs = `
web3._extend({
	property: 'txpool',
	modfods: [],
	properties:
	[
		new web3._extend.Property({
			name: 'content',
			getter: 'txpool_content'
		}),
		new web3._extend.Property({
			name: 'inspect',
			getter: 'txpool_inspect'
		}),
		new web3._extend.Property({
			name: 'status',
			getter: 'txpool_status',
			outputFormatter: function(status) {
				status.pending = web3._extend.utils.toDecimal(status.pending);
				status.queued = web3._extend.utils.toDecimal(status.queued);
				return status;
			}
		}),
	]
});
`

const AccountingJs = `
web3._extend({
	property: 'accounting',
	modfods: [
		new web3._extend.Property({
			name: 'balance',
			getter: 'account_balance'
		}),
		new web3._extend.Property({
			name: 'balanceCredit',
			getter: 'account_balanceCredit'
		}),
		new web3._extend.Property({
			name: 'balanceDebit',
			getter: 'account_balanceDebit'
		}),
		new web3._extend.Property({
			name: 'bytesCredit',
			getter: 'account_bytesCredit'
		}),
		new web3._extend.Property({
			name: 'bytesDebit',
			getter: 'account_bytesDebit'
		}),
		new web3._extend.Property({
			name: 'msgCredit',
			getter: 'account_msgCredit'
		}),
		new web3._extend.Property({
			name: 'msgDebit',
			getter: 'account_msgDebit'
		}),
		new web3._extend.Property({
			name: 'peerDrops',
			getter: 'account_peerDrops'
		}),
		new web3._extend.Property({
			name: 'selfDrops',
			getter: 'account_selfDrops'
		}),
	]
});
`

const LESJs = `
web3._extend({
	property: 'les',
	modfods:
	[
		new web3._extend.Modfod({
			name: 'getCheckpoint',
			call: 'les_getCheckpoint',
			params: 1
		}),
		new web3._extend.Modfod({
			name: 'clientInfo',
			call: 'les_clientInfo',
			params: 1
		}),
		new web3._extend.Modfod({
			name: 'priorityClientInfo',
			call: 'les_priorityClientInfo',
			params: 3
		}),
		new web3._extend.Modfod({
			name: 'setClientParams',
			call: 'les_setClientParams',
			params: 2
		}),
		new web3._extend.Modfod({
			name: 'setDefaultParams',
			call: 'les_setDefaultParams',
			params: 1
		}),
		new web3._extend.Modfod({
			name: 'addBalance',
			call: 'les_addBalance',
			params: 2
		}),
	],
	properties:
	[
		new web3._extend.Property({
			name: 'latestCheckpoint',
			getter: 'les_latestCheckpoint'
		}),
		new web3._extend.Property({
			name: 'checkpointContractAddress',
			getter: 'les_getCheckpointContractAddress'
		}),
		new web3._extend.Property({
			name: 'serverInfo',
			getter: 'les_serverInfo'
		}),
	]
});
`

const LESPayJs = `
web3._extend({
	property: 'lespay',
	modfods:
	[
		new web3._extend.Modfod({
			name: 'distribution',
			call: 'lespay_distribution',
			params: 2
		}),
		new web3._extend.Modfod({
			name: 'timeout',
			call: 'lespay_timeout',
			params: 2
		}),
		new web3._extend.Modfod({
			name: 'value',
			call: 'lespay_value',
			params: 2
		}),
	],
	properties:
	[
		new web3._extend.Property({
			name: 'requestStats',
			getter: 'lespay_requestStats'
		}),
	]
});
`
