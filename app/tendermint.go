package app

import (
	"reflect"

	"github.com/Fantom-foundation/go-lachesis/inter/idx"

	"github.com/ethereum/go-ethereum/common"
	eth "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/tendermint/tendermint/abci/types"
)

// InitChain implements ABCIApplication.InitChain.
// It should be Called once upon genesis.
func (a *App) InitChain(req types.RequestInitChain) types.ResponseInitChain {
	chain := a.config.Net.ChainInfo()
	if !reflect.DeepEqual(req, chain) {
		panic("incompatible chain")
	}

	a.Bootstrap()
	return types.ResponseInitChain{}
}

// DeliverTx implements ABCIApplication.DeliverTx
func (a *App) DeliverTx(req types.RequestDeliverTx) types.ResponseDeliverTx {
	const strict = false

	dagTx := BytesToDagTx(req.Tx)
	tx := dagTx.Transaction

	receipt, fee, skip, err := a.ctx.evmProcessor.
		ProcessTx(tx, a.ctx.txCount, a.ctx.gp, &a.ctx.header.GasUsed, a.ctx.header, a.ctx.statedb, vm.Config{}, strict)
	a.ctx.txCount++
	if !strict && (skip || err != nil) {
		return types.ResponseDeliverTx{
			Info:      "skipped",
			GasWanted: int64(tx.Gas()),
			GasUsed:   0,
		}
	}

	a.ctx.txs = append(a.ctx.txs, tx)
	a.ctx.receipts = append(a.ctx.receipts, receipt)
	a.ctx.totalFee.Add(a.ctx.totalFee, fee)
	a.store.AddDirtyOriginationScore(dagTx.Originator, fee)

	return types.ResponseDeliverTx{
		Info:      "ok",
		GasWanted: int64(tx.Gas()),
		GasUsed:   int64(receipt.GasUsed),
	}
}

// EndBlock implements ABCIApplication.EndBlock
func (a *App) EndBlock(
	req types.RequestEndBlock,
) (
	// TODO: return only types.ResponseEndBlock
	root common.Hash,
	receipts eth.Receipts,
	sealEpoch bool,
) {
	if a.ctx.block.Index != idx.Block(req.Height) {
		a.Log.Crit("missed block", "current", a.ctx.block.Index, "got", req.Height)
	}

	return a.endBlock()
	/*
		return types.ResponseEndBlock{
			ValidatorUpdates: types.ValidatorUpdates{
				types.ValidatorUpdate{},
				types.ValidatorUpdate{},
			},
		}
	*/
}
