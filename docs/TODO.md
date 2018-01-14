# TODO

[] fix backend compilation errors
[] update frontend boilerplate

[] finish working out encryption key stuff

[?] need a consistent interface for managing encryption keys - retrieving the encrypted encryption key & decrypting it for use for any db
[?] state engine so domain objects can validate whether actions can succeed
    - don't think an FSM actually makes sense.

[] create default notebook during account creation
[] add shutdown rpc handler
[] need to separate note metadata & note content in db so that loading a list of notes gets the metadata & loading a single note gets the content
[?] add FSM for account/user state
[] create shelf
[] save shelf
[] get shelf
[] get list of shelves
[] delete shelf
[] create note
[] get note
[] delete note
[] get list of notes
[] create tag
[] delete tag
[] get list of tags
[] title - tests
[] load notebook
[] notebook - tests
[] get list of notebooks
[] delete notebook
[] encrypt bucket names
[] delete db file when deleting objects
[] Document how error handling & logging will work (in general & in rpc responses)
[] Update general error handling & logging
[] Update rpc error handling & rpc responses
[] fix ui state operations to check that db exists first
[] add copyright headers to source files
[] get notebooks logic
[] close account db after period of inactivity
[] ping rpc keepalive method to keep account db from being closed
