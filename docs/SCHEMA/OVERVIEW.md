# Overview

* Databases are stored in the user data directory of the frontend application.
* BoltDB is used as the database engine.
* There is a single master database with a common filename.
* All other databases are named using the UUID of the container they store data for.


## Database Types

* Master
* Account (1 per account UUID)
* User (1 per user UUID)
* Shelf (1 per shelf UUID)
* Collection (1 per collection UUID)


## Database Indices

Each database contains index buckets for any databases directly below it.


* Master DB contains index of accounts
* account DBs contain index of users, & account-level shelves
* user dbs contain index of user-level shelves
* shelf dbs contain index of collections & notebook content
* collection dbs contain notebook content


master db -> account db -> user db -> shelf -> collection

