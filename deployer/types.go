package deployer

// DeploymentOptions set the desired state of your deployment
type DeploymentOptions struct {
	// Application name as it exists in ECS
	Application string `json:"Application"`
	// Desired version of ECS Application
	Version string `json:"Version"`
	// Environment (ECS Cluster Name) you would like to deploy to
	Environment string `json:"Environment"`
	// Description is an optional parameter adding context to the change
	Description string `json:"Description"`
	// DisableMemoryReservation determines whether or not to automatically adjust the memory reservation
	DisableMemoryReservation bool `json:"DisableMemoryReservation"`
}

// DeploymentResults maintain the depyments latest results
type DeploymentResults struct {
	// SuccessfullyInvoked bool value depicting a successful deployment invocation
	SuccessfullyInvoked bool   `json:"SuccessfullyInvoked"`
	ClusterArn          string `json:"ClusterArn"`
	ServiceArn          string `json:"ServiceArn"`
	ServiceName         string `json:"ServiceName"`
	TaskDefinition      string `json:"TaskDefinition"`
}
