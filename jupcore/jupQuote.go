package jupcore

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
  "encoding/json"

	"github.com/oapi-codegen/runtime"
	//"github.com/davecgh/go-spew/spew"
)

// Defines values for swap mode
type SwapMode string

const (
	SwapModeExactIn  SwapMode = "ExactIn"
	SwapModeExactOut SwapMode = "ExactOut"
)

// Defines values for GetQuoteParamsSwapMode.
type GetQuoteParamsSwapMode string

const (
	ExactIn  GetQuoteParamsSwapMode = "ExactIn"
	ExactOut GetQuoteParamsSwapMode = "ExactOut"
)

type PlatformFee struct {
	Amount *string `json:"amount,omitempty"`
	FeeBps *int32  `json:"feeBps,omitempty"`
}

type RoutePlanStep struct {
	Percent  int32    `json:"percent"`
	SwapInfo SwapInfo `json:"swapInfo"`
}

type SwapInfo struct {
	AmmKey     string  `json:"ammKey"`
	FeeAmount  string  `json:"feeAmount"`
	FeeMint    string  `json:"feeMint"`
	InAmount   string  `json:"inAmount"`
	InputMint  string  `json:"inputMint"`
	Label      *string `json:"label,omitempty"`
	OutAmount  string  `json:"outAmount"`
	OutputMint string  `json:"outputMint"`
}

type GetQuoteParams struct {
	// InputMint Input token mint address
	InputMint string `form:"inputMint" json:"inputMint"`

	// OutputMint Output token mint address
	OutputMint string `form:"outputMint" json:"outputMint"`

	// Amount The amount to swap, have to factor in the token decimals.
	Amount int64 `form:"amount" json:"amount"`

	// SlippageBps The slippage in basis points, 1 basis point is 0.01%.
  // If the output token amount exceeds the slippage then the swap transaction will fail.
	SlippageBps *int `form:"slippageBps,omitempty" json:"slippageBps,omitempty"`

	// DynamicSlippage Set to true to indicate the usage of dynamic slippage.
	DynamicSlippage *bool `form:"dynamicSlippage,omitempty" json:"dynamicSlippage,omitempty"`

	// AutoSlippage Automatically calculate the slippage based on pairs.
	AutoSlippage *bool `form:"autoSlippage,omitempty" json:"autoSlippage,omitempty"`

	// AutoSlippageCollisionUsdValue Automatic slippage collision value.
	AutoSlippageCollisionUsdValue *int `form:"autoSlippageCollisionUsdValue,omitempty" json:"autoSlippageCollisionUsdValue,omitempty"`

	// MaxAutoSlippageBps Max slippage in basis points for auto slippage calculation. Default is 400.
	MaxAutoSlippageBps *int `form:"maxAutoSlippageBps,omitempty" json:"maxAutoSlippageBps,omitempty"`

	// (ExactIn or ExactOut). Defaults to ExactIn. ExactOut is for supporting use cases where you need an exact token amount, 
  // like payments. In this case the slippage is on the input token.
  // Ex: I need to spend exactly 1 SOL to get x Bonk. I need get exactly 10 Bonk for sending x SOL.
	SwapMode *GetQuoteParamsSwapMode `form:"swapMode,omitempty" json:"swapMode,omitempty"`

	// Default is that all DEXes are included. You can pass in the DEXes that you want to include only and separate them by `,`. 
  // You can check out the full list [here](https://quote-api.jup.ag/v6/program-id-to-label).
	Dexes *[]string `form:"dexes,omitempty" json:"dexes,omitempty"`

	// Default is that all DEXes are included.
  // You can pass in the DEXes that you want to exclude and separate them by `,`. 
  // You can check out the full list [here](https://quote-api.jup.ag/v6/program-id-to-label).
	ExcludeDexes *[]string `form:"excludeDexes,omitempty" json:"excludeDexes,omitempty"`

	// Restrict intermediate tokens to a top token set that has stable liquidity. 
  // This will help to ease potential high slippage error rate when swapping with minimal impact on pricing.
	RestrictIntermediateTokens *bool `form:"restrictIntermediateTokens,omitempty" json:"restrictIntermediateTokens,omitempty"`

	// OnlyDirectRoutes Default is false. Direct Routes limits Jupiter routing to single hop routes only.
	OnlyDirectRoutes *bool `form:"onlyDirectRoutes,omitempty" json:"onlyDirectRoutes,omitempty"`

	// AsLegacyTransaction Default is false. Instead of using versioned transaction, this will use the legacy transaction.
	AsLegacyTransaction *bool `form:"asLegacyTransaction,omitempty" json:"asLegacyTransaction,omitempty"`

	// PlatformFeeBps If you want to charge the user a fee, you can specify the fee in BPS. Fee % is taken out of the output token.
	PlatformFeeBps *int `form:"platformFeeBps,omitempty" json:"platformFeeBps,omitempty"`

	// MaxAccounts Rough estimate of the max accounts to be used for the quote, so that you can compose with your own accounts
	MaxAccounts *int `form:"maxAccounts,omitempty" json:"maxAccounts,omitempty"`

	// MinimizeSlippage Default is false. Miminize slippage attempts to find routes with lower slippage.
	MinimizeSlippage *bool `form:"minimizeSlippage,omitempty" json:"minimizeSlippage,omitempty"`

	// PreferLiquidDexes Default is false. Enabling it would only consider markets with high liquidity to reduce slippage.
	PreferLiquidDexes *bool `form:"preferLiquidDexes,omitempty" json:"preferLiquidDexes,omitempty"`

	// Default is false.
  // Uses categorized top token lists as intermediate tokens to optimize routing paths, replacing the old static top token list.
  // This helps achieve better pricing while maintaining route reliability.
	TokenCategoryBasedIntermediateTokens *bool `form:"tokenCategoryBasedIntermediateTokens,omitempty" json:"tokenCategoryBasedIntermediateTokens,omitempty"`
}

func (cl *JupClient) GetQuoteWithResponse(ctx context.Context, params *GetQuoteParams, reqEditors ...RequestEditorFunction) (*GetQuoteResponse, error) {
	resp, err := cl.GetQuote(ctx, params, reqEditors...)
	if err != nil {
		return nil, err
	}

	return ParseGetQuoteResponse(resp)
}

func (cl *JupClient) GetQuote(ctx context.Context, params *GetQuoteParams, reqEditors ...RequestEditorFunction) (*http.Response, error) {
	req, err := NewQuoteRequest(cl.Endpoint, params)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := cl.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}

	return cl.HTTPClient.Do(req)
}

func NewQuoteRequest(endpoint string, params *GetQuoteParams) (*http.Request, error) {
	var err error
	endpointUrl, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/quote")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}
	queryURL, err := endpointUrl.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	// Setting up parameters
	if params != nil {
		queryValues := queryURL.Query()

		if queryFrag, err := runtime.StyleParamWithLocation("form", true, "inputMint", runtime.ParamLocationQuery, params.InputMint); err != nil {
			return nil, err
		} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
			return nil, err
		} else {
			for key, val := range parsed {
				for _, v2 := range val {
					queryValues.Add(key, v2)
				}
			}
		}

		if queryFrag, err := runtime.StyleParamWithLocation("form", true, "outputMint", runtime.ParamLocationQuery, params.OutputMint); err != nil {
			return nil, err
		} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
			return nil, err
		} else {
			for key, val := range parsed {
				for _, v2 := range val {
					queryValues.Add(key, v2)
				}
			}
		}

		if queryFrag, err := runtime.StyleParamWithLocation("form", true, "amount", runtime.ParamLocationQuery, params.Amount); err != nil {
			return nil, err
		} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
			return nil, err
		} else {
			for key, val := range parsed {
				for _, v2 := range val {
					queryValues.Add(key, v2)
				}
			}
		}

		if params.SlippageBps != nil {
			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "slippageBps", runtime.ParamLocationQuery, *params.SlippageBps); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for key, val := range parsed {
					for _, v2 := range val {
						queryValues.Add(key, v2)
					}
				}
			}
		}
    
		if params.PlatformFeeBps != nil {
			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "platformFeeBps", runtime.ParamLocationQuery, *params.PlatformFeeBps); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for key, val := range parsed {
					for _, v2 := range val {
						queryValues.Add(key, v2)
					}
				}
			}
		}

		queryURL.RawQuery = queryValues.Encode()
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

type QuoteResponse struct {
	ComputedAutoSlippage *int32          `json:"computedAutoSlippage,omitempty"`
	ContextSlot          *float32        `json:"contextSlot,omitempty"`
	InAmount             string          `json:"inAmount"`
	InputMint            string          `json:"inputMint"`
	OtherAmountThreshold string          `json:"otherAmountThreshold"`
	OutAmount            string          `json:"outAmount"`
	OutputMint           string          `json:"outputMint"`
	PlatformFee          *PlatformFee    `json:"platformFee,omitempty"`
	PriceImpactPct       string          `json:"priceImpactPct"`
	RoutePlan            []RoutePlanStep `json:"routePlan"`
	SlippageBps          int32           `json:"slippageBps"`
	SwapMode             SwapMode        `json:"swapMode"`
	TimeTaken            *float32        `json:"timeTaken,omitempty"`
}

type GetQuoteResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *QuoteResponse
}

func ParseGetQuoteResponse(resp *http.Response) (*GetQuoteResponse, error) {
	bodyBytes, err := io.ReadAll(resp.Body)
	defer func() { resp.Body.Close() }()

	if err != nil {
		return nil, err
	}

	out := &GetQuoteResponse{
		Body:         bodyBytes,
		HTTPResponse: resp,
	}

	switch {
	case strings.Contains(resp.Header.Get("Content-Type"), "json") && resp.StatusCode == 200:
		var dest QuoteResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		out.JSON200 = &dest
	}

	return out, nil
}
