# Overview

There are two separate databases - a master system database and an account
specific database.  Both databases are stored in the user data directory
of the application.  BoltDB is used as the database engine.


# Master Database

The master database file is named `notekeeper.db`.  It stores only minimal
application global information along with required reference points in order
to locate any linked account databases.


## Buckets

### ui_state

This bucket stores generic information about the UI state - window size,
position, etc.  It contains no PII data.

* `key` - "ui_state"
* `value` - unencrypted serialized JSON

### account_map

This bucket provides a mapping between account names and account UUID's.

* `key` - encryption key derived from account name w/ salt embedded
* `value` - unencrypted account UUID


### accounts

This bucket stores common account-level data for any accounts that have
been created.

* `key` - unencrypted account UUID
* `value` - serialized JSON encrypted w/ account-level encryption key


# Account Database

The account database file uses the account UUID as its filename with the
`.db` file extension.


## Buckets

### user_map

This bucket provides a mapping between user email addresses and user UUID's.

* `key` - encryption key derived from user email address w/ salt embedded
* `value` - unencrypted user UUID

The embedded salt is the user's fixed encryption key salt.


### users

This bucket stores user profile data for account users.

* `key` - unencrypted user UUID
* `value` - serialized JSON encrypted w/ user-level encryption key
