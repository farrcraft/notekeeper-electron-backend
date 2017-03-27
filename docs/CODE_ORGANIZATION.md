# Code Organization

## domain

The `domain` module is a layer of glue between core domain types and the RPC
layer.  Any time the RPC layer might need to perform some logic that isn't
immediately related to building the response or that might cause a lot of
additional exit points for the RPC call, that logic should be put into the
domain module.

This is a lousy name for the module. Other names like `logic` or `api` might
be at least as appropriate if not better.

The goal is to keep as much logic out of the RPC layer and keep it as thin as
possible.  The domain or api layer takes on all of that logic and gives us an
independent surface to test against without having to involve RPC in the test
harness.
