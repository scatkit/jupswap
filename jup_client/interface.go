package jup_client
import(
  "context"
  "github.com/scatkit/pumpdexer/rpc"
  "github.com/scatkit/pumpdexer/solana"
)

type rpcService interface{
  SendTransactionWithOpts(
    ctx context.Context,
		transaction *solana.Transaction,
		opts rpc.TransactionOpts,
	) (signature solana.Signature, err error)
  GetLatestBlockhash(
		ctx context.Context,
		commitment rpc.CommitmentType,
	) (out *rpc.GetLatestBlockhashResult, err error)
	GetSignatureStatuses(
		ctx context.Context,
		searchTransactionHistory bool,
		transactionSignatures ...solana.Signature,
	) (out *rpc.GetSignatureStatusesResult, err error)
	GetTokenAccountBalance(
		ctx context.Context,
		account solana.PublicKey,
		commitment rpc.CommitmentType, // optional
	) (out *rpc.GetTokenAccountBalanceResult, err error)
	Close() error
}

