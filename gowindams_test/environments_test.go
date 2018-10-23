package gowindams_test

import (
	"github.com/Inspectools/gowindams"
	"testing"
)

const TEST_FILE = "environments_test.yaml"
var TEST_DATA = [1]gowindams.Environment{
	{
		Name: "InspecTools Dev",
		ServiceAppId: "InspecTools-Dev-Services-App-Id",
		ServiceURI: "https://servicesdev.inspectools.net/",
	},
}

func TestLoadFromFile(testing *testing.T) {
	environments, err := gowindams.LoadEnvironments(TEST_FILE)
	if err != nil {
		testing.Fatalf("Unable to open test file: %s\n", err)
	}
	ecount := len(TEST_DATA)
	rcount := len(*environments)
	if rcount != ecount {
		testing.Fatalf("Loaded %d environments, %d were expected\n", rcount, ecount)
	}

	for i := 0; i < rcount; i++ {
		compareEnvironments(testing, TEST_DATA[i], (*environments)[i])
	}
}

func compareEnvironments(testing *testing.T, expected gowindams.Environment, got gowindams.Environment) {
	compareStrings(testing, expected.Name, got.Name)
	compareStrings(testing, expected.ServiceAppId, got.ServiceAppId)
	compareStrings(testing, expected.ServiceURI, got.ServiceURI)
}

func compareStrings(testing *testing.T, expected string, got string) {
	if got != expected {
		testing.Fatalf("Expected %s but got %s\n", expected, got)
	}
}
