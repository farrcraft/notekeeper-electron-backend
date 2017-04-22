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


## Decryption Flow

Here is how everything is decrypted going all the way down from the start.


Given:
* account name
* user email
* use passphrase

Database Files:

* notekeeper.db
    - contains account index
* <account UUID>.db
    - contains account
    - contains user index
    - contains shelf index
* <user UUID>.db
    - contains user
    - contains shelf index
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

(Signin Flow)
Look up the account name in the master DB account_map bucket to find the <account UUID>
Open <account UUID>.db
Look up the user email in the <account UUID>.db user_map bucket to find the <user UUID>
Open <user UUID>.db
Load user from <user UUID>.db; decrypt using passphrase key
Set <user UUID>.db EncryptedKey from user UserKey
Load account from <account UUID>.db
Set <account UUID>.db Encrypted key from user AccountKey
-----
Account & User are loaded in memory
Account & User db are open & encrypted keys are in memory

// Need to decide what persistent memory model looks like
// How much state do we persist between RPC calls?
// Is it more of a RESTful model?
// Or is there something of an active session?
// Or is there something of an FSM that transitions between states?
For everything below the user db we cache open db's in memory and provide a way
to look up the db key index bucket, but don't actually store any of the content
in memory. We just load it and return it in an rpc response and forget about it.
That way we're not duplicating a bunch of stuff in memory in frontend & backend
and don't have to worry about supporting future workflows.

(Retrieve account shelf list)

