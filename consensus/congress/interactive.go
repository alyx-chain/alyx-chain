package congress

import (
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/params"
)

type chainContext struct {
	chainReader consensus.ChainHeaderReader
	engine      consensus.Engine
}

func newChainContext(chainReader consensus.ChainHeaderReader, engine consensus.Engine) *chainContext {
	return &chainContext{
		chainReader: chainReader,
		engine:      engine,
	}
}

// Engine retrieves the chain's consensus engine.
func (cc *chainContext) Engine() consensus.Engine {
	return cc.engine
}

// GetHeader returns the hash corresponding to their hash.
func (cc *chainContext) GetHeader(hash common.Hash, number uint64) *types.Header {
	return cc.chainReader.GetHeader(hash, number)
}

func getInteractiveABI() map[string]abi.ABI {
	abiMap := make(map[string]abi.ABI, 0)
	tmpABI, _ := abi.JSON(strings.NewReader(validatorsInteractiveABI))
	abiMap[validatorsContractName] = tmpABI
	tmpABI, _ = abi.JSON(strings.NewReader(punishInteractiveABI))
	abiMap[punishContractName] = tmpABI
	tmpABI, _ = abi.JSON(strings.NewReader(proposalInteractiveABI))
	abiMap[proposalContractName] = tmpABI

	return abiMap
}

// executeMsg executes transaction sent to system contracts.
func executeMsg(msg core.Message, state *state.StateDB, header *types.Header, chainContext core.ChainContext, chainConfig *params.ChainConfig) (ret []byte, err error) {
	// Set gas price to zero
	// NewEVM(blockCtx BlockContext, txCtx TxContext, statedb StateDB, chainConfig *params.ChainConfig, config Config) *EVM {
	context := core.NewEVMBlockContext(header, chainContext, nil)
	vmenv := vm.NewEVM(context, vm.TxContext{}, state, chainConfig, vm.Config{})
	// Call(caller ContractRef, addr common.Address, input []byte, gas uint64, value *big.Int) (ret []byte, leftOverGas uint64, err error) {
	ret, _, err = vmenv.Call(vm.AccountRef(msg.From), *msg.To, msg.Data, msg.GasLimit, msg.Value)

	if err != nil {
		return []byte{}, err
	}

	return ret, nil
}
