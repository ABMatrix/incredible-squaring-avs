package main

import (
	"context"
	"fmt"
	"github.com/Layr-Labs/eigensdk-go/chainio/clients/eth"
	"github.com/Layr-Labs/eigensdk-go/chainio/clients/wallet"
	"github.com/Layr-Labs/eigensdk-go/chainio/txmgr"
	strategymanager "github.com/Layr-Labs/eigensdk-go/contracts/bindings/StrategyManager"
	sdklogging "github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/Layr-Labs/eigensdk-go/signerv2"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"math/big"
	"os"
)

func main() {
	logger, err := sdklogging.NewZapLogger(sdklogging.Development)
	if err != nil {
		panic(err)
	}

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

	pkWallet, err := wallet.NewPrivateKeyWallet(ethclient, signer, senderAddress, logger)
	txMgr := txmgr.NewSimpleTxManager(pkWallet, ethclient, logger, senderAddress)
	noSendTxOpts, err := txMgr.GetNoSendTxOpts()
	if err != nil {
		panic(err)
	}

	strategy := common.HexToAddress("0x3B2fB6Fb6d96fC823aFB1D2579411E2AFd3eC204")
	token := common.HexToAddress("0xf538309aCdcD0C4EC0707D664e9c2F7570026019")

	fmt.Println("senderAddress: ", senderAddress)
	fmt.Println("strategy: ", strategy)
	fmt.Println("token: ", token)

	fmt.Printf("%+v\n", noSendTxOpts)

	tx, err := contractStrategyManager.DepositIntoStrategy(noSendTxOpts, strategy, token, big.NewInt(5))
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
