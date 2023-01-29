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
	// Role is the IAM role to use when invoking a deployment.
	Role string `json:"Role"`
	// MaxAttempts is the Number of attempts to wait for service to become stable, with subsequent 15 sec pause.
	MaxAttempts int `json:"MaxAttempts"`
	// RefreshSecrets will update all container definition secrets to include ssm parameters that exists with the prefix "/<cluster>/service/*"
	RefreshSecrets bool `json:"RefreshSecrets"`
	// The ssm parameter store prefix to pull secrets from. Default: "/<cluster>/service/*"
	SecretsPrefix string `json:"SecretsPrefix"`
	// DryRun will preview changes
	DryRun bool `json:"DryRun"`
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
