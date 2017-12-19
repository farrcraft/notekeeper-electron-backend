All data should be encrypted at rest
User controls encryption key

What about an application locked / unlocked state like 1Password?
In memory data is decrypted when application is in unlocked state

How would searching work for encrypted data?
Bleve supports pluggable index data stores as long as they're key/value-based.
A boltdb store already exists. We could just copy that & add the encryption layer.

Shelves can be owned by either an Account a User
All Users in an Account can access content in Account Shelves
User Shelves can only be accessed by the User they belong to
each Account has an Account-level Encryption Key
each User has their own personal Encryption Key
When new Users are created, they receive a copy of the Account-level encryption key
The first User is created at the same time as the Account
The Account Key is created at the same time as the Account & passed down to the next new User

Notes must be organized inside Notebooks
Notebooks can be organized inside Collections
Notebooks must be organized inside Collections or Shelves
Notebooks and Collections must be organized inside Shelves
Accounts can have one or more Shelves
Users can have one or more Shelves
There is a default 'My Shelf' shelf
There is a default 'My Notebook' notebook
Notes, Notebooks, Shelves, and Collections can have Tags
There is a special built-in 'Trash' shelf
Deleted Notes, Notebooks, Collections, Tags default to the Trash shelf before permanent deletion
Notes, Notebooks, Shelves, Collections, Tags all have a Title property
Titles have TitleFormatting attributes
TitleFormatting applies to the entire Title content
TitleFormatting attributes are bold, italics, underscore, strike, background color, font color
Notes have a NoteType
The NoteType applies to the content of a Note
The NoteType affects the interactions available on a Note
NoteTypes include plain text, rich text, markdown, html, renderable image, file, pdf, reminder, list
Markdown and List Notes have an editing mode and a rendered view
List Note content consists of a collection of Lists
A List contains a collection of List Entries
Lists can be nested
A List Entry contains a single line of rich text content
The rendered view of List Notes can filter rendering of checked or unchecked List Entries
Lists support automatically moving checked or unchecked List Entries above or below the other type
Notes can have Revisions
Revisions are copies of Note content saved at a point in time
Revisions cannot be edited
Revisions can be deleted
Revisions can be manually created
Revisions can be automatically created prior to Sync updates
Notes can be created from NoteTemplates
Notes and Notebooks can be locked to prevent accidental editing
Notes and Notebooks can have individual password access and edit unlock codes
Notes and Notebooks can have individual synchronization settings 

Concepts (domain objects):

Account
User
Team
Note
NoteType
Notebook
Collection
Tag
Trash
TitleFormatting
Shelf
Revision
List
ListEntry
ListNote
ImageNote
RichNote
TextNote
NoteTemplate
