package frontend

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/getaceres/payment-demo/payment"
	"github.com/getaceres/payment-demo/persistence"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

const numPayments = 10

var frontend = FrontendV1{
	Router:            mux.NewRouter(),
	PaymentRepository: persistence.NewMemoryPaymentRepository(),
}

func TestMain(m *testing.M) {
	frontend.InitializeRoutes()
	os.Exit(m.Run())
}

func getDefaultPayment(t *testing.T) payment.Payment {
	payment, err := payment.GetDefaultTestPayment()
	if err != nil {
		t.Fatalf("Error getting default payment: %s", err.Error())
	}
	return payment
}

func addPayment(t *testing.T) payment.Payment {
	pay := getDefaultPayment(t)
	pay, err := frontend.PaymentRepository.AddPayment(pay)
	if err != nil {
		t.Fatalf("Error creating payment: %s", err.Error())
	}
	return pay
}

func executeRequest(t *testing.T, operation, path string, payload interface{}) *httptest.ResponseRecorder {
	var reader io.Reader
	if payload != nil {
		text, err := json.Marshal(payload)
		if err != nil {
			t.Fatalf("Error marshaling payload: %s", err.Error())
		}
		reader = bytes.NewBuffer(text)
	}
	req, err := http.NewRequest(operation, path, reader)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	frontend.Router.ServeHTTP(rr, req)
	return rr
}

func checkResponseCode(t *testing.T, result *httptest.ResponseRecorder, expectedCode int) {
	if result.Code != expectedCode {
		t.Fatalf("Unexpected status code returned: Expected %d but got %d", http.StatusCreated, result.Code)
	}
}

func checkResponse(t *testing.T, result *httptest.ResponseRecorder, expectedCode int, bodyOut Response) {
	checkResponseCode(t, result, expectedCode)
	err := ReadBody(result.Body, bodyOut)
	if err != nil {
		t.Fatalf("Error deserializing response body: %s", err.Error())
	}
	links := bodyOut.GetLinks()
	if links == nil || len(links) == 0 {
		t.Fatal("Links section is empty")
	}

	self, ok := links["self"]
	if !ok {
		t.Fatal("Self link not found")
	}

	if self == "" {
		t.Fatal("Self link is empty")
	}
}

func checkPaymentResponse(t *testing.T, result *httptest.ResponseRecorder, expectedCode int) payment.Payment {
	var returned PaymentResponse
	checkResponse(t, result, expectedCode, &returned)
	return returned.Data
}

func checkPaymentListResponse(t *testing.T, result *httptest.ResponseRecorder, expectedCode int) []payment.Payment {
	var returned PaymentListResponse
	checkResponse(t, result, expectedCode, &returned)
	return returned.Data
}

func TestAdd(t *testing.T) {
	pay := getDefaultPayment(t)
	result := executeRequest(t, "POST", "/v1/payments", pay)
	returned := checkPaymentResponse(t, result, http.StatusCreated)

	if pay.ID == returned.ID {
		t.Fatal("Input and output identifiers are the same")
	}

	pay.ID = returned.ID
	if !cmp.Equal(pay, returned) {
		t.Fatalf("Created payment differs from expected.\nExpected:\n%v\nBut got:\n%v", pay, returned)
	}
}

func TestUpdate(t *testing.T) {
	pay := addPayment(t)

	pay.Version = 2
	result := executeRequest(t, "PUT", "/v1/payments", pay)
	returned := checkPaymentResponse(t, result, http.StatusOK)

	if pay.ID != returned.ID {
		t.Fatal("Input and output identifiers are not the same")
	}

	if !cmp.Equal(pay, returned) {
		t.Fatalf("Updated payment differs from expected.\nExpected:\n%v\nBut got:\n%v", pay, returned)
	}

	id := uuid.New().String()
	pay.ID = id
	result = executeRequest(t, "PUT", "/v1/payments", pay)
	checkResponseCode(t, result, http.StatusNotFound)
}

func TestDelete(t *testing.T) {
	pay := addPayment(t)

	result := executeRequest(t, "DELETE", fmt.Sprintf("/v1/payments/%s", pay.ID), nil)
	returned := checkPaymentResponse(t, result, http.StatusOK)

	if !cmp.Equal(pay, returned) {
		t.Fatalf("Deleted payment differs from expected.\nExpected:\n%v\nBut got:\n%v", pay, returned)
	}

	result = executeRequest(t, "DELETE", fmt.Sprintf("/v1/payments/%s", uuid.New().String()), pay)
	checkResponseCode(t, result, http.StatusNotFound)
}

func TestGet(t *testing.T) {
	pay := addPayment(t)

	result := executeRequest(t, "GET", fmt.Sprintf("/v1/payments/%s", pay.ID), nil)
	returned := checkPaymentResponse(t, result, http.StatusOK)

	if !cmp.Equal(pay, returned) {
		t.Fatalf("Deleted payment differs from expected.\nExpected:\n%v\nBut got:\n%v", pay, returned)
	}

	result = executeRequest(t, "GET", fmt.Sprintf("/v1/payments/%s", uuid.New().String()), pay)
	checkResponseCode(t, result, http.StatusNotFound)
}

func TestGetList(t *testing.T) {
	initial, err := frontend.PaymentRepository.GetPayments(nil)
	if err != nil {
		t.Fatalf("Error getting the initial list of payments: %s", err.Error())
	}

	numInitial := len(initial)

	for i := 0; i < numPayments; i++ {
		addPayment(t)
	}

	result := executeRequest(t, "GET", "/v1/payments", nil)
	returned := checkPaymentListResponse(t, result, http.StatusOK)
	expected := numInitial + numPayments
	if len(returned) != expected {
		t.Fatalf("Unexpected number of payments returned. Expected %d but got %d", expected, len(returned))
	}
}
