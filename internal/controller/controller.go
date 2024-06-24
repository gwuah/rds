package controller

import (
	"context"

	"github.com/gwuah/rds/internal/config"
	"github.com/sirupsen/logrus"
)

type Controller struct {
	logger *logrus.Logger
	cfg    *config.Config
}

func (c *Controller) Run(ctx context.Context) error {

	return nil
}
