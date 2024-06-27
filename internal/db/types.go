package db

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"github.com/superfly/fly-go"
)

type DeploymentMetadata struct {
	LockID string `json:"lock_id"`
	Token  string `json:"token"`
}

func (dm *DeploymentMetadata) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("type assertion to []byte failed")
	}

	return json.Unmarshal(bytes, dm)
}

func (dm *DeploymentMetadata) Value() (driver.Value, error) {
	return json.Marshal(dm)
}

type AppStateSnapshot struct {
	Machines []*fly.Machine `json:"machines"`
}

func (ass *AppStateSnapshot) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("type assertion to []byte failed")
	}

	return json.Unmarshal(bytes, ass)
}

func (ass *AppStateSnapshot) Value() (driver.Value, error) {
	return json.Marshal(ass)
}

type Deployment struct {
	ID        string             `json:"id"`
	AppID     string             `json:"app_id"`
	Status    string             `json:"status"`
	Metadata  DeploymentMetadata `json:"metadata"`
	State     string             `json:"state"`
	Snapshot  AppStateSnapshot   `json:"snapshot"`
	WorkerID  string             `json:"worker_id"`
	CreatedAt string             `json:"created_at"`
	UpdatedAt string             `json:"updated_at"`
}
