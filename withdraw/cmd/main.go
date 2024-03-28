package main

import (
	"context"
	"fmt"
	"github.com/Layr-Labs/eigensdk-go/chainio/clients/eth"
	"github.com/Layr-Labs/eigensdk-go/chainio/clients/wallet"
	"github.com/Layr-Labs/eigensdk-go/chainio/txmgr"
	delegationManger "github.com/Layr-Labs/eigensdk-go/contracts/bindings/DelegationManager"
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

	delegationManager1 := common.HexToAddress("0x193e0aaD287FD43fBd14A55e44A5dCFC31B0e432")

	ethclient, err := eth.NewClient("https://ethereum-holesky-rpc.publicnode.com")
	if err != nil {
		panic(err)
	}

	signer := func(ctx context.Context, address common.Address) (bind.SignerFn, error) {
		return signerv2.PrivateKeySignerFn(ecdsaPrivateKey, big.NewInt(17000))
	}

	contractDelegationManager, err := delegationManger.NewContractDelegationManager(delegationManager1, ethclient)
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

	queuedWithdrawalParams := delegationManger.IDelegationManagerQueuedWithdrawalParams{
		Strategies: []common.Address{strategy},
		Shares:     []*big.Int{big.NewInt(5)},
		Withdrawer: senderAddress,
	}

	tx, err := contractDelegationManager.QueueWithdrawals(noSendTxOpts, []delegationManger.IDelegationManagerQueuedWithdrawalParams{queuedWithdrawalParams})
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
