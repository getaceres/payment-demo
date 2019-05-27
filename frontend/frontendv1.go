package frontend

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/getaceres/payment-demo/payment"
	"github.com/getaceres/payment-demo/persistence"
	"github.com/gorilla/mux"
	"github.com/imdario/mergo"
)

const (
	basePath = "/v1"
)

type FrontendV1 struct {
	Router            *mux.Router
	PaymentRepository persistence.PaymentRepository
}

func (a *FrontendV1) InitializeRoutes() {
	a.Router.HandleFunc(basePath+"/payments", a.AddPayment).Methods("POST")
	a.Router.HandleFunc(basePath+"/payments", a.GetPaymentList).Methods("GET")
	a.Router.HandleFunc(basePath+"/payments/{paymentID}", a.UpdatePayment).Methods("PUT")
	a.Router.HandleFunc(basePath+"/payments/{paymentID}", a.DeletePayment).Methods("DELETE")
	a.Router.HandleFunc(basePath+"/payments/{paymentID}", a.GetPayment).Methods("GET")
}

func (a *FrontendV1) doPaymentOperation(w http.ResponseWriter, r *http.Request, function func(id string) (payment.Payment, error), verb string) {
	paymentID, ok := mux.Vars(r)["paymentID"]
	if !ok {
		RespondWithError(w, http.StatusBadRequest, errors.New("Payment identifier is mandatory for this operation"))
		return
	}

	payment, err := function(paymentID)
	if err != nil {
		RespondWithError(w, GetPersistenceErrorCode(err), fmt.Errorf("Error %s payment %s: %s", verb, paymentID, err.Error()))
		return
	}

	RespondWithJSON(w, http.StatusOK, PaymentResponse{
		Data: payment,
		Links: map[string]string{
			"self": r.URL.String(),
		},
	})
	return
}

// AddPayment saves a new payment information into the database
// swagger:operation POST /payments addPayment
//
// Adds the payment object passed in the operation body to the database
//
// ---
// produces:
// - application/json
// - application/text
// parameters:
// - name: payment
//   in: body
//   description: The payment to add
//   required: true
//   schema:
//     "$ref": "#/definitions/Payment"
// responses:
//   '201':
//     description: The payment with an assigned identifier
//     schema:
//       "$ref": "#/definitions/PaymentResponse"
//   500:
//     description: Unexpected error
func (a *FrontendV1) AddPayment(w http.ResponseWriter, r *http.Request) {
	var pay payment.Payment
	err := ReadBody(r.Body, &pay)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, fmt.Errorf("Error reading payment body: %s", err.Error()))
		return
	}

	updated, err := a.PaymentRepository.AddPayment(pay)
	if err != nil {
		RespondWithError(w, GetPersistenceErrorCode(err), fmt.Errorf("Error saving payment: %s", err.Error()))
		return
	}

	RespondWithJSON(w, http.StatusCreated, PaymentResponse{
		Data: updated,
		Links: map[string]string{
			"self": fmt.Sprintf("%s/%s", r.URL.String(), updated.ID),
		},
	})
	return
}

// UpdatePayment updates the information of a payment by providing a partial document that will be merged with the original one
// swagger:operation PUT /payments/{paymentID} updatePayment
//
// Receives a partial Payment document with the fields that should be updated. It will be merged with the information already present.
//
// ---
// produces:
// - application/json
// - application/text
// parameters:
// - name: paymentID
//   in: path
//   description: The identifier of the payment to update
//   required: true
//   type: string
// - name: payment
//   in: body
//   description: A partial payment document with the fields to update
//   required: true
//   schema:
//     "$ref": "#/definitions/Payment"
// responses:
//   '200':
//     description: The updated payment
//     schema:
//       "$ref": "#/definitions/PaymentResponse"
//   500:
//     description: Unexpected error
//   404:
//     description: Payment not found
func (a *FrontendV1) UpdatePayment(w http.ResponseWriter, r *http.Request) {
	a.doPaymentOperation(w, r, func(id string) (payment.Payment, error) {
		existing, err := a.PaymentRepository.GetPayment(id)
		if err != nil {
			return existing, err
		}

		var partial payment.Payment
		err = ReadBody(r.Body, &partial)
		if err != nil {
			return existing, fmt.Errorf("Error reading partial update: %s", err.Error())
		}

		err = mergo.Merge(&partial, existing)
		if err != nil {
			return existing, fmt.Errorf("Error applying new values to existing payment %s: %s", id, err.Error())
		}

		partial.ID = id
		return a.PaymentRepository.UpdatePayment(partial)
	}, "updating")
}

// DeletePayment deletes the information of a payment given its identifier
// swagger:operation DELETE /payments/{paymentID} deletePayment
//
// Receives the identifier of the payment to delete and returns the deleted document.
//
// ---
// produces:
// - application/json
// - application/text
// parameters:
// - name: paymentID
//   in: path
//   description: The identifier of the payment to update
//   required: true
//   type: string
// responses:
//   '200':
//     description: The deleted payment
//     schema:
//       "$ref": "#/definitions/PaymentResponse"
//   500:
//     description: Unexpected error
//   404:
//     description: Payment not found
func (a *FrontendV1) DeletePayment(w http.ResponseWriter, r *http.Request) {
	a.doPaymentOperation(w, r, a.PaymentRepository.DeletePayment, "deleting")
}

// GetPayment retrieves the information of a payment given its identifier
// swagger:operation GET /payments/{paymentID} getPayment
//
// Receives the identifier of the payment to delete and returns the information of the document document.
//
// ---
// produces:
// - application/json
// - application/text
// parameters:
// - name: paymentID
//   in: path
//   description: The identifier of the requiered payment
//   required: true
//   type: string
// responses:
//   '200':
//     description: The required payment
//     schema:
//       "$ref": "#/definitions/PaymentResponse"
//   500:
//     description: Unexpected error
//   404:
//     description: Payment not found
func (a *FrontendV1) GetPayment(w http.ResponseWriter, r *http.Request) {
	a.doPaymentOperation(w, r, a.PaymentRepository.GetPayment, "getting")
}

// GetPaymentList retrieves a list with all the registered payments
// swagger:operation GET /payments getPaymentList
//
// This implementation doesn't receive any extra parameter but in future implementations the list could be filtered with query parameters.
//
// ---
// produces:
// - application/json
// - application/text
// responses:
//   '200':
//     description: The list of registered payments
//     schema:
//       type: array
//       items:
//         "$ref": "#/definitions/PaymentListResponse"
//   500:
//     description: Unexpected error
func (a *FrontendV1) GetPaymentList(w http.ResponseWriter, r *http.Request) {
	payments, err := a.PaymentRepository.GetPayments(nil)
	if err != nil {
		RespondWithError(w, GetPersistenceErrorCode(err), fmt.Errorf("Error getting payment list: %s", err.Error()))
		return
	}

	RespondWithJSON(w, http.StatusOK, PaymentListResponse{
		Data: payments,
		Links: map[string]string{
			"self": r.URL.String(),
		},
	})
}
