package bluegreen

import (
	"context"
	"time"

	lib "github.com/gwuah/rds/libs"
	"github.com/gwuah/rds/libs/circuit_breaker"

	"github.com/sirupsen/logrus"
	fly "github.com/superfly/fly-go"
	"github.com/superfly/fly-go/flaps"
)

const (
	webURL = "https://api.fly.io"
)

type BlueGreen struct {
	cancelChan chan lib.CancelMessage
	ctx        context.Context
	token      string
	tasks      map[Name]Task
	// deploymentRequest *protov1.CreateDeploymentRequest
	graphqlClient *fly.Client
	flapsClient   *flaps.Client
	logger        *logrus.Logger
	deployment    Deployment
}

func Run(
	ctx context.Context,
	logger *logrus.Logger,
	deployment Deployment,
	circuitBreaker *circuit_breaker.CircuitBreaker,
) (chan lib.CancelMessage, error) {

	flapsClient, err := flaps.NewWithOptions(context.Background(), flaps.NewClientOpts{
		Transport: &fly.Transport{
			UnderlyingTransport: circuitBreaker,
		},
	})
	if err != nil {
		return nil, err
	}

	bg := BlueGreen{
		// deploymentRequest: req,
		logger:      logger,
		flapsClient: flapsClient,
		graphqlClient: fly.NewClientFromOptions(fly.ClientOptions{
			BaseURL:     webURL,
			AccessToken: deployment.Token,
			Transport: &fly.Transport{
				UnderlyingTransport: circuitBreaker,
			},
		}),
	}

	cancelChan := make(chan lib.CancelMessage)
	workerCtx, cancel := context.WithCancel(ctx)
	tasks := map[Name]Task{
		LockApp: {
			Current:  LockApp,
			Previous: Noop,
			Next:     CanPerformDeployment,
			Action: func(ctx context.Context) error {
				return nil
			},
			Rollback: func(ctx context.Context) error {
				return nil
			},
		},
	}

	bg.tasks = tasks
	bg.cancelChan = cancelChan

	go func() {
		ticker := time.NewTicker(time.Second * 5)
		for {
			select {
			case <-workerCtx.Done():
			case <-ticker.C:
				// update deployment with heartbeat
			}
		}
	}()
	go func() {
		select {
		case <-workerCtx.Done():
		case msg := <-cancelChan:
			if msg.DeploymentID == bg.deployment.ID {
				cancel()
			}
		}
	}()

	return cancelChan, nil
}
