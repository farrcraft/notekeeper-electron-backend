package uistate

import (
	"encoding/json"

	"notekeeper-electron-backend/codes"
	"notekeeper-electron-backend/db"

	"github.com/sirupsen/logrus"
	"go.etcd.io/bbolt"
)

// UIState contains the saved settings for the UI
// The UI state includes only the most generic persistent UI settings
// !!!!The stored representation **IS NOT ENCRYPTED**!!!!
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
	DBRegistry       *db.Registry   `json:"-"`
	Logger           *logrus.Logger `json:"-"`
}

// NewUIState returns a new UIState object
func NewUIState(registry *db.Registry, logger *logrus.Logger) *UIState {
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
		DBRegistry:       registry,
		Logger:           logger,
	}
	return state
}

// Create creates a default UI state if none exists yet
func (state *UIState) Create() error {
	// for now there is only a single state that is kept in the master db
	if state.DBRegistry == nil || state.DBRegistry.Master == nil {
		state.Logger.Warn("ui state create - missing db")
		code := codes.New(codes.ScopeUIState, codes.ErrorMissingDB)
		return code
	}
	err := state.DBRegistry.Master.DB.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte("ui_state"))
		if bucket == nil {
			bucket, err := tx.CreateBucket([]byte("ui_state"))
			if err != nil {
				state.Logger.Warn("Error creating ui_state bucket - ", err)
				code := codes.New(codes.ScopeUIState, codes.ErrorCreateBucket)
				return code
			}
			data, err := json.Marshal(state)
			if err != nil {
				state.Logger.Warn("Error marshaling default UI State - ", err)
				code := codes.New(codes.ScopeUIState, codes.ErrorMarshal)
				return code
			}

			err = bucket.Put([]byte("ui_state"), data)
			if err != nil {
				state.Logger.Warn("Error writing default UI State - ", err)
				code := codes.New(codes.ScopeUIState, codes.ErrorWriteBucket)
				return code
			}
		}
		return nil
	})
	if err != nil {
		if codes.IsInternalError(err) {
			return err
		}
		state.Logger.Warn("error saving default ui state - ", err)
		code := codes.New(codes.ScopeUIState, codes.ErrorSave)
		return code
	}
	return nil
}

// Load loads the UI's saved state from the database
func (state *UIState) Load() error {
	// [FIXME] precondition - expecting open master db here
	err := state.DBRegistry.Master.DB.View(func(tx *bbolt.Tx) error {
		// Assume bucket exists and has keys
		bucket := tx.Bucket([]byte("ui_state"))
		cursor := bucket.Cursor()
		key, value := cursor.Seek([]byte("ui_state"))
		if key == nil {
			state.Logger.Warn("Error loading UI State")
			code := codes.New(codes.ScopeUIState, codes.ErrorLoad)
			return code
		}

		err := json.Unmarshal(value, state)
		if err != nil {
			state.Logger.Warn("Error decoding UI State json - ", err)
			code := codes.New(codes.ScopeUIState, codes.ErrorDecode)
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
	// [FIXME] precondition - expecting open master db here
	err := state.DBRegistry.Master.DB.Update(func(tx *bbolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte("ui_state"))
		if err != nil {
			state.Logger.Warn("Error creating UI State Bucket - ", err)
			code := codes.New(codes.ScopeUIState, codes.ErrorCreateBucket)
			return code
		}

		data, err := json.Marshal(state)
		if err != nil {
			state.Logger.Warn("Error marshaling UI State - ", err)
			code := codes.New(codes.ScopeUIState, codes.ErrorMarshal)
			return code
		}

		err = bucket.Put([]byte("ui_state"), data)
		if err != nil {
			state.Logger.Warn("Error writing UI State - ", err)
			code := codes.New(codes.ScopeUIState, codes.ErrorWriteBucket)
			return code
		}

		return nil
	})
	if err != nil {
		if codes.IsInternalError(err) {
			return err
		}
		state.Logger.Warn("Error saving UI State - ", err)
		code := codes.New(codes.ScopeUIState, codes.ErrorSave)
		return code
	}
	return nil
}
