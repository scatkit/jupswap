package jup_client
import(
  "github.com/scatkit/pumpdexer/solana"
  "encoding/base64"
	"fmt"

  bin "github.com/gagliardetto/binary"
)

func NewTransactionFromBase64(trxBase64 string) (solana.Transaction, error){
  txBytes,err := base64.StdEncoding.DecodeString(trxBase64)
  if err != nil{
    return solana.Transaction{}, fmt.Errorf("could not decode transaction: %w", err)
  }
  tx, err := solana.TransactionFromDecoder(bin.NewBinDecoder(txBytes))
  if err != nil{
    return solana.Transaction{}, fmt.Errorf("could not deserialize transaction from Decoder: %w", err)
  }
  
  return *tx, nil
}
