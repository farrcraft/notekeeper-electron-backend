package db

import (
	"encoding/json"
	"fmt"

	"../codes"
	"../crypto"
	"go.etcd.io/bbolt"
)

// LoadEncryptedKey loads the encrypted key from an index bucket
// [FIXME] doesn't work for account-scoped db's - not sure we actually need this method anyway. maybe delete?
func (registry *Registry) LoadEncryptedKey(dbKey Key, passphraseKey []byte, handle *Handle) ([]byte, error) {
	var encryptedKey []byte
	var bucketName []byte
	bucketName = []byte(fmt.Sprint(TypeToStr(dbKey.Type), "_index"))
	err := handle.DB.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket(bucketName)
		if bucket == nil {
			registry.Logger.Debug(bucketName, " bucket does not exist")
			code := codes.New(codes.ScopeDB, codes.ErrorBucketMissing)
			return code
		}

		cursor := bucket.Cursor()
		key, value := cursor.Seek(dbKey.ID.Bytes())
		if key == nil {
			registry.Logger.Debug("Error loading record from index [", bucketName, "]")
			code := codes.New(codes.ScopeDB, codes.ErrorLoad)
			return code
		}

		c := crypto.New(registry.Logger)
		encryptionKey, err := c.Open(passphraseKey, handle.EncryptedKey)
		if err != nil {
			registry.Logger.Debug("Error opening key - ", err)
			code := codes.New(codes.ScopeDB, codes.ErrorOpenKey)
			return code
		}

		// decrypt value
		decryptedData, err := c.Open(encryptionKey, value)
		if err != nil {
			registry.Logger.Debug("Error decrypting data - ", err)
			code := codes.New(codes.ScopeDB, codes.ErrorDecrypt)
			return code
		}

		entry := &IndexEntry{}
		err = json.Unmarshal(decryptedData, entry)
		if err != nil {
			registry.Logger.Debug("Error decoding json - ", err)
			code := codes.New(codes.ScopeDB, codes.ErrorDecode)
			return code
		}

		encryptedKey = entry.EncryptedKey
		return nil
	})
	return encryptedKey, err
}
