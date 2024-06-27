package controller

import (
	"context"
	"time"

	"github.com/gwuah/rds/internal/config"
	"github.com/gwuah/rds/internal/db"
	"github.com/sirupsen/logrus"
)

type Controller struct {
	logger *logrus.Logger
	cfg    *config.Config
	db     *db.DB
}

func New(logger *logrus.Logger, cfg *config.Config, db *db.DB) *Controller {
	return &Controller{
		logger: logger,
		cfg:    cfg,
		db:     db,
	}
}

func (c *Controller) Run(ctx context.Context) error {
	every5Seconds := time.NewTicker(time.Second * 5)
	logger := c.logger.WithField("component", "controller")

	for range every5Seconds.C {
		deployments, err := c.db.GetDeployments(ctx, []string{"pending", "recovering"})
		if err != nil {
			logger.WithError(err).Error("failed to get all 'pending' and 'recovering' events")
			continue
		}

		for _, deployment := range *deployments {
			if time.Since(time.Unix(deployment.LastHeartbeat, 0)) > 5*time.Second {
				logger.Infof("deployment %s needs attention. state=%s", deployment.ID, deployment.Status)
				switch deployment.Strategy {
				case "bluegreen":
				}
			}
		}
	}

	return nil
}
