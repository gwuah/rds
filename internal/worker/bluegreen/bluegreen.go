package worker

import (
	"context"
	"fmt"
	"time"

	"github.com/cenkalti/backoff"
	protov1 "github.com/gwuah/rds/api/gen/proto/v1"
	"github.com/gwuah/rds/internal/breaker"

	"github.com/sirupsen/logrus"
	fly "github.com/superfly/fly-go"
)

const (
	webURL = "https://api.fly.io"
)

type BlueGreen struct {
	ctx               context.Context
	token             string
	tasks             map[Name]Task
	deploymentRequest *protov1.CreateDeploymentRequest
	graphqlClient     *fly.Client
	logger            *logrus.Logger
}

func New(
	ctx context.Context,
	logger *logrus.Logger,
	req *protov1.CreateDeploymentRequest,
	breaker *breaker.Breaker,
) *BlueGreen {
	bg := BlueGreen{
		deploymentRequest: req,
		logger:            logger,
		graphqlClient: fly.NewClientFromOptions(fly.ClientOptions{
			BaseURL:     webURL,
			AccessToken: req.Token,
			Transport: &fly.Transport{
				UnderlyingTransport: breaker,
			},
		}),
	}

	tasks := map[Name]Task{
		LockApp: {
			Current:  LockApp,
			Previous: Noop,
			Next:     CanPerformBluegreenDeployment,
			Action: func(ctx context.Context) error {
				return bg.lockApp(ctx)
			},
			Rollback: func(ctx context.Context) error {
				return bg.unlockApp(ctx)
			},
		},
	}

	bg.tasks = tasks
	return &bg
}

func (bg *BlueGreen) lockApp(ctx context.Context) error {
	eb := backoff.NewExponentialBackOff()
	eb.MaxElapsedTime = time.Second * 5

	return backoff.RetryNotify(func() error {
		runCtx, cancel := context.WithCancel(bg.ctx)
		defer cancel()

		response, err := bg.graphqlClient.LockApp(runCtx, fly.AppLockInput{
			AppID: fmt.Sprint(bg.deploymentRequest.GetAppId()),
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

func (bg *BlueGreen) unlockApp(ctx context.Context) error {
	return nil
}

func (bg *BlueGreen) Run(ctx context.Context) error {
	return nil
}
