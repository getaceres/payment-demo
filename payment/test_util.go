package payment

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

const (
	TestPaymentFileName = "test_payment.json"
)

func GetTestPayment(resourcesPath, filename string) (Payment, error) {
	var testPayment Payment
	testPaymentFilePath := fmt.Sprintf("%s/%s", resourcesPath, filename)
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

func GetDefaultTestPayment(resourcesPath string) (Payment, error) {
	return GetTestPayment(resourcesPath, TestPaymentFileName)
}
