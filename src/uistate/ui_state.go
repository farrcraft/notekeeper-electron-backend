package uistate

import (
	"encoding/json"

	"../codes"
	"../db"
	"github.com/Sirupsen/logrus"
	"github.com/boltdb/bolt"
)

// UIState contains the saved settings for the UI
// The UI state includes only the most generic persistent UI settings
// !!!!The stored representation **IS NOT ENCRYPTED**!!!
type UIState struct {
	WindowWidth      int32          `json:"window_width" mapstructure:"window_width"`
	WindowHeight     int32          `json:"window_height" mapstructure:"window_height"`
	WindowXPosition  int32          `json:"window_x_position" mapstructure:"window_x_position"`
	WindowYPosition  int32          `json:"window_y_position" mapstructure:"window_y_position"`
	WindowMaximized  bool           `json:"window_maximized" mapstructure:"window_maximized"`
	WindowMinimized  bool           `json:"window_minimized" mapstructure:"window_minimized"`
	WindowFullscreen bool           `json:"window_fullscreen" mapstructure:"window_fullscreen"`
	DisplayWidth     int32          `json:"display_width" mapstructure:"display_width"`
	DisplayHeight    int32          `json:"display_height" mapstructure:"display_height"`
	DisplayXPosition int32          `json:"display_x_position" mapstructure:"display_x_position"`
	DisplayYPosition int32          `json:"display_y_position" mapstructure:"display_y_position"`
	DB               *db.DB         `json:"-"`
	Logger           *logrus.Logger `json:"-"`
}

// NewUIState returns a new UIState object
func NewUIState(db *db.DB, logger *logrus.Logger) *UIState {
	state := &UIState{
		WindowWidth:      -1,
		WindowHeight:     -1,
		WindowXPosition:  0,
		WindowYPosition:  0,
		WindowMaximized:  false,
		WindowMinimized:  false,
		WindowFullscreen: false,
		DisplayWidth:     -1,
		DisplayHeight:    -1,
		DisplayXPosition: -1,
		DisplayYPosition: -1,
		DB:               db,
		Logger:           logger,
	}
	return state
}

// Create creates a default UI state if none exists yet
func (state *UIState) Create() error {
	if state.DB == nil {
		state.Logger.Debug("ui state create - missing db")
		code := codes.New(codes.ErrorUIStateMissingDb)
		return code
	}
	err := state.DB.DB.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("ui_state"))
		if bucket == nil {
			bucket, err := tx.CreateBucket([]byte("ui_state"))
			if err != nil {
				state.Logger.Debug("Error creating ui_state bucket - ", err)
				code := codes.New(codes.ErrorUIStateCreateBucket)
				return code
			}
			data, err := json.Marshal(state)
			if err != nil {
				state.Logger.Debug("Error marshaling default UI State - ", err)
				code := codes.New(codes.ErrorDefaultUIStateMarshal)
				return code
			}

			err = bucket.Put([]byte("ui_state"), data)
			if err != nil {
				state.Logger.Debug("Error writing default UI State - ", err)
				code := codes.New(codes.ErrorDefaultUIStateWrite)
				return code
			}
		}
		return nil
	})
	if err != nil {
		if codes.IsInternalError(err) {
			return err
		}
		state.Logger.Debug("error saving default ui state - ", err)
		code := codes.New(codes.ErrorDefaultUIStateSave)
		return code
	}
	return nil
}

// Load loads the UI's saved state from the database
func (state *UIState) Load() error {
	err := state.DB.DB.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		bucket := tx.Bucket([]byte("ui_state"))
		cursor := bucket.Cursor()
		key, value := cursor.Seek([]byte("ui_state"))
		if key == nil {
			state.Logger.Debug("Error loading UI State")
			code := codes.New(codes.ErrorLoadUIState)
			return code
		}

		err := json.Unmarshal(value, state)
		if err != nil {
			state.Logger.Debug("Error decoding UI State json - ", err)
			code := codes.New(codes.ErrorUIStateDecode)
			return code
		}
		return nil
	})
	if err != nil {
		if codes.IsInternalError(err) {
			return err
		}
		// [FIXME] - handle unknown error
	}

	return nil
}

// Save saves the UI's state to the database
func (state *UIState) Save() error {
	err := state.DB.DB.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte("ui_state"))
		if err != nil {
			state.Logger.Debug("Error creating UI State Bucket - ", err)
			code := codes.New(codes.ErrorUIStateBucket)
			return code
		}

		data, err := json.Marshal(state)
		if err != nil {
			state.Logger.Debug("Error marshaling UI State - ", err)
			code := codes.New(codes.ErrorUIStatemarshal)
			return code
		}

		err = bucket.Put([]byte("ui_state"), data)
		if err != nil {
			state.Logger.Debug("Error writing UI State - ", err)
			code := codes.New(codes.ErrorUIStateWrite)
			return code
		}

		return nil
	})
	if err != nil {
		if codes.IsInternalError(err) {
			return err
		}
		state.Logger.Debug("Error saving UI State - ", err)
		code := codes.New(codes.ErrorUIStateSave)
		return code
	}
	return nil
}
