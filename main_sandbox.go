package main

import (
	"context"
  //"time"
  "fmt"
  "log"

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

func main() {
  userPubkey := solana.MustPubkeyFromBase58("CJMTJWF97jd3dspsN5qhPp4EpKBHMTnkRvDkpSHUWSGJ")
  // maybe not return an error
	jupiterClient, err := jup.NewClient(URL)
	if err != nil {
		panic(err)
	}
  solClient := rpc.New("https://api.mainnet-beta.solana.com")
  
  // Selling case
  var inputTokenMint = "6p6xgHyF7AeE6TZkSmFsko444wqoP15icUSqi2jfGiPN"
  var outputTokenMint = "So11111111111111111111111111111111111111112"
  var inputAmount = "1.394"
  var slippageBps = 250
  
  var decimalPrec uint
  if inputTokenMint == "So11111111111111111111111111111111111111112"{
    decimalPrec = uint(1e9) // should do this as constant
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
  quoteResp, err := jupiterClient.GetQuoteWithResponse(context.Background(), &jup.GetQuoteParams{
    InputMint:   inputTokenMint,
    OutputMint:  outputTokenMint, 
    Amount:      formattedUserInputAmount,
    SlippageBps: &slippageBps,
  })

  if err != nil {
    panic(err)
  }
  if quoteResp.JSON200 == nil {
    panic("invalid GetQuoteResponse response")
  }
  quote := quoteResp.JSON200
  //spew.Dump(quote)
  
  swapInstrResp, err := jupiterClient. PostSwapInstructionsWithResponse(context.TODO(), jup.SwapInstructionsRequest{
  	QuoteResponse:             *quote,
  	UserPublicKey:             userPubkey.String(),
  })
  
  if swapInstrResp.JSON200 == nil {
  	panic("invalid PostSwapWithResponse{} response")
  }
  swapInstr := swapInstrResp.JSON200

  spew.Dump(swapInstr)
}

