package mongo

import (
	"flag"
	"fmt"
	"os"
	"testing"

	"github.com/getaceres/payment-demo/persistence"
)

var integrationMongo = flag.Bool("mongo", false, "run MongoDB tests")

var tester = persistence.PaymentRepositoryTester{
	ResourcesPath: "../../test_resources",
}

func TestMain(m *testing.M) {
	repo, err := NewMongoPaymentRepository("mongodb://localhost:27017", "payment-demo-test")
	if err != nil {
		fmt.Printf("Error getting MongoDB repository: %s", err.Error())
	}
	tester.Repository = repo
	os.Exit(m.Run())
}

func TestAdd(t *testing.T) {
	if *integrationMongo {
		tester.TestAdd(t)
	}
}

func TestUpdate(t *testing.T) {
	if *integrationMongo {
		tester.TestUpdate(t)
	}
}

func TestDelete(t *testing.T) {
	if *integrationMongo {
		tester.TestDelete(t)
	}
}

func TestGetId(t *testing.T) {
	if *integrationMongo {
		tester.TestGetId(t)
	}
}

func TestGetList(t *testing.T) {
	if *integrationMongo {
		tester.TestGetList(t, 10)
	}
}
