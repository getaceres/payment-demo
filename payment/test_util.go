package payment

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

const (
	TestResourcesPath   = "../test_resources"
	TestPaymentFileName = "test_payment.json"
)

func GetTestPayment(filename string) (Payment, error) {
	var testPayment Payment
	testPaymentFilePath := fmt.Sprintf("%s/%s", TestResourcesPath, filename)
	testPaymentFileContent, err := ioutil.ReadFile(testPaymentFilePath)
	if err != nil {
		return testPayment, err
	}

	err = json.Unmarshal([]byte(testPaymentFileContent), &testPayment)
	if err != nil {
		return testPayment, err
	}

	return testPayment, nil
}

func GetDefaultTestPayment() (Payment, error) {
	return GetTestPayment(TestPaymentFileName)
}
