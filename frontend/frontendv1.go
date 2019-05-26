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

func (a *FrontendV1) DeletePayment(w http.ResponseWriter, r *http.Request) {
	a.doPaymentOperation(w, r, a.PaymentRepository.DeletePayment, "deleting")
}

func (a *FrontendV1) GetPayment(w http.ResponseWriter, r *http.Request) {
	a.doPaymentOperation(w, r, a.PaymentRepository.GetPayment, "getting")
}

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
