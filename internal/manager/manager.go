package manager

import (
	"context"
	"time"

	"connectrpc.com/connect"
	"github.com/sirupsen/logrus"
	"github.com/sqids/sqids-go"
	"github.com/superfly/fly-go"
	"github.com/superfly/fly-go/flaps"
	"github.com/superfly/fly-go/tokens"

	protov1 "github.com/gwuah/rds/api/gen/proto/v1"
	"github.com/gwuah/rds/internal/db"
	"github.com/gwuah/rds/libs/circuit_breaker"
)

type Manager struct {
	circuitBreaker *circuit_breaker.CircuitBreaker
	logger         *logrus.Logger
	db             *db.DB
}

func New(logger *logrus.Logger, db *db.DB, circuitBreaker *circuit_breaker.CircuitBreaker) *Manager {
	return &Manager{
		logger:         logger,
		db:             db,
		circuitBreaker: circuitBreaker,
	}
}

func createID() (string, error) {
	s, err := sqids.New(sqids.Options{
		MinLength: 6,
	})

	if err != nil {
		return "", err
	}

	id, err := s.Encode([]uint64{uint64(time.Now().UTC().Unix())}) // "86Rf07"
	if err != nil {
		return "", err
	}

	return id, err
}

func (m Manager) CreateDeployment(ctx context.Context, req *connect.Request[protov1.CreateDeploymentRequest]) (*connect.Response[protov1.CreateDeploymentResponse], error) {
	logger := m.logger.WithFields(logrus.Fields{
		"method": "CreateDeployment",
	})

	id, err := createID()
	if err != nil {
		logger.WithError(err).Error("failed to create deployment id")
		return connect.NewResponse(&protov1.CreateDeploymentResponse{}), err
	}

	flapsClient, err := flaps.NewWithOptions(context.Background(), flaps.NewClientOpts{
		Tokens: tokens.Parse(req.Msg.Token),
		Transport: &fly.Transport{
			UnderlyingTransport: m.circuitBreaker,
		},
	})
	if err != nil {
		return nil, err
	}

	machines, err := flapsClient.ListActive(ctx)
	if err != nil {
		return nil, err
	}

	deployment, err := m.db.CreateDeployment(ctx, db.Deployment{
		ID:     id,
		AppID:  req.Msg.AppId,
		Status: "pending",
		Metadata: db.DeploymentMetadata{
			Token: req.Msg.Token,
		},
		Snapshot: db.AppStateSnapshot{
			Machines: machines,
		},
	})

	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&protov1.CreateDeploymentResponse{
		Id: deployment.ID,
	}), nil
}

func (m Manager) GetDeployment(ctx context.Context, req *connect.Request[protov1.GetDeploymentRequest]) (*connect.Response[protov1.GetDeploymentResponse], error) {
	m.logger.Info("[GetDeployment]")
	return connect.NewResponse(&protov1.GetDeploymentResponse{}), nil
}

func (m Manager) StopDeployment(ctx context.Context, req *connect.Request[protov1.StopDeploymentRequest]) (*connect.Response[protov1.StopDeploymentResponse], error) {
	m.logger.Info("[StopDeployment]")
	return connect.NewResponse(&protov1.StopDeploymentResponse{}), nil
}
