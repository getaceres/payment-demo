package frontend

import "github.com/getaceres/payment-demo/payment"

type Response interface {
	GetLinks() map[string]string
}

type PaymentResponse struct {
	Data  payment.Payment   `json:"data"`
	Links map[string]string `json:"links"`
}

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
