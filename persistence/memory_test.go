package persistence

import "testing"

var tester = PaymentRepositoryTester{
	Repository: NewMemoryPaymentRepository(),
	ResourcesPath: "../test_resources",
}

func TestAdd(t *testing.T) {
	tester.TestAdd(t)
}

func TestUpdate(t *testing.T) {
	tester.TestUpdate(t)
}

func TestDelete(t *testing.T) {
	tester.TestDelete(t)
}

func TestGetId(t *testing.T) {
	tester.TestGetId(t)
}

func TestGetList(t *testing.T) {
	tester.TestGetList(t, 10)
}
