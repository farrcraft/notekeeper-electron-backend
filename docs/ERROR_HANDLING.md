# Error Handling

Errors generated from 3rd party libraries or stdlib are referred to here as
*External Errors*.


*External Errors* **MUST**:
* Immediately be logged to the `DEBUG` logger.
* Never be passed along any further after they have been logged.


After an *External Error* has been logged, a new generic error should be
created.  *Generic Errors*:

* **MAY** be passed up to the client application.
* **SHOULD** be logged to the `ERROR` logger at the RPC boundary.
