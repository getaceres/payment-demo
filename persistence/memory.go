package persistence

import (
	"errors"

	"github.com/getaceres/payment-demo/payment"
	"github.com/google/uuid"
)

type MemoryPaymentRepository struct {
	Payments map[string]payment.Payment
}

func NewMemoryPaymentRepository() *MemoryPaymentRepository {
	return &MemoryPaymentRepository{
		Payments: make(map[string]payment.Payment),
	}
}

func (m *MemoryPaymentRepository) AddPayment(pay payment.Payment) (payment.Payment, error) {
	pay.ID = uuid.New().String()
	m.Payments[pay.ID] = pay
	return pay, nil
}

func (m *MemoryPaymentRepository) UpdatePayment(pay payment.Payment) (payment.Payment, error) {
	id := pay.ID
	if id == "" {
		return pay, errors.New("Payment with empty identifier passed")
	}

	_, ok := m.Payments[id]
	if !ok {
		return pay, NotFoundError{PaymentElementType, id}
	}
	m.Payments[id] = pay
	return pay, nil
}

func (m *MemoryPaymentRepository) DeletePayment(id string) (payment.Payment, error) {
	pay, ok := m.Payments[id]
	if !ok {
		return pay, NotFoundError{PaymentElementType, id}
	}
	delete(m.Payments, id)
	return pay, nil
}

func (m *MemoryPaymentRepository) GetPayment(id string) (payment.Payment, error) {
	pay, ok := m.Payments[id]
	if !ok {
		return pay, NotFoundError{PaymentElementType, id}
	}
	return pay, nil
}

func (m *MemoryPaymentRepository) GetPayments(filter map[string]string) ([]payment.Payment, error) {
	result := make([]payment.Payment, 0, len(m.Payments))
	for _, pay := range m.Payments {
		result = append(result, pay)
	}
	return result, nil
}
