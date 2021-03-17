package tests

import (
	"bytes"
	"fmt"
	"log"
	"net/http"

	"github.com/nimrodshn/cs-load-test/pkg/helpers"
	"github.com/nimrodshn/cs-load-test/pkg/result"
	amsv1 "github.com/openshift-online/ocm-sdk-go/accountsmgmt/v1"
	v1 "github.com/openshift-online/ocm-sdk-go/accountsmgmt/v1"
	uuid "github.com/satori/go.uuid"
	vegeta "github.com/tsenart/vegeta/v12/lib"
)

func TestRegisterNewCluster(options *helpers.TestOptions) error {

	testName := options.TestName
	log.Printf("Executing Test: %s", testName)

	// Fetch the authorization token and create a dynamic Target generator for
	// building valid HTTP Requests
	authorizationToken := "TODO: Fetch me programattically"
	targeter := generateClusterRegistrationTargeter(authorizationToken)

	// Create a file to store results
	fileName := fmt.Sprintf("%s_%s.json", options.ID, testName)
	resultFile, err := createFile(fileName, options.OutputDirectory)
	defer resultFile.Close()
	if err != nil {
		return err
	}

	// Store Metrics from load test
	options.Metrics[testName] = new(vegeta.Metrics)
	defer options.Metrics[testName].Close()

	for res := range options.Attacker.Attack(targeter, options.Rate, options.Duration, testName) {
		result.Write(res, resultFile)
		options.Metrics[testName].Add(res)
	}

	fmt.Printf("Results written to: %s/%s\n", options.OutputDirectory, fileName)

	return nil
}

func getAuthorizationToken() {
	// TODO: Implement me
}

func generateClusterRegistrationTargeter(authorizationToken string) vegeta.Targeter {

	targeter := func(t *vegeta.Target) error {

		// Each Cluster uses a UUID to ensure uniqueness
		clusterId := uuid.NewV4().String()
		body, err := amsv1.NewClusterRegistrationRequest().AuthorizationToken(authorizationToken).ClusterID(clusterId).Build()
		if err != nil {
			return err
		}

		var raw bytes.Buffer
		err = v1.MarshalClusterRegistrationRequest(body, &raw)
		if err != nil {
			return err
		}

		t.Method = http.MethodPost
		t.URL = helpers.ClusterRegistrationEndpoint
		t.Body = raw.Bytes()

		return nil
	}

	return targeter
}
