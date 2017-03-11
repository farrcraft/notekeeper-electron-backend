package account

import (
	"crypto/subtle"
	"encoding/json"
	"errors"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/boltdb/bolt"
	uuid "github.com/satori/go.uuid"
)

// UserProfile contains the minimal user information that is visible to all users of an account
type UserProfile struct {
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

// User is a single user in an account
type User struct {
	ID            uuid.UUID      `json:"id"`          // ID is the unique identifier of the user
	Profile       *UserProfile   `json:"profile"`     // Profile is the user information that is visible to all users in an account
	Active        bool           `json:"-"`           // Active indicates whether the user is active or not
	Account       *Account       `json:"-"`           // Account is the account that the user belongs to
	Created       time.Time      `json:"created"`     // Created is the time when the user was created
	Updated       time.Time      `json:"updated"`     // Updated is the time when the user was last created
	AccountKey    []byte         `json:"account_key"` // AccountKey is the encrypted version of the account-level encryption key
	PassphraseKey []byte         `json:"-"`           // PassphraseKey is the key derived from the passphrase
	Salt          []byte         `json:"-"`           // Salt is the unique salt for generating the passphrase key
	Shelves       []*Shelf       `json:"-"`           // Shelves is the set of shelves that belong to the user
	Logger        *logrus.Logger `json:"-"`
	DB            *bolt.DB       `json:"-"`
}

// NewUser creates a new user object
func NewUser(db *bolt.DB, logger *logrus.Logger, email string) *User {
	now := time.Now()
	user := &User{
		ID: uuid.NewV4(),
		Profile: &UserProfile{
			Email: email,
		},
		Active:  true,
		Created: now,
		Updated: now,
		DB:      db,
		Logger:  logger,
	}
	return user
}

// Lookup looks for a user id in the account database
func (user *User) Lookup() error {
	originalID := user.ID
	user.ID = uuid.Nil
	user.DB.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		bucket := tx.Bucket([]byte("user_map"))
		cursor := bucket.Cursor()

		for key, value := cursor.First(); key != nil; key, value = cursor.Next() {
			// extract the salt from the existing encrypted email
			encryptedEmail, salt := ExtractSalt(key)
			// create a new key using the extracted salt and the unencrypted email we're searching for
			checkEmail, err := DeriveKey([]byte(user.Profile.Email), salt[:])
			if err != nil {
				user.Logger.Error("Error deriving user map key - ", err)
				return err
			}
			// the new key should match the existing key if we have the right email and salt
			if subtle.ConstantTimeCompare(encryptedEmail, checkEmail[:]) == 1 {
				user.ID, err = uuid.FromBytes(value)
				// the salt stored in the email key is also the primary user passphrase key salt
				user.Salt = salt[:]
				if err != nil {
					user.Logger.Error("Error converting email map uuid - ", err)
					return err
				}
				return nil
			}
		}
		return nil
	})
	if user.ID == uuid.Nil {
		user.ID = originalID
		return errors.New("no user found")
	}
	return nil
}

// Load loads the user data for a user from the account database
func (user *User) Load(passphrase string) error {
	user.DB.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		bucket := tx.Bucket([]byte("users"))
		cursor := bucket.Cursor()
		key, value := cursor.Seek(user.ID.Bytes())
		if key == nil {
			err := errors.New("Error loading user")
			user.Logger.Error(err)
			return err
		}

		// salt will have already been provided from a previous Lookup() operation
		passphraseKey, err := DeriveKey([]byte(passphrase), user.Salt)
		if err != nil {
			user.Logger.Error("Error deriving key from passphrase - ", err)
			return err
		}
		// need to decrypt value
		decryptedData, err := Open(passphraseKey[:], value)
		if err != nil {
			// this error condition is most likely caused by an incorrect passphrase
			user.Logger.Error("Error decrypting user data - ", err)
			return err
		}

		err = json.Unmarshal(decryptedData, user)
		if err != nil {
			user.Logger.Error("Error decoding user json - ", err)
			return err
		}
		user.PassphraseKey = passphraseKey[:]
		Zero(passphraseKey[:])
		return nil
	})
	return nil
}

// Save saves the user to the database
func (user *User) Save() error {
	err := user.DB.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte("users"))
		if err != nil {
			user.Logger.Error("Error creating users bucket - ", err)
			return err
		}
		data, err := json.Marshal(user)
		if err != nil {
			user.Logger.Error("Error marshaling user - ", err)
			return err
		}

		encryptedData, err := Seal(user.PassphraseKey, data)
		if err != nil {
			user.Logger.Error("Error encrypting user content - ", err)
			return err
		}

		err = bucket.Put(user.ID.Bytes(), encryptedData)
		if err != nil {
			user.Logger.Error("Error saving user - ", err)
			return err
		}
		return nil
	})
	if err != nil {
		user.Logger.Error("Error saving user - ", err)
		return err
	}
	// Since users are keyed by only an unencrypted id in the db
	// we also need to store a mapping between a key derived from the email address and the id
	// otherwise there is no way to look up a user without taking a brute force decryption test approach
	err = user.DB.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte("user_map"))
		if err != nil {
			user.Logger.Error("Error creating user_map bucket - ", err)
			return err
		}
		encryptedEmail, err := DeriveKey([]byte(user.Profile.Email), user.Salt)
		if err != nil {
			user.Logger.Error("Error creating user map key - ", err)
			return err
		}
		err = bucket.Put(encryptedEmail[:], user.ID.Bytes())
		if err != nil {
			user.Logger.Error("Error saving user map - ", err)
			return err
		}
		return nil
	})
	if err != nil {
		user.Logger.Error("Error mapping user - ", err)
		return err
	}
	return err
}

// CreateKeys creates the account and user key from a passphrase
func (user *User) CreateKeys(passphrase []byte) error {
	// generate account-level encryption key
	accountKey, err := GenerateKey()
	if err != nil {
		return err
	}
	// derive key from passphrase
	var key = new([KeySize]byte)
	key, user.Salt, err = DeriveKeyAndSalt(passphrase)
	if err != nil {
		return err
	}
	slicedKey := key[:]
	user.PassphraseKey = append(user.Salt, slicedKey...)
	user.AccountKey, err = Seal(user.PassphraseKey, accountKey[:])
	if err != nil {
		return err
	}
	return nil
}
