package db

import (
	"github.com/superfly/fly-go"
	"gorm.io/datatypes"
)

type DeploymentMetadata struct {
	LockID string `json:"lock_id"`
	Token  string `json:"token"`
}

type AppStateSnapshot struct {
	Machines []*fly.Machine `json:"machines"`
}

type Deployment struct {
	ID        string                                 `json:"id"`
	AppID     string                                 `json:"app_id"`
	Status    string                                 `json:"status"`
	Metadata  datatypes.JSONType[DeploymentMetadata] `json:"metadata"`
	State     string                                 `json:"state"`
	Snapshot  datatypes.JSONType[AppStateSnapshot]   `json:"snapshot"`
	WorkerID  string                                 `json:"worker_id"`
	CreatedAt string                                 `json:"created_at"`
	UpdatedAt string                                 `json:"updated_at"`
}

type Event struct {
	ID           string `json:"id"`
	DeploymentID string `json:"deployment_id"`
	EntityID     string `json:"entity_id"`
	Message      string `json:"message"`
	Action       string `json:"action"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}
