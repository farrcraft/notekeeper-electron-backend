# User Database

## Buckets

### profile

This bucket stores user profile data.

* `key` - unencrypted <user UUID>
* `value` - serialized JSON encrypted w/ passphrase derived key

### shelf_index

### tags

This bucket contains user-level tags.
