

# Account Database

The account database file uses the account UUID as its filename with the
`.db` file extension.


## Buckets

### profile

This bucket stores common account profile data.

* `key` - unencrypted account UUID
* `value` - serialized JSON encrypted w/ account-level encryption key

### user_index

This bucket provides a mapping between user email addresses and user UUID's.

* `key` - encryption key derived from user email address w/ salt embedded
* `value` - unencrypted user UUID

The embedded salt is the user's fixed encryption key salt.

### user_profiles

This bucket stores user profile data for account users.
The user profile information stored here is visble to all users in the account

* `key` - unencrypted user UUID
* `value` - serialized JSON encrypted w/ account-level encryption key

### shelf_index

### tags

This bucket contains user-level tags.
