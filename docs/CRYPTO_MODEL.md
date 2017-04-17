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

Look up the account UUID in the master DB account_map bucket


