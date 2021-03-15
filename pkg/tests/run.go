package tests

import (
	"net/http"
	"time"

	"github.com/nimrodshn/cs-load-test/pkg/helpers"
	uuid "github.com/satori/go.uuid"
	vegeta "github.com/tsenart/vegeta/v12/lib"
)

type testCase func(attacker *vegeta.Attacker,
	testID string,
	metrics map[string]*vegeta.Metrics,
	rate vegeta.Pacer,
	outputDirectory string,
	duration time.Duration) error

func Run(
	attacker *vegeta.Attacker,
	metrics map[string]*vegeta.Metrics,
	rate vegeta.Pacer,
	outputDirectory string,
	duration time.Duration) error {

	// testId provides a common value to associate all output data from running
	// the full test suite with a single test run.
	testID := uuid.NewV4().String()

	// Specify Test Cases
	tests := []helpers.TestOptions{

		{
			TestName: "self-access-token",
			Path:     helpers.SelfAccessTokenEndpoint,
			Method:   http.MethodPost,
			Rate:     vegeta.Rate{Freq: helpers.SelfAccessTokenRate, Per: time.Second},
			Handler:  TestGenericEndpoint,
		},

		{
			TestName: "list-subscriptions",
			Path:     helpers.ListSubscriptionsEndpoint,
			Method:   http.MethodGet,
			Rate:     vegeta.Rate{Freq: helpers.ListSubscriptionsRate, Per: time.Second},
			Handler:  TestGenericEndpoint,
		},

		{
			TestName: "access-review",
			Path:     helpers.AccessReviewEndpoint,
			Method:   http.MethodPost,
			Body:     "{\"account_username\": \"rhn-support-tiwillia\", \"action\": \"get\", \"resource_type\": \"Subscription\"}",
			Rate:     vegeta.Rate{Freq: helpers.AccessReviewRate, Per: time.Second},
			Handler:  TestGenericEndpoint,
		},
	}

	// Include "Test Infrastructure" in each TestOptions.
	// This is done here to avoid duplicating a lot of code above.
	for i := range tests {
		tests[i].ID = testID
		tests[i].Duration = duration
		tests[i].OutputDirectory = outputDirectory
		tests[i].Attacker = attacker
		tests[i].Metrics = metrics
	}

	/*

		testCases := []testCase{
			TestCreateCluster,
			TestListClusters,
			TestSelfAccessToken,
			TestListSubscriptions,
			TestAccessReview,
			TestRegisterNewCluster,
		}

		for _, testCase := range testCases {
			err := testCase(attacker,
				testID,
				metrics,
				rate,
				outputDirectory,
				duration)
			if err != nil {
				return err
			}
		}

	*/

	for _, test := range tests {
		err := test.Handler(&test)
		if err != nil {
			return err
		}
	}

	return nil
}
