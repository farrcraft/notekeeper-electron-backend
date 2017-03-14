# RPC Error Handling

* Never return an `error` type in a gRPC handler
* If a gRPC handler can return an error state, include the correct status response in the return payload
* Use the standard `codes.Error*` and `codes.Status*` constants
* Any time a new error is created (`errors.New(...)`), it should be logged immediately before it gets `return`ed to the caller
