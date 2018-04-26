package main

func inputValidation(depOpts DeploymentOptions) bool {
	if depOpts.Application == "" {
		return false
	}

	if depOpts.Version == "" {
		return false
	}

	if depOpts.Environment == "" {
		return false
	}
	return true
}
