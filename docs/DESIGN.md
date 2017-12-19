# Features

Encryption at rest
Encryption in transit
Decentralized (mostly)
verified message authentication when syncing
Windows application
iOS application
MacOS application
Web application
Cross Application syncing
Notebooks
Nested notebooks
Tagging
Search
Lists
Spreadsheet notes
Rich Text Formatting
Note/Notebook sharing
2 factor authentication
a view of columns of notes (columns being notebooks?) to enable project management workflows? e.g. - https://waffle.io/

Ability to nuke encryption keys from an application install
Ability to sync encryption keys to web
Ability to have web setting for burn after read for encryption keys
Enable/disable sync
3 types of sync: key-only sync, data only sync, full sync
Nuking encryption keys is an advanced feature that must be opted-in
Default behavior is to sync encryption keys when they are needed (if a device is missing them or has an outdated copy)
Nuking encryption keys from an application install makes it impossible to access the account's notebooks on that install*
*Ignoring forensic reconstruction of the key material via external recovery methods
Encryption keys must be synced from another install in order to recover access to the notebooks
If all copies of the encryption keys are nuked from all application installs (including web), notebook data is lost forever
With granular control of how, when, and where key material is transmitted between device installs and our servers you have greater peace of mind of the safety of your encrypted notebooks.  Even if you do choose to transit and store your key material on our servers, it is still protected with your own unique passphrase which we never store.


# Application Layers

RPC Server
RPC Method Handlers
API
Domain Objects
Database
Crypto Primitives

# Account / User Design

How do accounts & users interact?
sharing of notes, notebooks, shelves between users & accounts
share individual notes or notebooks or whole shelves?

if there is sharing granularity at all levels:

private shelf
     public notebook
          private note
          public note
     private notebook
          public note
          private note
public shelf
     public notebook
          private note
          public note
     private notebook
          public note
          private note

a lot of these combinations don't make sense

There are two separate issues:
- accessibility of content between account level & user level
- how content can be shared

shelfs can be owned by either the account or a single user
all users in an account can access content in account shelves
user shelves can only be accessed by the user they belong to
each account has an account-level encryption key
each user has their own personal encryption key
when new users are created, they receive a copy of the account-level encryption key
the first user is created at the same time as the account
the account key is created at the same time & passed down to the next new user

