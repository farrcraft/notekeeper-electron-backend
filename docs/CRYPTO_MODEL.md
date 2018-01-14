# Crypto Model

## RPC

RPC communication uses TLS.


Client & Server sign requests & responses using ed25519 keys.  The first request
a client makes to the server must be a key exchange request.


Client & server keep track of the sequence number of messages sent & received.
The sequence is part of the message envelope & not the payload, so sequence
tampering will not be detected by signature verification.


## DB

* Bolt DB file itself isn't encrypted
* Values in Bolt DB buckets are encrypted using nacl secretbox w/ random nonce per value (prepended to value)
* The content of every DB file is encrypted with its own randomly generated encryption key
* We use scrypt to derive a key from the user's passphrase
* We encrypt DB-file specific encryption keys using the scrypt derived key
* The encrypted encryption key for each DB is stored in its parent DB
* The scrypt derived key is not stored in any db
* Notebooks also have their own encryption keys (even though they don't get their own db's)
* Notebooks encrypt their content (i.e. notes), but not their metadata
* Notebook metadata is encrypted using the encryption key of the DB where the notebook is stored


To access content:

* Take the passphrase input
* Recompute the scrypt derived key which can then 
* Decrypt the DB file encryption key which can 
* Decrypt db values

Decrypted keys can be cached in memory while application is active (like 1password)

## Encryption Keys

There are 3 special encryption keys:

- Passphrase Derived Key
- User Encryption Key
- Account Encryption Key

* The user and account encryption keys are randomly generated when the user/account is created.
* The passphrase key is derived from the passphrase content.
* The user & account encryption key are stored encrypted in the <user UUID>.db as part of the profile bucket value.
* The passphrase key is used to encrypt/decrypt the account & user keys. It is never used for any other content.
* The user key encrypts content in the user DB.
* The account key encrypts content in the account DB.
* Remaining DB types each have their own encryption key, sealed with either the user or account key.

The encrypted version of a DB's encryption key is stored in the metadata index bucket for that DB type.

Key types & locations:

* Master - There is no master key
* Account - <user UUID>.db as part of the profile data, sealed with the passphrase key
* User - <user UUID>.db as part of the profile data, sealed with the passphrase key
* Shelf
* Collection
* Notebook


## Decryption Flow

Here is how everything is decrypted going all the way down from the start.


Given:
* account name
* user email
* user passphrase

Database Files:

* notekeeper.db (master)
    - contains account index
* <account UUID>.db
    - contains profile (account)
    - contains user index
    - contains shelf index
    - contains tag index
* <user UUID>.db
    - contains profile (user)
    - contains shelf index
    - contains tag index
* <user owned shelf UUID>.db
    - contains collection index
    - contains notebooks
    - contains notes
* <account owned shelf UUID>.db
    - contains collection index
    - contains notebooks
    - contains notes
* <user owned collection UUID>.db
    - contains notebooks
    - contains notes
* <account owned collection UUID>.db
    - contains notebooks
    - contains notes

- Open <master DB>

[Signin Flow]
- Look up the account name in the <master DB> *account_index* bucket to find the <account UUID>
    - Iterate <master DB> *account_index* bucket
        - Extract salt embedded in key field
        - Encrypt given *account name* using extracted salt
        - If new encrypted value matches key field value
            - value field contains <account UUID>

- Open <account UUID>.db
- Look up the user email in the <account UUID>.db *user_index* bucket to find the <user UUID>
    - Iterate <account UUID>.db *user_index* bucket
        - Extract salt embedded in key field
        - Encrypt given *user email address* using extracted salt
        - If new encrypted value matches key field value
            - value field contains <user UUID>

- Open <user UUID>.db

- Load user profile data from <user UUID>.db *profile* bucket
- Decrypt profile data using key derived from passphrase
- Set <user UUID>.db EncryptedKey from user UserKey
- Set <account UUID>.db Encrypted key from user AccountKey
- Load account profile data from <account UUID>.db *profile* bucket

-----
Account & User are loaded in memory
Account & User db are open & encrypted keys are in memory

// Need to decide what persistent memory model looks like
// How much state do we persist between RPC calls?
// Q: Is it more of a RESTful model? - A: No
// Or is there something of an active session?
// Q: Or is there something of an FSM that transitions between states? - A: No
For everything below the user db we cache open db's in memory and provide a way
to look up the db key index bucket, but don't actually store any of the content
in memory. We just load it and return it in an rpc response and forget about it.
That way we're not duplicating a bunch of stuff in memory in frontend & backend
and don't have to worry about supporting future workflows.

(Retrieve account shelf list)

