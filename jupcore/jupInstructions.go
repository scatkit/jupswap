package jupcore
import(
  "context"
  "encoding/json"
  "net/http"
  "net/url"
  "fmt"
  "bytes"
  "io"
  "strings"
)

type AccountMeta struct{
  IsSigner   bool   `json:"isSigner"`
	IsWritable bool   `json:"isWritable"`
	Pubkey     string `json:"pubkey"`
}

// Instruction defines model for Instruction.
type Instruction struct {
	Accounts  []AccountMeta `json:"accounts"`
	Data      string        `json:"data"`
	ProgramId string        `json:"programId"`
}

type SwapInstructionsRequest = SwapRequest

type SwapInstructionsResponse struct{
  // AddressLookupTableAddresses The lookup table addresses that you can use if you are using versioned transaction.
	AddressLookupTableAddresses []string     `json:"addressLookupTableAddresses"`
	CleanupInstruction          *Instruction `json:"cleanupInstruction,omitempty"`

	// ComputeBudgetInstructions The necessary instructions to setup the compute budget.
	ComputeBudgetInstructions []Instruction `json:"computeBudgetInstructions"`

	// SetupInstructions Setup missing ATA for the users.
	SetupInstructions      []Instruction `json:"setupInstructions"`
	SwapInstruction        Instruction   `json:"swapInstruction"`
	TokenLedgerInstruction *Instruction  `json:"tokenLedgerInstruction,omitempty"`
}

type PostSwapInstructionsResponse struct{
  Body          []byte
  HTTPResponse  *http.Response
  JSON200       *SwapInstructionsResponse
}

func (cl *JupClient) PostSwapInstructionsWithResponse(ctx context.Context, reqBody SwapInstructionsRequest, reqEditors ...RequestEditorFunction,
) (*PostSwapInstructionsResponse, error) {
	resp, err := cl.PostSwapInstructions(ctx, reqBody, reqEditors...)
	if err != nil {
		return nil, err
	}
  
  bodyBytes, err := io.ReadAll(resp.Body)
  defer func(){resp.Body.Close()}()
  if err != nil{
    return nil, err
  }
  
  response := &PostSwapInstructionsResponse{
    Body: bodyBytes,
    HTTPResponse: resp,
  }

  switch {
	case strings.Contains(resp.Header.Get("Content-Type"), "json") && resp.StatusCode == 200:
		var dest SwapInstructionsResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest
	}

	return response, nil
}

func (cl *JupClient) PostSwapInstructions(ctx context.Context, reqBody SwapInstructionsRequest, reqEditors ...RequestEditorFunction,
) (*http.Response, error){ 
  req, err := NewPostSwapInstructionsRequest(ctx, cl.Endpoint, reqBody)
  if err != nil{
    return nil, err
  }
  req = req.WithContext(ctx)
  if err = cl.applyEditors(ctx, req, reqEditors); err != nil{
    return nil, err
  }

  return cl.HTTPClient.Do(req)
} 

func NewPostSwapInstructionsRequest(ctx context.Context, endpointURL string, reqBody SwapInstructionsRequest,
) (*http.Request, error){

  reqBytes, err := json.Marshal(reqBody)
  if err != nil{
    return nil, err
  }

  Endpoint, err := url.Parse(endpointURL)
  if err != nil{
    return nil, err
  } 

  operationPath := fmt.Sprintf("/swap-instructions")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := Endpoint.Parse(operationPath)
	if err != nil {
		return nil, err
	}
  
  //endpointURL += "/swap"
  request, err := http.NewRequestWithContext(ctx, "POST", queryURL.String(), bytes.NewReader(reqBytes))
  
  request.Header.Set("Content-Type", "application/json")
  request.Header.Set("Accept", "application/json")
  
  return request, nil
}

