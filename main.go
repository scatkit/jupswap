package main

import(
  jup "github.com/scatkit/jupswap/jupcore"
  //"fmt"
  "context"
  "github.com/davecgh/go-spew/spew"
) 

const URL = "https://quote-api.jup.ag/v6"

func main(){
  jupClient, err := jup.NewClient(URL)
  if err != nil{
    panic(err)
  }
  
  slippageBps := 200
  
  quoteResp, err := jupClient.GetQuoteWithResponse(context.Background(), &jup.GetQuoteParams{
    InputMint:   "So11111111111111111111111111111111111111112",
		OutputMint:  "BYXgSMha7DJkMuAA5ZD9UEFoRfqExEF3fJfDv7qZpump",
		Amount:      1000500, 
		SlippageBps: &slippageBps,
  })
  
  if err != nil{
    panic(err)
  }
  if quoteResp.JSON200 == nil{
    panic("invalid GetQuoteResponse response")
  }

  quote := quoteResp.JSON200
   
}
