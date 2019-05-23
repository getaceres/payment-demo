package payment

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/google/go-cmp/cmp"
)

type PaymentListType struct {
	Data []Payment `json:"data"`
}

func TestPaymensUnmarshal(t *testing.T) {
	var payments PaymentListType
	fileName := "../test_resources/payment_list.json"
	fileContent, err := ioutil.ReadFile(fileName)
	if err != nil {
		t.Fatalf("Error reading test payment list: %s", err.Error())
	}
	err = json.Unmarshal(fileContent, &payments)
	if err != nil {
		t.Fatalf("Error unmarshaling payment list: %s", err.Error())
	}

	if payments.Data == nil {
		t.Fatalf("Unmarshaled a nil list of payments from %s", fileName)
	}

	if len(payments.Data) != 14 {
		t.Fatalf("Wrong number of payments unmarshaled. Expected 14 but got %d", len(payments.Data))
	}
}

func TestPaymentSerialization(t *testing.T) {
	testPayment, err := GetDefaultTestPayment()
	if err != nil {
		t.Fatalf("Error getting test payment: %s", err.Error())
	}

	testAccountType := 0

	comparePayment := Payment{
		Type:           "Payment",
		ID:             "4ee3a8d8-ca7b-4290-a52c-dd5b6165ec43",
		Version:        0,
		OrganisationID: "743d5b63-8e6f-432e-a8fa-c5d8d2ee5fcb",
		Attributes: PaymentAttributesType{
			Amount: "100.21",
			BeneficiaryParty: PaymentPartyType{
				AccountName:       "W Owens",
				AccountNumber:     "31926819",
				AccountNumberCode: "BBAN",
				AccountType:       &testAccountType,
				Address:           "1 The Beneficiary Localtown SE2",
				BankID:            "403000",
				BankIDCode:        "GBDSC",
				Name:              "Wilfred Jeremiah Owens",
			},
			ChargesInformation: PaymentChargesInformationType{
				BearerCode: "SHAR",
				SenderCharges: []PaymentAmountType{
					PaymentAmountType{
						Amount:   "5.00",
						Currency: "GBP",
					},
					PaymentAmountType{
						Amount:   "10.00",
						Currency: "USD",
					},
				},
				ReceiverChargesAmount:   "1.00",
				ReceiverChargesCurrency: "USD",
			},
			Currency: "GBP",
			DebtorParty: PaymentPartyType{
				AccountName:       "EJ Brown Black",
				AccountNumber:     "GB29XABC10161234567801",
				AccountNumberCode: "IBAN",
				Address:           "10 Debtor Crescent Sourcetown NE1",
				BankID:            "203301",
				BankIDCode:        "GBDSC",
				Name:              "Emelia Jane Brown",
			},
			EndToEndReference: "Wil piano Jan",
			FX: PaymentExchangeInformationType{
				ContractReference: "FX123",
				ExchangeRate:      "2.00000",
				OriginalAmount:    "200.42",
				OriginalCurrency:  "USD",
			},
			NumericReference:     "1002001",
			PaymentID:            "123456789012345678",
			PaymentPurpose:       "Paying for goods/services",
			PaymentScheme:        "FPS",
			PaymentType:          "Credit",
			ProcessingDate:       "2017-01-18",
			Reference:            "Payment for Em's piano lessons",
			SchemePaymentSubType: "InternetBanking",
			SchemePaymentType:    "ImmediatePayment",
			SponsorParty: PaymentPartyType{
				AccountNumber: "56781234",
				BankID:        "123123",
				BankIDCode:    "GBDSC",
			},
		},
	}

	if !cmp.Equal(testPayment, comparePayment) {
		t.Fatalf("Unmarshaled payment differs from expected.\nExpected:\n%v\nBut got:\n%v", comparePayment, testPayment)
	}

}
