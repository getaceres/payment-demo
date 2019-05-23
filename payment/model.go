package payment

type Payment struct {
	Type           string                `json:"type,omitempty"`
	ID             string                `json:"id,omitempty"`
	Version        int                   `json:"version,omitempty"`
	OrganisationID string                `json:"organisation_id,omitempty"`
	Attributes     PaymentAttributesType `json:"attributes,omitempty"`
}

type PaymentAttributesType struct {
	Amount               string                         `json:"amount,omitempty"`
	BeneficiaryParty     PaymentPartyType               `json:"beneficiary_party,omitempty"`
	ChargesInformation   PaymentChargesInformationType  `json:"charges_information,omitempty"`
	Currency             string                         `json:"currency,omitempty"`
	DebtorParty          PaymentPartyType               `json:"debtor_party,omitempty"`
	EndToEndReference    string                         `json:"end_to_end_reference,omitempty"`
	FX                   PaymentExchangeInformationType `json:"fx,omitempty"`
	NumericReference     string                         `json:"numeric_reference,omitempty"`
	PaymentID            string                         `json:"payment_id,omitempty"`
	PaymentPurpose       string                         `json:"payment_purpose,omitempty"`
	PaymentScheme        string                         `json:"payment_scheme,omitempty"`
	PaymentType          string                         `json:"payment_type,omitempty"`
	ProcessingDate       string                         `json:"processing_date,omitempty"`
	Reference            string                         `json:"reference,omitempty"`
	SchemePaymentSubType string                         `json:"scheme_payment_sub_type,omitempty"`
	SchemePaymentType    string                         `json:"scheme_payment_type,omitempty"`
	SponsorParty         PaymentPartyType               `json:"sponsor_party,omitempty"`
}

type PaymentPartyType struct {
	AccountName       string `json:"account_name,omitempty"`
	AccountNumber     string `json:"account_number,omitempty"`
	AccountNumberCode string `json:"account_number_code,omitempty"`
	AccountType       *int   `json:"account_type,omitempty"`
	Address           string `json:"address,omitempty"`
	BankID            string `json:"bank_id,omitempty"`
	BankIDCode        string `json:"bank_id_code,omitempty"`
	Name              string `json:"name,omitempty"`
}

type PaymentChargesInformationType struct {
	BearerCode              string              `json:"bearer_code,omitempty"`
	SenderCharges           []PaymentAmountType `json:"sender_charges,omitempty"`
	ReceiverChargesAmount   string              `json:"receiver_charges_amount,omitempty"`
	ReceiverChargesCurrency string              `json:"receiver_charges_currency,omitempty"`
}

type PaymentAmountType struct {
	Amount   string `json:"amount,omitempty"`
	Currency string `json:"currency,omitempty"`
}

type PaymentExchangeInformationType struct {
	ContractReference string `json:"contract_reference,omitempty"`
	ExchangeRate      string `json:"exchange_rate,omitempty"`
	OriginalAmount    string `json:"original_amount,omitempty"`
	OriginalCurrency  string `json:"original_currency,omitempty"`
}
