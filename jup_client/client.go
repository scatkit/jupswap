package jup_client
import(
  //"github.com/scatkit/pumpdexer/solana"
  "github.com/scatkit/pumpdexer/rpc" 
  "context"
  "fmt"
  "github.com/shopspring/decimal"
)

type Client interface{
  SendTransactionOnChain(ctx context.Context, txBasa64 string) (TxID, error)
}

const defaultMaxRetries = uint(20)

type TxID string

type TokenAccount struct{
  Amount decimal.Decimal
  Decimal uint8
}

type client struct{
  maxRetries uint
  clientRPC  rpcService
  wallet     Wallet
}

func newClient(wallet Wallet, rpcEndpoint string, opts ...ClientOption) (*client, error){
  cl := &client{
    maxRetries: defaultMaxRetries,
    wallet:     wallet,
  }
  
  for _, opt := range opts{
    if err := opt(cl); err != nil{
      return nil, fmt.Errorf("could not apply option")
    }
  }
  
  if cl.clientRPC == nil{
    if rpcEndpoint == ""{
      return nil, fmt.Errorf("endpoint url is required")
    }
    cl.clientRPC = rpc.New(rpcEndpoint)
  }
  
  return cl, nil 
}

func NewClient(wallet Wallet, rpcEndpoint string, opts ...ClientOption) (Client, error){
  return newClient(wallet, rpcEndpoint, opts...)
}

func (cl client) SendTransactionOnChain(ctx context.Context, txBase64 string) (TxID, error){
  latestBlockhash, err := cl.clientRPC.GetLatestBlockhash(ctx, "")
  if err != nil{
    return "", fmt.Errorf("could not get the latest blockhash: %w", err)
  }
  
  tx, err := NewTransactionFromBase64(txBase64)
  if err != nil{
    return "", fmt.Errorf("could not deseeialize swap transaction: %w", err)
  }
  tx.Message.RecentBlockhash = latestBlockhash.Value.Blockhash 

  // Signs the transaction with the wallet's private key
  tx, err = cl.wallet.SignTransaction(tx)
  if err != nil{
    return "", fmt.Errorf("could not sign swap transaction: %w", err)
  } 
  
  sig, err := cl.clientRPC.SendTransactionWithOpts(ctx, &tx, rpc.TransactionOpts{
    MaxRetries:           &cl.maxRetries,  
    MinContextSlot:       &latestBlockhash.Context.Slot,
    PreflightCommitment:  rpc.CommitmentProcessed,
  })
  
  if err != nil{
    return "", fmt.Errorf("could not send the transaction: %w", err)
  }
  
  return TxID(sig.String()), nil
}

