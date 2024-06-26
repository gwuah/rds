package bluegreen

import (
	"context"
	"time"

	"github.com/gwuah/rds/internal/circuit_breaker"
	"github.com/gwuah/rds/lib"

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

func New(
	ctx context.Context,
	logger *logrus.Logger,
	deployment Deployment,
	circuitBreaker *circuit_breaker.CircuitBreaker,
) (*BlueGreen, error) {

	cancelChan := make(chan lib.CancelMessage)
	workerCtx, cancel := context.WithCancel(ctx)

	flapsClient, err := flaps.NewWithOptions(context.Background(), flaps.NewClientOpts{
		Transport: &fly.Transport{
			UnderlyingTransport: circuitBreaker,
		},
	})
	if err != nil {
		return nil, err
	}

	bg := BlueGreen{
		cancelChan: cancelChan,
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

	go bg.RunHeartbeat(workerCtx)
	go bg.HandleUserCancel(workerCtx, cancelChan, cancel)
	return &bg, nil
}

func (bg *BlueGreen) RunHeartbeat(ctx context.Context) error {
	ticker := time.NewTicker(time.Second * 5)

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			// update deployment with heartbeat
		}
	}
}

func (bg *BlueGreen) HandleUserCancel(ctx context.Context, cancelChan chan lib.CancelMessage, cancel context.CancelFunc) error {

	select {
	case <-ctx.Done():
		return nil
	case msg := <-cancelChan:
		if msg.DeploymentID == bg.deployment.ID {
			cancel()
		}
	}

	return nil
}

func (bg *BlueGreen) Run(ctx context.Context) error {

	return nil
}
