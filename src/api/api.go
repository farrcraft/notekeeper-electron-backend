package api

import (
	"../db"
	"github.com/Sirupsen/logrus"
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
