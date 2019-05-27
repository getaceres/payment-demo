package frontend

import "github.com/getaceres/payment-demo/payment"

type Response interface {
	GetLinks() map[string]string
}

// PaymentResponse is the response of a REST operation which returns a single payment
// swagger:model
type PaymentResponse struct {
	Data  payment.Payment   `json:"data"`
	Links map[string]string `json:"links"`
}

// PaymentListResponse is the response of a REST operation which returns a list of payments
// swagger:model
type PaymentListResponse struct {
	Data  []payment.Payment `json:"data"`
	Links map[string]string `json:"links"`
}

func (r PaymentResponse) GetLinks() map[string]string {
	return r.Links
}

func (r PaymentListResponse) GetLinks() map[string]string {
	return r.Links
}
