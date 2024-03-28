package main

import (
	"context"
	"fmt"
	"github.com/Layr-Labs/eigensdk-go/chainio/clients/eth"
	"github.com/Layr-Labs/eigensdk-go/chainio/clients/wallet"
	"github.com/Layr-Labs/eigensdk-go/chainio/txmgr"
	strategymanager "github.com/Layr-Labs/eigensdk-go/contracts/bindings/StrategyManager"
	"github.com/Layr-Labs/eigensdk-go/signerv2"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"math/big"
	"os"
)

func main() {
	ecdsaKey := os.Getenv("ECDSA_KEY")

	if ecdsaKey == "" {
		panic("ecdsaKey not set")
	}

	ecdsaPrivateKey, err := crypto.HexToECDSA(ecdsaKey)
	if err != nil {
		panic("Cannot parse ecdsa private key")
	}
	senderAddress := crypto.PubkeyToAddress(ecdsaPrivateKey.PublicKey)

	strategyManagerAddr := common.HexToAddress("0xed9df70f6E01F6eAD317346855808980c457EbeC")

	ethclient, err := eth.NewClient("https://ethereum-holesky-rpc.publicnode.com")
	if err != nil {
		panic(err)
	}

	signer := func(ctx context.Context, address common.Address) (bind.SignerFn, error) {
		return signerv2.PrivateKeySignerFn(ecdsaPrivateKey, big.NewInt(17000))
	}

	contractStrategyManager, err := strategymanager.NewContractStrategyManager(strategyManagerAddr, ethclient)
	if err != nil {
		panic(err)
	}

	pkWallet, err := wallet.NewPrivateKeyWallet(ethclient, signer, senderAddress, nil)
	txMgr := txmgr.NewSimpleTxManager(pkWallet, ethclient, nil, senderAddress)
	noSendTxOpts, err := txMgr.GetNoSendTxOpts()
	if err != nil {
		panic(err)
	}

	strategy := common.HexToAddress("0xb35a9763CC9DA3F1fabD425C50a79b6f58295ADc")
	token := common.HexToAddress("0x94373a4919B3240D86eA41593D5eBa789FEF3848")

	fmt.Println("senderAddress: ", senderAddress)
	fmt.Println("strategy: ", strategy)
	fmt.Println("token: ", token)

	fmt.Printf("%+v\n", noSendTxOpts)

	tx, err := contractStrategyManager.DepositIntoStrategy(noSendTxOpts, strategy, token, big.NewInt(100000000000000))
	if err != nil {
		panic(err)
	}

	fmt.Println(tx.Hash())

	send, err := txMgr.Send(context.Background(), tx)
	if err != nil {
		panic(err)
	}

	fmt.Println(send.TxHash)
}
