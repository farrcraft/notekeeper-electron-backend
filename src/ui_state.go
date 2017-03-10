package main

import (
	"encoding/json"
	"errors"

	"github.com/Sirupsen/logrus"
	"github.com/boltdb/bolt"
)

// UIState contains the saved settings for the UI
type UIState struct {
	WindowWidth  int            `json:"window_width"`
	WindowHeight int            `json:"window_height"`
	DB           *bolt.DB       `json:"-"`
	Logger       *logrus.Logger `json:"-"`
}

// NewUIState returns a new UIState object
func NewUIState(db *bolt.DB, logger *logrus.Logger) *UIState {
	state := &UIState{
		WindowWidth:  -1,
		WindowHeight: -1,
		DB:           db,
		Logger:       logger,
	}
	return state
}

// Load loads the UI's saved state from the database
func (state *UIState) Load() error {
	state.DB.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		bucket := tx.Bucket([]byte("ui_state"))
		cursor := bucket.Cursor()
		key, value := cursor.Seek([]byte("ui_state"))
		if key == nil {
			err := errors.New("Error loading UI State")
			state.Logger.Error(err)
			return err
		}

		err := json.Unmarshal(value, state)
		if err != nil {
			state.Logger.Error("Error decoding UI State json - ", err)
			return err
		}
		return nil
	})
	return nil
}

// Save saves the UI's state to the database
func (state *UIState) Save() error {
	err := state.DB.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte("ui_state"))
		if err != nil {
			state.Logger.Error("Error creating UI State Bucket - ", err)
			return err
		}

		data, err := json.Marshal(state)
		if err != nil {
			state.Logger.Error("Error marshaling UI State - ", err)
			return err
		}

		err = bucket.Put([]byte("ui_state"), data)
		if err != nil {
			state.Logger.Error("Error saving UI State - ", err)
			return err
		}

		return nil
	})
	if err != nil {
		state.Logger.Error("Error saving UI State - ", err)
		return err
	}
	return nil
}
