package api

import (
	"notekeeper-electron-backend/db"

	"github.com/sirupsen/logrus"
)

// API interface
type API struct {
	DBRegistry *db.Registry
	Logger     *logrus.Logger
}

// New creates a new API object
func New(registry *db.Registry, logger *logrus.Logger) *API {
	api := &API{
		DBRegistry: registry,
		Logger:     logger,
	}
	return api
}
