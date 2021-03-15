package helpers

// Default rate (requests/second) for each endpoint
const (
	CreateClusterRate           = 10
	ListClustersRate            = 10
	SelfAccessTokenRate         = 17 // ~1000/hour
	ListSubscriptionsRate       = 34 // ~2000/hour
	AccessReviewRate            = 100
	RegisterNewClusterRate      = 17 // ~1000/hour
	RegisterExistingClusterRate = 25 // "reauth"
)
