# Apple Notes CLI

A command-line tool written in Go for managing the macOS Notes app through AppleScript.

## Installation

### Build from source

```bash
# Clone the repository
git clone https://github.com/kclin/auto_notes.git
cd auto_notes

# Build
go build -o notes .

# Move to PATH (optional)
sudo mv notes /usr/local/bin/
```

## Usage

### Help

```bash
./notes --help
./notes <command> --help
```

### List notes

```bash
# List all notes
./notes list
./notes ls

# List notes in a specific folder
./notes list -f "Work"
```

### Create a note

```bash
# Create a new note
./notes create -t "Note title" -b "Note body"

# Create a note in a specific folder
./notes create -t "Note title" -b "Note body" -f "Work"
```

### Show a note

```bash
# Show by note title
./notes show "Meeting Notes"

# Show by note ID
./notes show "x-coredata://..."
```

### Search notes

```bash
# Search all notes
./notes search "keyword"

# Search within a specific folder
./notes search "keyword" -f "Work"
```

### Delete a note

```bash
# Delete by title (moves the note to Recently Deleted)
./notes delete "Meeting Notes"

# Delete by ID
./notes delete "x-coredata://..."
```

### Move notes

```bash
# Move a single note
./notes move "Meeting Notes" -t "Archive"

# Move multiple notes
./notes move "Note 1" "Note 2" "Note 3" -t "Work"

# Move by ID to avoid duplicate-title conflicts
./notes move "x-coredata://..." -t "Personal"
```

### Export notes

```bash
# Export by note ID to avoid duplicate-title conflicts
./notes export "x-coredata://..." --format md -o output.md

# Export as Markdown
./notes export "Meeting Notes" --format md -o output.md

# Export as HTML
./notes export "Meeting Notes" --format html -o output.html

# If no format is specified and output is stdout, Markdown is used by default
./notes export "Meeting Notes"
```

### Folder management

```bash
# List all folders
./notes folder list

# Create a new folder
./notes folder create "New Folder"
```

## Commands

| Command | Description |
|------|------|
| `notes list` | List notes |
| `notes create` | Create a note |
| `notes show` | Show note details |
| `notes search` | Search notes |
| `notes delete` | Delete a note |
| `notes move` | Move notes to another folder |
| `notes export` | Export a note |
| `notes folder list` | List folders |
| `notes folder create` | Create a folder |

## Requirements

- macOS with the Apple Notes app
- Go 1.21+ for building from source

## Testing

The project currently separates tests into two layers:

- Unit tests: do not require Apple Notes and can run in a normal `go test` workflow
- Integration tests: require macOS, `osascript`, and access to Notes.app

```bash
# Run all unit tests
GOCACHE=$(pwd)/.gocache go test ./...

# Run Apple Notes integration tests
GOCACHE=$(pwd)/.gocache go test -tags=integration ./internal/apple
```

Current test coverage includes:

- `escapeAppleScriptString()` escaping behavior
- `parseAppleDate()` Apple date parsing
- `textToHTML()` plain text to HTML conversion
- `NewNotesClient()` client creation
- `ListFolders()` integration test
- `CreateNote()` / `DeleteNote()` integration tests
- `ShowNote()` integration test
- `ExportNote()` integration test by note ID
- `SearchNotes()` integration test
- `FindNotesByName()` integration test
- `MoveNote()` integration test
- `export` format resolution and HTML-to-Markdown conversion unit tests

## Technical Notes

This tool communicates with the macOS Notes app through AppleScript and currently supports:

- Note CRUD operations
- Note moves
- Folder management
- Search
- Export

## License

MIT
