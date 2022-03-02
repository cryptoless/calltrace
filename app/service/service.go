package service

import (
	"fmt"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/state/snapshot"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/trie"

	"github.com/cryptoless/chain-raw-api-server/callserver"
	"github.com/cryptoless/chain-raw-api-server/message"
)

var Service *callserver.ServiceRegister
var once sync.Once

func init() {
	once.Do(func() {
		Service = callserver.NewService()
		adp := &Trace{}
		Service.Registration(adp)
	})
}

type TransactionArgs struct {
	From                 *common.Address `json:"from"`
	To                   *common.Address `json:"to"`
	Gas                  *hexutil.Uint64 `json:"gas"`
	GasPrice             *hexutil.Big    `json:"gasPrice"`
	MaxFeePerGas         *hexutil.Big    `json:"maxFeePerGas"`
	MaxPriorityFeePerGas *hexutil.Big    `json:"maxPriorityFeePerGas"`
	Value                *hexutil.Big    `json:"value"`
	Nonce                *hexutil.Uint64 `json:"nonce"`

	// We accept "data" and "input" for backwards-compatibility reasons.
	// "input" is the newer name and should be preferred by clients.
	// Issue detail: https://github.com/ethereum/go-ethereum/issues/15628
	Data  *hexutil.Bytes `json:"data"`
	Input *hexutil.Bytes `json:"input"`

	// Introduced by AccessListTxType transaction.
	AccessList *types.AccessList `json:"accessList,omitempty"`
	ChainID    *hexutil.Big      `json:"chainId,omitempty"`
}

type dataBase struct {
}

func newDataBase() state.Database {
	//todo: db
	db := rawdb.NewMemoryDatabase()
	return state.NewDatabaseWithConfig(db, &trie.Config{
		// Cache:     cacheConfig.TrieCleanLimit,
		// Journal:   cacheConfig.TrieCleanJournal,
		// Preimages: cacheConfig.Preimages,
	})
}

///

type callTracerTest struct {
	Genesis *core.Genesis `json:"genesis"`
	Context *callContext  `json:"context"`
	Input   string        `json:"input"`
}
type callContext struct {
	Number     math.HexOrDecimal64   `json:"number"`
	Difficulty *math.HexOrDecimal256 `json:"difficulty"`
	Time       math.HexOrDecimal64   `json:"timestamp"`
	GasLimit   math.HexOrDecimal64   `json:"gasLimit"`
	Miner      common.Address        `json:"miner"`
}

func MakePreState(db ethdb.Database, accounts core.GenesisAlloc, snapshotter bool) (*snapshot.Tree, *state.StateDB) {
	sdb := state.NewDatabase(db)
	statedb, _ := state.New(common.Hash{}, sdb, nil)
	for addr, a := range accounts {
		statedb.SetCode(addr, a.Code)
		statedb.SetNonce(addr, a.Nonce)
		statedb.SetBalance(addr, a.Balance)
		for k, v := range a.Storage {
			statedb.SetState(addr, k, v)
		}
	}
	// Commit and re-open to start with a clean state.
	root, _ := statedb.Commit(false)

	var snaps *snapshot.Tree
	if snapshotter {
		snaps, _ = snapshot.New(db, sdb.TrieDB(), 1, root, false, true, false)
	}
	statedb, _ = state.New(root, sdb, snaps)
	return snaps, statedb
}

// func getBlockContractState() core.GenesisAlloc {
// 	alloc := core.GenesisAlloc{}

// 	addr := common.HexToAddress("0x2a98c5f40bfa3dee83431103c535f6fae9a8ad38")
// 	ga := core.GenesisAccount{
// 		Balance: big.NewInt(0),
// 		Code:    common.Hex2Bytes(code),
// 		Nonce:   uint64(1),
// 		Storage: make(map[common.Hash]common.Hash),
// 	}
// 	ga.Storage[common.HexToHash("0000000000000000000000000000000000000000000000000000000000000002")] = common.HexToHash("0000000000000000000000002cccf5e0538493c235d1c5ef6580f77d99e91396")
// 	alloc[addr] = ga
// 	return alloc
// }
// func getConf() *params.ChainConfig {
// 	return nil
// }
// func getContext() *callContext {
// 	return nil
// }

type Trace struct {
}

func (e *Trace) Call(test *callTracerTest) ([]byte, message.Error) {

	// test := new(callTracerTest)
	// json.Unmarshal([]byte(calljs), test)

	tx := new(types.Transaction)

	err := rlp.DecodeBytes(common.FromHex(test.Input), tx)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(common.Bytes2Hex(tx.Data()))
	fmt.Println(tx.To().Hex())
	_, statedb := MakePreState(rawdb.NewMemoryDatabase(), test.Genesis.Alloc, false)
	blockCtx := vm.BlockContext{
		CanTransfer: core.CanTransfer,
		Transfer:    core.Transfer,
		Coinbase:    test.Context.Miner,
		BlockNumber: new(big.Int).SetUint64(uint64(test.Context.Number)),
		Time:        new(big.Int).SetUint64(uint64(test.Context.Time)),
		Difficulty:  (*big.Int)(test.Context.Difficulty),
		GasLimit:    uint64(test.Context.GasLimit),
	}

	signer := types.MakeSigner(test.Genesis.Config, new(big.Int).SetUint64(uint64(test.Context.Number)))
	msg, err := tx.AsMessage(signer, nil)
	fmt.Println(msg.From().Hex())
	if err != nil {
		panic(err)
	}

	evm := vm.NewEVM(blockCtx, core.NewEVMTxContext(msg), statedb, test.Genesis.Config, vm.Config{})

	gp := new(core.GasPool).AddGas(math.MaxUint64)
	result, err := core.ApplyMessage(evm, msg, gp)
	if err != nil {
		return nil, message.InternalError(err.Error())
	}

	return result.Return(), nil
}
