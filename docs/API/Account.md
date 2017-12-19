# Account API Methods

## Account::create

Request Arguments:

* `name` - the name of the account to be created.
* `email` - the email address of the initial account user.
* `passphrase` - the password of the initial account user.

Response:

## Account::unlock

Request Arguments:

* `passphrase` - the password of the user that locked the account.

Response:

## Account::signin

Request Arguments:

* `name` - the name of the account.
* `email` - the email address of an account user.
* `passphrase` - the password of an account user.

Response:

## Account::signout

Request Arguments:

Response:

## Account::lock

Request Arguments:

Response:

## AccountState::get

Request Arguments:

Response:
