package bluegreen

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/cenkalti/backoff"
	protov1 "github.com/gwuah/rds/api/gen/proto/v1"
	"github.com/sirupsen/logrus"
	fly "github.com/superfly/fly-go"
	"github.com/superfly/fly-go/flaps"
)

var (
	ErrAborted               = errors.New("deployment aborted by user")
	ErrWaitTimeout           = errors.New("wait timeout")
	ErrCreateGreenMachine    = errors.New("failed to create green machines")
	ErrWaitForStartedState   = errors.New("could not get all green machines into started state")
	ErrWaitForHealthy        = errors.New("could not get all green machines to be healthy")
	ErrMarkReadyForTraffic   = errors.New("failed to mark green machines as ready")
	ErrCordonBlueMachines    = errors.New("failed to cordon blue machines")
	ErrStopBlueMachines      = errors.New("failed to stop blue machines")
	ErrWaitForStoppedState   = errors.New("could not get all blue machines into stopped state")
	ErrDestroyBlueMachines   = errors.New("failed to destroy previous deployment")
	ErrValidationError       = errors.New("app not in valid state for bluegreen deployments")
	ErrOrgLimit              = errors.New("app can't undergo bluegreen deployment due to org limits")
	ErrMultipleImageVersions = errors.New("found multiple image versions")

	safeToDestroyValue = "safe_to_destroy"
)

type WorkflowCtx struct {
	ctx               context.Context
	logger            *logrus.Logger
	deploymentRequest *protov1.CreateDeploymentRequest
	deployment        Deployment
	graphqlClient     *fly.Client
	flapsClient       *flaps.Client
	blueMachines      []fly.Machine
	launchInputs      *fly.LaunchMachineInput
}

type Workflow struct {
}

func (w *Workflow) lockApp(wctx WorkflowCtx) error {
	eb := backoff.NewExponentialBackOff()
	eb.MaxElapsedTime = time.Second * 5

	return backoff.RetryNotify(func() error {
		runCtx, cancel := context.WithCancel(wctx.ctx)
		defer cancel()

		response, err := wctx.graphqlClient.LockApp(runCtx, fly.AppLockInput{
			AppID: fmt.Sprint(wctx.deploymentRequest.GetAppId()),
		})
		if err != nil {
			return err
		}

		if response.LockID != "" {
			return nil
		}

		// update deployment in db with lock_id

		return fmt.Errorf("failed to lock app")
	}, eb, func(err error, d time.Duration) {
		// log this somewhere, can be used to build alerting
	})
}

func (w *Workflow) unlockApp(wctx WorkflowCtx) error {
	eb := backoff.NewExponentialBackOff()
	eb.MaxElapsedTime = time.Second * 5

	return backoff.RetryNotify(func() error {
		runCtx, cancel := context.WithCancel(wctx.ctx)
		defer cancel()

		response, err := wctx.graphqlClient.UnlockApp(runCtx, fly.AppLockInput{
			AppID:  fmt.Sprint(wctx.deploymentRequest.GetAppId()),
			LockID: wctx.deployment.LockID,
		})
		if err != nil {
			return err
		}

		if response.Name != "" {
			return nil
		}

		// delete lock_id from deployment in db
		return fmt.Errorf("failed to lock app")
	}, eb, func(err error, d time.Duration) {
		// log this somewhere, can be used to build alerting
	})
}

func (w *Workflow) canPerformDeployment(wctx WorkflowCtx) error {
	eb := backoff.NewExponentialBackOff()
	eb.MaxElapsedTime = time.Second * 5

	return backoff.RetryNotify(func() error {
		runCtx, cancel := context.WithCancel(wctx.ctx)
		defer cancel()

		canPerform, err := wctx.graphqlClient.CanPerformBluegreenDeployment(runCtx, fmt.Sprint(wctx.deploymentRequest.AppId))
		if err != nil {
			return err
		}

		if canPerform {
			return nil
		}

		return backoff.Permanent(ErrOrgLimit)
	}, eb, func(err error, d time.Duration) {
		// log this somewhere, can be used to build alerting
	})
}

func (w *Workflow) detectMultipleImageVersions(wctx WorkflowCtx) error {
	imageToMachineIDs := map[string][]string{}

	for _, mach := range wctx.blueMachines {
		image := mach.ImageRefWithVersion()
		imageToMachineIDs[image] = append(imageToMachineIDs[image], mach.ID)
	}

	if len(imageToMachineIDs) == 1 {
		return nil
	}

	return ErrMultipleImageVersions
}

func (w *Workflow) snapshotStateOfWorld(wctx WorkflowCtx) error {
	eb := backoff.NewExponentialBackOff()
	eb.MaxElapsedTime = time.Second * 5

	return backoff.RetryNotify(func() error {
		runCtx, cancel := context.WithCancel(wctx.ctx)
		defer cancel()

		machines, err := wctx.flapsClient.ListActive(runCtx)
		if err != nil {
			return err
		}

		fmt.Println(len(machines))

		// store snapshots in db
		return fmt.Errorf("failed to lock app")
	}, eb, func(err error, d time.Duration) {
		// log this somewhere, can be used to build alerting
	})
}

func (w *Workflow) createGreenMachines(wctx WorkflowCtx) error {
	eb := backoff.NewExponentialBackOff()
	eb.MaxElapsedTime = time.Second * 5

	return backoff.RetryNotify(func() error {
		runCtx, cancel := context.WithCancel(wctx.ctx)
		defer cancel()

		machines, err := wctx.flapsClient.ListActive(runCtx)
		if err != nil {
			return err
		}

		fmt.Println(len(machines))

		// store snapshots in db
		return fmt.Errorf("failed to lock app")
	}, eb, func(err error, d time.Duration) {
		// log this somewhere, can be used to build alerting
	})
}

// func (w *Workflow) attachCustomTopLevelChecks(wctx WorkflowCtx) error {
// 	for _, entry := range wctx.blueMachines {
// 		for _, service := range entry.launchInput.Config.Services {
// 			servicePort := service.InternalPort
// 			serviceProtocol := service.Protocol

// 			for _, check := range service.Checks {
// 				cc := fly.MachineCheck{
// 					Port:              check.Port,
// 					Type:              check.Type,
// 					Interval:          check.Interval,
// 					Timeout:           check.Timeout,
// 					GracePeriod:       check.GracePeriod,
// 					HTTPMethod:        check.HTTPMethod,
// 					HTTPPath:          check.HTTPPath,
// 					HTTPProtocol:      check.HTTPProtocol,
// 					HTTPSkipTLSVerify: check.HTTPSkipTLSVerify,
// 					HTTPHeaders:       check.HTTPHeaders,
// 				}

// 				if cc.Port == nil {
// 					cc.Port = &servicePort
// 				}

// 				if cc.Type == nil {
// 					cc.Type = &serviceProtocol
// 				}

// 				if entry.launchInput.Config.Checks == nil {
// 					entry.launchInput.Config.Checks = make(map[string]fly.MachineCheck)
// 				}
// 				entry.launchInput.Config.Checks[fmt.Sprintf("bg_deployments_%s", *check.Type)] = cc
// 			}
// 		}
// 	}
// }
