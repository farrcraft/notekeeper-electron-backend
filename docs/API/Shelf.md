# Shelf API Methods

## User::shelves

Request Arguments:

* `id` - User UUID

Response:

* `shelves`

## Account::shelves

Request Arguments:

* `id` - Account UUID

Response:

* `shelves`

## User::Shelf::create

Request Arguments:

* `name` - Name of the shelf
* `id` - User UUID

Response:

An Id Response

## Account::Shelf::create

Request Arguments:

* `name` - Name of the shelf
* `id` - Account UUID

Response:

An Id Response

## User::Shelf::save

Request Arguments:

* `id` - Shelf UUID
* `ownerId` - User UUID
* `name` - Name of the shelf described as a Title object
* `locked` - Wether the shelf is locked

Response:

An Empty Response

## Account::Shelf::save

Request Arguments:

* `id` - Shelf UUID
* `ownerId` - Account UUID
* `name` - Name of the shelf described as a Title object
* `locked` - Wether the shelf is locked

Response:

An Empty Response

## User::Shelf::delete

Request Arguments:

* `id` - Shelf UUID
* `ownerId` - User UUID

Response:

An Empty Response

## Account::Shelf::delete

Request Arguments:

* `id` - Shelf UUID
* `ownerId` - Account UUID

Response:

An Empty Response
