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

// PaymentRepository is the interface that any persistence backend must implement.
// It contains the basic CRUD operations for individual payments
// (AddPayment, GetPayment, UpdatePayment and DeletePayment)
// plus an operation that must return a list of payments filtered by arbitrary parameters.
type PaymentRepository interface {
	// AddPayment must save the payment information passed as parameter in the persistence backend asigning it a unique identifier.
	// It must return the saved document with this new identifier or an error if something unexpected happens
	AddPayment(pay payment.Payment) (payment.Payment, error)
	// UpdatePayment must replace the payment information whose identifier is in the input object with the information present in the input parameter.
	// It must return the updated payment information or an error if something goes wrong or the payment with such identifier does not exist.
	UpdatePayment(pay payment.Payment) (payment.Payment, error)
	// DeletePayment must delete the payment information whose identifier matches with the one passed as parameter.
	// It must return the deleted object or an error if something goes wrong or the payment with such identifier does not exist.
	DeletePayment(id string) (payment.Payment, error)
	// GetPayment must return the payment information whose identifier matches with the one passed as parameter.
	// It must return the payment information or an error if something goes wrong or the payment with such identifier does not exist.
	GetPayment(id string) (payment.Payment, error)
	// GetPayments must return a list of payments which match the filters passed as parameter.
	// If this parameter is nil or empty, it must return the whole list of payments available in the persistence backend.
	// In case of error, it must be returned as second parameter.
	GetPayments(filter map[string]string) ([]payment.Payment, error)
}

func (e NotFoundError) Error() string {
	return fmt.Sprintf("%s %s not found", e.ElementType, e.ID)
}

func (e AlreadyExistsError) Error() string {
	return fmt.Sprintf("%s %s already exists", e.ElementType, e.ID)
}
