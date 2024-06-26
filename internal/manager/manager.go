package manager

import (
	"context"

	"connectrpc.com/connect"
	"gorm.io/gorm"

	protov1 "github.com/gwuah/rds/api/gen/proto/v1"
	"github.com/sirupsen/logrus"
)

type Manager struct {
	logger *logrus.Logger
	db     *gorm.DB
}

func New(logger *logrus.Logger, db *gorm.DB) *Manager {
	return &Manager{
		logger: logger,
		db:     db,
	}
}

func (m Manager) CreateDeployment(ctx context.Context, req *connect.Request[protov1.CreateDeploymentRequest]) (*connect.Response[protov1.CreateDeploymentResponse], error) {
	m.logger.Info("[CreateDeployment]")
	return connect.NewResponse(&protov1.CreateDeploymentResponse{}), nil
}

func (m Manager) GetDeployment(ctx context.Context, req *connect.Request[protov1.GetDeploymentRequest]) (*connect.Response[protov1.GetDeploymentResponse], error) {
	m.logger.Info("[GetDeployment]")
	return connect.NewResponse(&protov1.GetDeploymentResponse{}), nil
}

func (m Manager) StopDeployment(ctx context.Context, req *connect.Request[protov1.StopDeploymentRequest]) (*connect.Response[protov1.StopDeploymentResponse], error) {
	m.logger.Info("[StopDeployment]")
	return connect.NewResponse(&protov1.StopDeploymentResponse{}), nil
}
