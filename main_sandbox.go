package main

import (
	"context"
  //"time"
  "fmt"
  "log"
  "encoding/json"
  "encoding/base64"

  bin "github.com/gagliardetto/binary"
  "github.com/scatkit/pumpdexer/rpc"
	jup "github.com/scatkit/jupswap/jupcore"
	//"github.com/scatkit/jupswap/jup_client"
	"github.com/scatkit/pumpdexer/solana"
	"github.com/davecgh/go-spew/spew"
  dec "github.com/shopspring/decimal"

)

const URL = "https://quote-api.jup.ag/v6"


func formatAmount(amountStr string, decimalNum uint) (formattedAmount int64, err error){
  amountFloat, err  := dec.NewFromString(amountStr)
  if err != nil{
    return 0, fmt.Errorf("invalid float string: %s", amountStr)
  }
  // Base 10^decimalNum 
  smallestUnitOfToken := dec.NewFromInt(10).Pow(dec.NewFromInt(int64(decimalNum)))
  
  return amountFloat.Mul(smallestUnitOfToken).IntPart(), nil
}

func main_tx() {
  userPubkey := solana.MustPubkeyFromBase58("CJMTJWF97jd3dspsN5qhPp4EpKBHMTnkRvDkpSHUWSGJ")
  // maybe not return an error
	jupiterClient, err := jup.NewClient(URL)
	if err != nil {
		panic(err)
	}
  solClient := rpc.New("https://api.mainnet-beta.solana.com")
  
  // Selling case
  var inputTokenMint = "So11111111111111111111111111111111111111112"
  var outputTokenMint = "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v"
  var inputAmount = "0.0001"
  var slippageBps = 150
  
  var decimalPrec uint
  if inputTokenMint == "So11111111111111111111111111111111111111112"{
    decimalPrec = uint(9) // should do this as constant
  } else{
    tokenAccountBalanceResponse, err := solClient.GetTokenSupply(
      context.Background(), 
      solana.MustPubkeyFromBase58(inputTokenMint),
      rpc.CommitmentFinalized)
      if err != nil{
        log.Fatal(err)
      }
      decimalPrec = uint(tokenAccountBalanceResponse.Value.Decimals) // casts uint8 to uint
    } 
  formattedUserInputAmount, err := formatAmount(inputAmount, decimalPrec)
  if err != nil{
    log.Fatal(err)
  }

  // JUP 
  //var asLegacy = true
  quoteResp, err := jupiterClient.GetQuoteWithResponse(context.Background(), &jup.GetQuoteParams{
    InputMint:   inputTokenMint,
    OutputMint:  outputTokenMint, 
    Amount:      formattedUserInputAmount,
    SlippageBps: &slippageBps,
    //AsLegacyTransaction:  &asLegacy,
  })

  if err != nil {
    panic(err)
  }
  if quoteResp.JSON200 == nil {
    panic("invalid GetQuoteResponse response")
  }
  quote := quoteResp.JSON200
  
  dynamicComputeUnitLimit := true 
  prioritizationFeeLamports := jup.SwapRequest_PrioritizationFeeLamports{}
  jtip := map[string]int64{
    "jitoTipLamports": 10000,
  }
  bts,err := json.Marshal(jtip)
  if err != nil{
    panic(err)
  }
  if err = prioritizationFeeLamports.UnmarshalJSON(bts); err != nil {
		panic(err)
	}
  swapResp, err := jupiterClient.PostSwapWithResponse(context.TODO(), jup.SwapRequest{
    PrioritizationFeeLamports: &prioritizationFeeLamports,
  	QuoteResponse:             *quote,
  	UserPublicKey:             userPubkey.String(),
    //AsLegacyTransaction:       &asLegacy, 
    DynamicComputeUnitLimit: &dynamicComputeUnitLimit,
  })
  
  if err != nil{
    log.Fatal(err)
  }
  
  if swapResp.JSON200 == nil {
  	panic("invalid PostSwapWithResponse{} response")
  }
  swap := swapResp.JSON200
  //spew.Dump(swap)
  txBytes, err := base64.StdEncoding.DecodeString(swap.SwapTransaction)
  if err != nil {
   	log.Fatal(err)
   }
  
  tx, err := solana.TransactionFromDecoder(bin.NewBinDecoder(txBytes))
  if err != nil{
    log.Fatal(err)
  }
  
  spew.Dump(tx)
  //wallet,err := jup_client.NewWalletFromPrivateKeyBase58("5rg7jXrAYXoAYt1ARV1RzuRFCsH948MyjMjKVG8Kiw7pdZZ7QBjuJnEfufvukPJ5hLyRHUXkPBuc9mP7AS35i5yC")
  //if err != nil{
  //  log.Fatal(err)
  //}
  //jupsolClient, err := jup_client.NewClient(wallet, "https://api.mainnet-beta.solana.com")
  //fmt.Println(jupsolClient)
  //Sign and send the transaction.
	//signedTx, err := jupsolClient.SendTransactionOnChain(context.TODO(),swap.SwapTransaction)
	//if err != nil {
	//	panic(err)
	//}
  //
  //fmt.Println(signedTx)
  
 // _, err = jupsolClient.CheckSignature(context.TODO(), signedTx)
 // if err != nil {
 // 	panic(err)
 // }
 // 
 // //Sending tx
 // 
 // txMessageBytes, err := tx.Message.MarshalBinary()
 // if err != nil {
 // 	fmt.Errorf("could not serialize transaction: %w", err)
 // }
 // signature, err := privkey.Sign(txMessageBytes)
 // if err != nil {
 // 	fmt.Errorf("could not sign transaction: %w", err)
 // }
 // tx.Signatures = []solana.Signature{signature} 
 // spew.Dump(tx)
 //   privkey := solana.MustPrivkeyFromBase58("5rg7jXrAYXoAYt1ARV1RzuRFCsH948MyjMjKVG8Kiw7pdZZ7QBjuJnEfufvukPJ5hLyRHUXkPBuc9mP7AS35i5yC")
 //   _, err = tx.Sign(func(key solana.PublicKey) *solana.PrivateKey {
 //  	if privkey.PublicKey().Equals(key) {
 //  		return &privkey
 //  	}
 //  	return nil
 //  })
 //  if err != nil {
 //  	log.Fatalf("Failed to sign transaction: %v", err)
 //  }
 // 
 //   bhash, err := solClient.GetLatestBlockhash(context.Background(), rpc.CommitmentFinalized)
 //   if err != nil{
 //     log.Fatal(err)
 //   }
 //   tx.Message.RecentBlockhash = bhash.Value.Blockhash
 // 
 // 
 //   spew.Dump(tx)
 //   sig, err := solClient.SendTransaction(context.Background(), tx)
 //   if err != nil{
 //     log.Fatal(err)
 //   }
 //   fmt.Println("Signature",sig)
}

