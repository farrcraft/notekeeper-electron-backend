# TODO

[] add shutdown rpc handler
[] need to separate note metadata & note content in db so that loading a list of notes gets the metadata & loading a single note gets the content
[] add FSM for account/user state
[] do we need index & map types to match the corresponding db buckets?
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
[] create default notebook during account creation
[] delete db file when deleting objects
[] Document how error handling & logging will work (in general & in rpc responses)
[] Update general error handling & logging
[] Update rpc error handling & rpc responses
[] fix ui state operations to check that db exists first
[] add copyright headers to source files
[] get notebooks logic
[] close account db after period of inactivity
[] ping rpc keepalive method to keep account db from being closed



# DONE

[+] rpc for shelves
[+] rpc for collections
[+] rpc - create notebook
[+] rpc - get notebooks
[+] rpc - get notebook notes
[+] remove service library (backend won't be a native service)
[+] add ready rpc handler
[+] add note handlers
[+] add tag handlers
[+] add notebook handlers
[+] generate code coverage reports
[+] move rpc handlers out of base rpc module & into their own module
[+] make sure all rpc error paths set the error scope
[+] refactor error codes w/ scope property
[+] work out encryption key layers
[+] add encryption keys in db module
[+] refactor db code to use db module
[+] move account data into account db instead of master db & just keep account_map in master db
