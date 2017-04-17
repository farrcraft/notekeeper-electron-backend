package api

import (
	"../db"
	"github.com/Sirupsen/logrus"
)

// API interface
type API struct {
	DBFactory *db.Factory
	Logger    *logrus.Logger
}

// New creates a new API object
func New(factory *db.Factory, logger *logrus.Logger) *API {
	api := &API{
		DBFactory: factory,
		Logger:    logger,
	}
	return api
}
