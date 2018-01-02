
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

### account_index

This bucket provides a mapping between account names and account UUID's.

* `key` - encryption key derived from account name w/ salt embedded
* `value` - unencrypted account UUID
