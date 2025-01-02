package jup_client

type ClientOption func(*client) error

/*
Transaction options an rpc client accepts

type TransactionOpts struct{
  Encoding solana.EncodingType        `json:"encoding,omitemprt"`
  SkipPreflight bool                  `json:"skipPreflight,omitempty"`
  PreflightCommitment CommitmentType  `json:"preflightCommitment,omitempty"`
  MaxRetries *uint                    `json:"maxRetries"`
  MinContextSlot *uint64              `json:"minContextSlot"`
}
*/

func WithMaxRetries(maxRetries uint) ClientOption{
  return func(cl *client) error{
    cl.maxRetries = maxRetries
    return nil
  }
}
