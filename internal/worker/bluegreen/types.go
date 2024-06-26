package bluegreen

import "context"

type Name string

const (
	Noop                             Name = "no_op"
	LockApp                          Name = "lock_app"
	CanPerformDeployment             Name = "can_perform_deployment"
	DetectMultipleImageVersions      Name = "detect_multiple_image_version"
	EnsureHealthchecksAreDefined     Name = "ensure_healthchecks_are_defined"
	SnapshotAppStateBeforeDeployment Name = "snapshot_app_state_before_deployment"
	CreateGreenMachines              Name = "create_green_machines"
	WaitForGreenMachinesToStart      Name = "wait_for_green_machines_to_start"
	WaitForGreenMachinesToBeHealthy  Name = "wait_for_green_machines_to_be_healthy"
	UncordonGreenMachines            Name = "uncordon_green_machines"
	WaitForMachineInfoPropagation    Name = "wait_for_machine_info_propagation"
	CordonBlueMachines               Name = "cordon_blue_machines"
	StopBlueMachines                 Name = "stop_blue_machines"
	WaitForBlueMachinesToStop        Name = "wait_for_blue_machines_to_stop"
	DestroyBlueMachines              Name = "destroy_blue_machines"
	UnlockApp                        Name = "unlock_app"
)

type Task struct {
	Current  Name
	Next     Name
	Previous Name
	Action   func(ctx context.Context) error
	Rollback func(ctx context.Context) error
}

type Deployment struct {
	ID              string `json:"id"`
	State           string `json:"state"`
	Token           string `json:"token"`
	Checkpoint      string `json:"checkpoint"`
	CurrentSnapshot string `json:"current_snapshot"`
	WorkerID        string `json:"worker_id"`
	CreatedAt       string `json:"created_at"`
	UpdatedAT       string `json:"updated_at"`
	Metadata        string `json:"metadata"`
	LockID          string `json:"lock_id"`
}
