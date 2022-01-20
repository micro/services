package domain

type Status string

const (
	StatusDeploying       = "Deploying"
	StatusUpdating        = "Updating"
	StatusDeleting        = "Deleting"
	StatusDeploymentError = "DeploymentError"
)
