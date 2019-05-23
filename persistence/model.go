package persistence

import (
	"fmt"

	"github.com/getaceres/payment-demo/payment"
)

const (
	PaymentElementType = "Payment"
)

type NotFoundError struct {
	ElementType string
	ID          string
}

type AlreadyExistsError struct {
	ElementType string
	ID          string
}

type PaymentRepository interface {
	AddPayment(pay payment.Payment) (payment.Payment, error)
	UpdatePayment(pay payment.Payment) (payment.Payment, error)
	DeletePayment(id string) (payment.Payment, error)
	GetPayment(id string) (payment.Payment, error)
	GetPayments(filter map[string]string) ([]payment.Payment, error)
}

func (e NotFoundError) Error() string {
	return fmt.Sprintf("%s %s not found", e.ElementType, e.ID)
}

func (e AlreadyExistsError) Error() string {
	return fmt.Sprintf("%s %s already exists", e.ElementType, e.ID)
}
