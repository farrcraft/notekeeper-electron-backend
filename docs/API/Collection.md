# Collection API Methods

## User::collections

Request Arguments:

* `shelfId` - Shelf UUID that contains Collections

Response:

* `collections`

## Account::collections

Request Arguments:

* `shelfId` - Shelf UUID that contains Collections

Response:

* `collections`

## User::Collection::create

Request Arguments:

* `name` - Title object describing the Collection name
* `shelfId` - Shelf UUID that contains Collection

Response:

An ID Response

## Account::Collection::create

Request Arguments:

* `name` - Title object describing the Collection name
* `shelfId` - Shelf UUID that contains Collection

Response:

An ID Response

## User::Collection::save

Request Arguments:

* `id` - Collection UUID
* `shelfId` - Shelf UUID
* `name` - Title object describing the Collection name
* `locked` - Whether or not the Collection is locked

Response:

An Empty Response

## Account::Collection::save

Request Arguments:

* `id` - Collection UUID
* `shelfId` - Shelf UUID
* `name` - Title object describing the Collection name
* `locked` - Whether or not the Collection is locked

Response:

An Empty Response

## User::Collection::delete

Request Arguments:

* `id` - Collection UUID
* `shelfId` - Shelf UUID

Response:

An Empty Response

## User::Collection::delete

Request Arguments:

* `id` - Collection UUID
* `shelfId` - Shelf UUID

Response:

An Empty Response