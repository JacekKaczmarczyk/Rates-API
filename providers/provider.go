package providers

type Response struct {
	AsOF     string      `json:"asOf"`
	Provider string      `json:"provider"`
	Rates    []RateValue `json:"rates"`
}

type RateValue struct {
	Code  string  `json:"code"`
	Value float64 `json:"value"`
}

type ErrorResponse struct {
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

type Provider interface {
	GetCurrencies(codes []string, date string) (Response, error, int)
}
