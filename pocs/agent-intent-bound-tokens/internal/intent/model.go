package intent

type Claims struct {
	Issuer         string `json:"iss"`
	Subject        string `json:"sub"`
	Audience       string `json:"aud"`
	Action         string `json:"action"`
	Resource       string `json:"resource"`
	MaxAmountCents int64  `json:"max_amount_cents"`
	NotBefore      int64  `json:"nbf"`
	ExpiresAt      int64  `json:"exp"`
	ID             string `json:"jti"`
}

type Request struct {
	Token       string `json:"token"`
	Audience    string `json:"audience"`
	Action      string `json:"action"`
	Resource    string `json:"resource"`
	AmountCents int64  `json:"amount_cents"`
}

type Decision struct {
	Allowed bool    `json:"allowed"`
	Code    string  `json:"code"`
	Reason  string  `json:"reason"`
	Claims  *Claims `json:"claims,omitempty"`
}
