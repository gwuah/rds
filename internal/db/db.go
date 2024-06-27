package db

import (
	"context"
	"embed"

	"github.com/pressly/goose/v3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

type DB struct {
	gorm *gorm.DB
}

func New(dbPath string) (*DB, error) {
	ourDb := &DB{}
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return ourDb, err
	}

	ourDb.gorm = db

	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect(string(goose.DialectSQLite3)); err != nil {
		return ourDb, err
	}

	rawDB, err := db.DB()
	if err != nil {
		return ourDb, err
	}

	if err := goose.Up(rawDB, "migrations"); err != nil {
		return ourDb, err
	}

	return ourDb, err
}

func (db *DB) CreateDeployment(ctx context.Context, deployment Deployment) (*Deployment, error) {
	result := db.gorm.Create(&deployment)
	if result.Error != nil {
		return &deployment, result.Error
	}

	newDeployment := Deployment{}
	result = db.gorm.Where("id = ?", deployment.ID).First(&newDeployment)
	if result.Error != nil {
		return &newDeployment, result.Error
	}

	return &newDeployment, nil
}

func (db *DB) GetDeploymentById(ctx context.Context, id string) (*Deployment, error) {
	deployment := Deployment{}

	result := db.gorm.Where("id = ?", id).First(&deployment)
	if result.Error != nil {
		return &deployment, result.Error
	}

	return &deployment, nil
}

func (db *DB) GetDeployments(ctx context.Context, appID string) (*[]Deployment, error) {
	deployments := []Deployment{}

	result := db.gorm.Where("app_id = ?", appID).Find(&deployments)
	if result.Error != nil {
		return &deployments, result.Error
	}

	return &deployments, nil
}

func (db *DB) GetDeploymentEvents(ctx context.Context, deploymentID string) (*[]Event, error) {
	events := []Event{}

	result := db.gorm.Where("deployment_id = ?", deploymentID).Find(&events)
	if result.Error != nil {
		return &events, result.Error
	}

	return &events, nil
}
