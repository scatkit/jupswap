package jup_client

import (
	"context"
	"fmt"
	"github.com/scatkit/pumpdexer/rpc"
	"github.com/scatkit/pumpdexer/solana"
	"github.com/shopspring/decimal"
  //"github.com/davecgh/go-spew/spew"

)

type Client interface {
	SendTransactionOnChain(ctx context.Context, txBasa64 string) (TxID, error)
  CheckSignature(ctx context.Context, signedTx TxID) (bool, error)
}

const defaultMaxRetries = uint(20)

type TxID string

type TokenAccount struct {
	Amount  decimal.Decimal
	Decimal uint8
}

type client struct {
	maxRetries uint
	clientRPC  rpcMethod // e.g SendTransaction
	wallet     Wallet    // user's wallet
}

func newClient(wallet Wallet, rpcEndpoint string, opts ...ClientOption) (*client, error) {
	cl := &client{
		maxRetries: defaultMaxRetries,
		wallet:     wallet,
	}

	for _, opt := range opts {
		if err := opt(cl); err != nil {
			return nil, fmt.Errorf("could not apply option")
		}
	}

	if cl.clientRPC == nil {
		if rpcEndpoint == "" {
			return nil, fmt.Errorf("endpoint url is required")
		}
		cl.clientRPC = rpc.New(rpcEndpoint)
	}

	return cl, nil
}

func NewClient(wallet Wallet, rpcEndpoint string, opts ...ClientOption) (Client, error) {
  return newClient(wallet, rpcEndpoint, opts...)
}

func (cl client) SendTransactionOnChain(ctx context.Context, txBase64 string) (TxID, error) {
	latestBlockhash, err := cl.clientRPC.GetLatestBlockhash(ctx, rpc.CommitmentFinalized)
	if err != nil {
		return "", fmt.Errorf("could not get the latest blockhash: %w", err)
	}

	tx, err := NewTransactionFromBase64(txBase64)
	if err != nil {
		return "", fmt.Errorf("could not deseeialize swap transaction: %w", err)
	}
	tx.Message.RecentBlockhash = latestBlockhash.Value.Blockhash

	// Signs the transaction with the wallet's private key
	tx, err = cl.wallet.SignTransaction(tx)
	if err != nil {
		return "", fmt.Errorf("could not sign swap transaction: %w", err)
	}
  
   //spew.Dump(tx)

	sig, err := cl.clientRPC.SendTransactionWithOpts(ctx, &tx, rpc.TransactionOpts{
		MaxRetries:          &cl.maxRetries,
		MinContextSlot:      &latestBlockhash.Context.Slot,
		PreflightCommitment: rpc.CommitmentConfirmed,
	})

	if err != nil {
		return "", fmt.Errorf("could not send the transaction: %w", err)
	}

	return TxID(sig.String()), nil
}

func (cl client) CheckSignature(ctx context.Context, tx TxID) (bool, error) {
	sig, err := solana.SignatureFromBase58(string(tx))
	if err != nil {
		return false, fmt.Errorf("could not decode the transaction from base58")
	}
  
	status, err := cl.clientRPC.GetSignatureStatuses(ctx, false, sig)
	if err != nil {
		return false, fmt.Errorf("could not get signature status: %w", err)
	}

	if len(status.Value) == 0 {
		return false, fmt.Errorf("could not confirm transaction: no valid status")
	}

	if status.Value[0] == nil || status.Value[0].ConfirmationStatus != rpc.ConfirmationStatusFinalized {
		return false, fmt.Errorf("transaction not finalized yet")
	}

	if status.Value[0].Err != nil {
		return true, fmt.Errorf("transaction confirmed with error: %s", status.Value[0].Err)
	}

	return true, nil
}


