package jup_client
import (
  "fmt"
  "github.com/scatkit/pumpdexer/solana"
)

type Wallet struct{
  *solana.Wallet
}

func NewWalletFromPrivateKeyBase58(privateKey string) (Wallet, error){
  // Decoding base58 string
  address, err := solana.WalletFromPrivateKeyBase58(privateKey)
  if err != nil{
    return Wallet{}, err
  }
  return Wallet{address}, nil
}

func (w Wallet) SignTransaction(tx solana.Transaction) (solana.Transaction, error){ 
  txMsgBytes, err := tx.Message.MarshalBinary()
  if err != nil{
    return solana.Transaction{}, fmt.Errorf("could not serialize transaction via Wallet %s: %w", w.PrivateKey, err)
  }
  
  sig, err := w.PrivateKey.Sign(txMsgBytes)
  if err != nil{
    return solana.Transaction{}, fmt.Errorf("could not sign transaction via Wallet %s: %w", w.PrivateKey, err)
  }
  tx.Signatures = []solana.Signature{sig}
  
  return tx, nil
}

