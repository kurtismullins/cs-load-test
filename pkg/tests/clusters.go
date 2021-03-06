package tests

import (
	"bytes"
	"fmt"
	"net/http"
	"time"

	"github.com/nimrodshn/cs-load-test/pkg/helpers"
	"github.com/nimrodshn/cs-load-test/pkg/report"

	v1 "github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1"
	vegeta "github.com/tsenart/vegeta/v12/lib"
)

func TestListClusters(attacker *vegeta.Attacker,
	metrics vegeta.Metrics,
	rate vegeta.Pacer,
	outputDirectory string,
	duration time.Duration) error {

	fakeClusterProps := map[string]string{
		"fake_cluster": "true",
	}
	body, err := v1.NewCluster().
		Name("load-test").
		Properties(fakeClusterProps).
		MultiAZ(false).Build()
	if err != nil {
		return err
	}
	var raw bytes.Buffer
	err = v1.MarshalCluster(body, &raw)

	targeter := vegeta.NewStaticTargeter(vegeta.Target{
		Method: http.MethodGet,
		URL:    helpers.ClustersEndpoint,
		Body:   nil,
	})
	for res := range attacker.Attack(targeter, rate, duration, "Create") {
		metrics.Add(res)
	}
	metrics.Close()

	return report.Write("list-clusters-report",
		outputDirectory,
		&metrics)
}

func TestCreateCluster(attacker *vegeta.Attacker,
	metrics vegeta.Metrics,
	rate vegeta.Pacer,
	outputDirectory string,
	duration time.Duration) error {

	targeter := generateCreateClusterTargeter()
	for res := range attacker.Attack(targeter, rate, duration, "Create") {
		metrics.Add(res)
	}
	metrics.Close()

	return report.Write("create-cluster-report",
		outputDirectory,
		&metrics)
}

// Generates a targeter for the "POST /api/clusters_mgmt/v1/clusters" endpoint
// with monotonic increasing indexes.
// The clusters created are "fake clusters", that is, do not consume any cloud-provider infrastructure.
func generateCreateClusterTargeter() vegeta.Targeter {
	idx := 0

	targeter := func(t *vegeta.Target) error {
		fakeClusterProps := map[string]string{
			"fake_cluster": "true",
		}
		body, err := v1.NewCluster().
			Name(fmt.Sprintf("test-cluster-%d", idx)).
			Properties(fakeClusterProps).
			MultiAZ(false).Build()
		if err != nil {
			return err
		}

		var raw bytes.Buffer
		err = v1.MarshalCluster(body, &raw)
		if err != nil {
			return err
		}

		fmt.Printf("Body length: %d\n", raw.Len())

		t.Method = http.MethodPost
		t.URL = helpers.ClustersEndpoint
		t.Body = raw.Bytes()

		idx += 1
		return nil
	}
	return targeter
}
