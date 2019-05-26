package persistence

import (
	"testing"

	"github.com/getaceres/payment-demo/payment"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
)

const (
	testPaymentFileName = "test_payment.json"
)

type PaymentRepositoryTester struct {
	Repository    PaymentRepository
	ResourcesPath string
}

func (p PaymentRepositoryTester) getDefaultPayment(t *testing.T) payment.Payment {
	payment, err := payment.GetDefaultTestPayment(p.ResourcesPath)
	if err != nil {
		t.Fatalf("Error getting default payment: %s", err.Error())
	}
	return payment
}

func (p PaymentRepositoryTester) TestAdd(t *testing.T) {
	payment := p.getDefaultPayment(t)

	newPayment, err := p.Repository.AddPayment(payment)
	if err != nil {
		t.Fatalf("Error adding payment: %s", err.Error())
	}

	if newPayment.ID == payment.ID {
		t.Fatalf("The payment was saved with the default ID but a new one was expected")
	}
}

func (p PaymentRepositoryTester) TestUpdate(t *testing.T) {
	payment := p.getDefaultPayment(t)

	newPayment, err := p.Repository.AddPayment(payment)
	if err != nil {
		t.Fatalf("Error adding payment: %s", err.Error())
	}

	newPayment.Version = 2

	updated, err := p.Repository.UpdatePayment(newPayment)
	if err != nil {
		t.Fatalf("Error updating payment: %s", err.Error())
	}

	if !cmp.Equal(updated, newPayment) {
		t.Fatalf("Returned updated payment differs from the passed one.\nReturned:\n%v\nBut expected:\n%v", updated, newPayment)
	}

	newId := uuid.New().String()
	updated.ID = newId
	updated, err = p.Repository.UpdatePayment(updated)
	p.checkNotFoundError(newId, "updating", err, t)
}

func (p PaymentRepositoryTester) TestDelete(t *testing.T) {
	payment := p.getDefaultPayment(t)

	newPayment, err := p.Repository.AddPayment(payment)
	if err != nil {
		t.Fatalf("Error adding payment: %s", err.Error())
	}

	deleted, err := p.Repository.DeletePayment(newPayment.ID)
	if err != nil {
		t.Fatalf("Error adding payment: %s", err.Error())
	}

	if !cmp.Equal(deleted, newPayment) {
		t.Fatalf("Returned deleted payment differs from the passed one.\nReturned:\n%v\nBut expected:\n%v", deleted, newPayment)
	}

	id := deleted.ID
	deleted, err = p.Repository.GetPayment(id)
	p.checkNotFoundError(id, "getting", err, t)
}

func (p PaymentRepositoryTester) TestGetId(t *testing.T) {
	payment := p.getDefaultPayment(t)

	newPayment, err := p.Repository.AddPayment(payment)
	if err != nil {
		t.Fatalf("Error adding payment: %s", err.Error())
	}

	got, err := p.Repository.GetPayment(newPayment.ID)
	if !cmp.Equal(got, newPayment) {
		t.Fatalf("Returned get payment differs from the created one.\nReturned:\n%v\nBut expected:\n%v", got, newPayment)
	}

	newID := uuid.New().String()
	got, err = p.Repository.GetPayment(newID)
	p.checkNotFoundError(newID, "getting", err, t)
}

func (p PaymentRepositoryTester) TestGetList(t *testing.T, toCreate int) {
	existing, err := p.Repository.GetPayments(nil)
	if err != nil {
		t.Errorf("Error listing all payments: %s", err.Error())
	}
	numInitial := len(existing)
	payment := p.getDefaultPayment(t)

	for i := 0; i < toCreate; i++ {
		_, err := p.Repository.AddPayment(payment)
		if err != nil {
			t.Fatalf("Error adding payment: %s", err.Error())
		}
	}

	existing, err = p.Repository.GetPayments(nil)
	if err != nil {
		t.Errorf("Error listing all payments after adding: %s", err.Error())
	}

	if len(existing) != numInitial+toCreate {
		t.Errorf("Expected %d elements after adding but got %d", numInitial+toCreate, len(existing))
	}

}

func (p PaymentRepositoryTester) checkNotFoundError(id, action string, err error, t *testing.T) {
	if err == nil {
		t.Fatalf("Expected NotFound error %s non existing payment %s but got nil", action, id)
	} else {
		_, ok := err.(NotFoundError)
		if !ok {
			t.Fatalf("Expected NotFound error %s non existing payment %s but got error %v", action, id, err)
		}
	}
}
