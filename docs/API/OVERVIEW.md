# RPC API

This document describes the RPC API contract that is callable by the frontend application.

API methods are organized by one or more category levels to describe the object type that the
method interacts with.

## Method naming conventions

* Category names are capitalized camel-case
* Category and action names are separated by `::`
* Action names are lower-case camel-case
* Method names start with a Category name and end with an Action name
* Method names must contain only a single Action name
* Method names may contain more than one Category name

## Categories

### RPC

The RPC category is used for general high-level client/server interactions,
primarily internal communication functionality.

### DB

Sub-categories:

* Master - actions pertaining to the master DB.

### User

Actions that apply directly to the user or actions that pertain to objects owned
by the user.  These are typically broken down into additional Sub-categories:

* Shelf
* Collection
* Tag

### Account

Actions that apply to the account or account-level objects.
