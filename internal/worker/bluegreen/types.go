package worker

import "context"

type Name string

const (
	Noop                          Name = "no_op"
	LockApp                       Name = "lock_app"
	CanPerformBluegreenDeployment Name = "lock_app"
)

type Task struct {
	Current  Name
	Next     Name
	Previous Name
	Action   func(ctx context.Context) error
	Rollback func(ctx context.Context) error
}
