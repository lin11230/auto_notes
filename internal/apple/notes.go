package apple

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

// Note represents an Apple Note
type Note struct {
	ID               string
	Name             string
	Body             string
	Container        string
	CreationDate     time.Time
	ModificationDate time.Time
}

// Folder represents an Apple Notes folder
type Folder struct {
	ID   string
	Name string
}

// NotesClient provides methods to interact with Apple Notes via AppleScript
type NotesClient struct{}

var debugMode bool

// NewNotesClient creates a new NotesClient
func NewNotesClient() *NotesClient {
	return &NotesClient{}
}

// SetDebug controls whether low-level AppleScript errors are exposed.
func SetDebug(enabled bool) {
	debugMode = enabled
}

// runAppleScript executes an AppleScript command and returns the output
func runAppleScript(script string) (string, error) {
	cmd := exec.Command("osascript", "-e", script)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		if debugMode {
			return "", fmt.Errorf("AppleScript error: %s", strings.TrimSpace(stderr.String()))
		}
		return "", fmt.Errorf("AppleScript execution failed")
	}

	return strings.TrimSpace(stdout.String()), nil
}

// ListNotes returns all notes or notes from a specific folder
func (c *NotesClient) ListNotes(folder string) ([]Note, error) {
	var script string
	if folder != "" {
		script = fmt.Sprintf(`
			tell application "Notes"
				set output to ""
				repeat with eachNote in notes of folder "%s"
					set noteId to id of eachNote
					set noteName to name of eachNote
					set noteCreated to creation date of eachNote
					set noteModified to modification date of eachNote
					try
						set noteFolder to name of folder of eachNote
					on error
						set noteFolder to "Notes"
					end try
					set output to output & noteId & "|||" & noteName & "|||" & noteFolder & "|||" & (noteCreated as string) & "|||" & (noteModified as string) & "||||"
				end repeat
				return output
			end tell
		`, escapeAppleScriptString(folder))
	} else {
		script = `
			tell application "Notes"
				set output to ""
				repeat with eachNote in notes
					set noteId to id of eachNote
					set noteName to name of eachNote
					set noteCreated to creation date of eachNote
					set noteModified to modification date of eachNote
					try
						set noteFolder to name of folder of eachNote
					on error
						set noteFolder to "Notes"
					end try
					set output to output & noteId & "|||" & noteName & "|||" & noteFolder & "|||" & (noteCreated as string) & "|||" & (noteModified as string) & "||||"
				end repeat
				return output
			end tell
		`
	}

	result, err := runAppleScript(script)
	if err != nil {
		return nil, err
	}

	if result == "" {
		return []Note{}, nil
	}

	notes := strings.Split(result, "||||")
	var noteList []Note

	for _, n := range notes {
		n = strings.TrimSpace(n)
		if n == "" {
			continue
		}
		parts := strings.Split(n, "|||")
		if len(parts) >= 5 {
			note := Note{
				ID:        parts[0],
				Name:      parts[1],
				Container: parts[2],
			}
			// Parse dates
			if t, err := parseAppleDate(parts[3]); err == nil {
				note.CreationDate = t
			}
			if t, err := parseAppleDate(parts[4]); err == nil {
				note.ModificationDate = t
			}
			noteList = append(noteList, note)
		}
	}

	return noteList, nil
}

// ListNotesPaginated returns paginated notes from all notes or a specific folder
func (c *NotesClient) ListNotesPaginated(folder string, limit int, offset int) ([]Note, error) {
	// AppleScript arrays are 1-indexed
	startIndex := offset + 1
	endIndex := offset + limit

	var script string
	if folder != "" {
		script = fmt.Sprintf(`
			tell application "Notes"
				set output to ""
				set totalNotes to count of notes of folder "%s"
				if %d > totalNotes then return ""
				
				set endIndex to %d
				if endIndex > totalNotes then set endIndex to totalNotes
				
				set targetNotes to notes %d thru endIndex of folder "%s"
				
				repeat with eachNote in targetNotes
					set noteId to id of eachNote
					set noteName to name of eachNote
					set noteCreated to creation date of eachNote
					set noteModified to modification date of eachNote
					try
						set noteFolder to name of folder of eachNote
					on error
						set noteFolder to "Notes"
					end try
					set output to output & noteId & "|||" & noteName & "|||" & noteFolder & "|||" & (noteCreated as string) & "|||" & (noteModified as string) & "||||"
				end repeat
				return output
			end tell
		`, escapeAppleScriptString(folder), startIndex, endIndex, startIndex, escapeAppleScriptString(folder))
	} else {
		script = fmt.Sprintf(`
			tell application "Notes"
				set output to ""
				set totalNotes to count of notes
				if %d > totalNotes then return ""
				
				set endIndex to %d
				if endIndex > totalNotes then set endIndex to totalNotes
				
				set targetNotes to notes %d thru endIndex
				
				repeat with eachNote in targetNotes
					set noteId to id of eachNote
					set noteName to name of eachNote
					set noteCreated to creation date of eachNote
					set noteModified to modification date of eachNote
					try
						set noteFolder to name of folder of eachNote
					on error
						set noteFolder to "Notes"
					end try
					set output to output & noteId & "|||" & noteName & "|||" & noteFolder & "|||" & (noteCreated as string) & "|||" & (noteModified as string) & "||||"
				end repeat
				return output
			end tell
		`, startIndex, endIndex, startIndex)
	}

	result, err := runAppleScript(script)
	if err != nil {
		return nil, err
	}

	if result == "" {
		return []Note{}, nil
	}

	notes := strings.Split(result, "||||")
	var noteList []Note

	for _, n := range notes {
		n = strings.TrimSpace(n)
		if n == "" {
			continue
		}
		parts := strings.Split(n, "|||")
		if len(parts) >= 5 {
			note := Note{
				ID:        parts[0],
				Name:      parts[1],
				Container: parts[2],
			}
			// Parse dates
			if t, err := parseAppleDate(parts[3]); err == nil {
				note.CreationDate = t
			}
			if t, err := parseAppleDate(parts[4]); err == nil {
				note.ModificationDate = t
			}
			noteList = append(noteList, note)
		}
	}

	return noteList, nil
}

// CreateNote creates a new note with the given title and body
func (c *NotesClient) CreateNote(title, body, folder string) (*Note, error) {
	// Convert plain text to HTML if needed (wrap in <p> tags)
	htmlBody := textToHTML(body)

	var script string
	if folder != "" {
		script = fmt.Sprintf(`
			tell application "Notes"
				tell folder "%s"
					set newNote to make new note with properties {name:"%s", body:"%s"}
					return id of newNote
				end tell
			end tell
		`, escapeAppleScriptString(folder), escapeAppleScriptString(title), escapeAppleScriptString(htmlBody))
	} else {
		script = fmt.Sprintf(`
			tell application "Notes"
				set newNote to make new note with properties {name:"%s", body:"%s"}
				return id of newNote
			end tell
		`, escapeAppleScriptString(title), escapeAppleScriptString(htmlBody))
	}

	id, err := runAppleScript(script)
	if err != nil {
		return nil, err
	}

	return &Note{
		ID:   id,
		Name: title,
		Body: body,
	}, nil
}

// ShowNote returns the details of a specific note
func (c *NotesClient) ShowNote(identifier string) (*Note, error) {
	// If the identifier looks like a coredata ID, look up by ID directly
	if strings.HasPrefix(identifier, "x-coredata://") {
		script := fmt.Sprintf(`
			tell application "Notes"
				try
					set foundNote to note id "%s"
					set noteId to id of foundNote
					set noteName to name of foundNote
					set noteBody to body of foundNote
					try
						set noteFolder to name of folder of foundNote
					on error
						set noteFolder to "Notes"
					end try
					set noteCreated to creation date of foundNote
					set noteModified to modification date of foundNote
					return noteId & "|||" & noteName & "|||" & noteFolder & "|||" & (noteCreated as string) & "|||" & (noteModified as string) & "|||" & noteBody
				on error
					return "NOT_FOUND"
				end try
			end tell
		`, escapeAppleScriptString(identifier))

		result, err := runAppleScript(script)
		if err != nil {
			return nil, err
		}

		if result == "NOT_FOUND" {
			return nil, fmt.Errorf("note '%s' not found", identifier)
		}

		parts := strings.SplitN(result, "|||", 6)
		if len(parts) < 6 {
			return nil, fmt.Errorf("failed to parse note data")
		}

		note := &Note{
			ID:        parts[0],
			Name:      parts[1],
			Container: parts[2],
			Body:      parts[5],
		}
		if t, err := parseAppleDate(parts[3]); err == nil {
			note.CreationDate = t
		}
		if t, err := parseAppleDate(parts[4]); err == nil {
			note.ModificationDate = t
		}
		return note, nil
	}

	// Otherwise, try to find by name
	script := fmt.Sprintf(`
		tell application "Notes"
			try
				set foundNote to note "%s"
				set noteId to id of foundNote
				set noteName to name of foundNote
				set noteBody to body of foundNote
				try
					set noteFolder to name of folder of foundNote
				on error
					set noteFolder to "Notes"
				end try
				set noteCreated to creation date of foundNote
				set noteModified to modification date of foundNote
				return noteId & "|||" & noteName & "|||" & noteFolder & "|||" & (noteCreated as string) & "|||" & (noteModified as string) & "|||" & noteBody
			on error
				return "NOT_FOUND"
			end try
		end tell
	`, escapeAppleScriptString(identifier))

	result, err := runAppleScript(script)
	if err != nil {
		return nil, err
	}

	if result == "NOT_FOUND" {
		return nil, fmt.Errorf("note '%s' not found", identifier)
	}

	parts := strings.SplitN(result, "|||", 6)
	if len(parts) < 6 {
		return nil, fmt.Errorf("failed to parse note data")
	}

	note := &Note{
		ID:        parts[0],
		Name:      parts[1],
		Container: parts[2],
		Body:      parts[5],
	}
	if t, err := parseAppleDate(parts[3]); err == nil {
		note.CreationDate = t
	}
	if t, err := parseAppleDate(parts[4]); err == nil {
		note.ModificationDate = t
	}
	return note, nil
}

// DeleteNote moves a note to trash or permanently deletes it
func (c *NotesClient) DeleteNote(identifier string, permanent bool) error {
	// Note: AppleScript's "delete" command moves note to trash (Recently Deleted)
	// For permanent deletion, we need to delete from trash separately

	var script string
	if strings.HasPrefix(identifier, "x-coredata://") {
		script = fmt.Sprintf(`
			tell application "Notes"
				try
					set foundNote to note id "%s"
					delete foundNote
					return "DELETED"
				on error
					return "NOT_FOUND"
				end try
			end tell
		`, escapeAppleScriptString(identifier))
	} else {
		script = fmt.Sprintf(`
			tell application "Notes"
				try
					set foundNote to note "%s"
					delete foundNote
					return "DELETED"
				on error
					return "NOT_FOUND"
				end try
			end tell
		`, escapeAppleScriptString(identifier))
	}

	result, err := runAppleScript(script)
	if err != nil {
		return err
	}

	if result == "NOT_FOUND" {
		return fmt.Errorf("note '%s' not found", identifier)
	}

	return nil
}

// SearchNotes searches for notes containing the keyword
func (c *NotesClient) SearchNotes(keyword, folder string) ([]Note, error) {
	var script string
	if folder != "" {
		script = fmt.Sprintf(`
			tell application "Notes"
				set output to ""
				set searchResults to (notes of folder "%s") whose body contains "%s" or name contains "%s"
				repeat with eachNote in searchResults
					set noteId to id of eachNote
					set noteName to name of eachNote
					try
						set noteFolder to name of folder of eachNote
					on error
						set noteFolder to "Notes"
					end try
					set noteCreated to creation date of eachNote
					set noteModified to modification date of eachNote
					set output to output & noteId & "|||" & noteName & "|||" & noteFolder & "|||" & (noteCreated as string) & "|||" & (noteModified as string) & "||||"
				end repeat
				return output
			end tell
		`, escapeAppleScriptString(folder), escapeAppleScriptString(keyword), escapeAppleScriptString(keyword))
	} else {
		script = fmt.Sprintf(`
			tell application "Notes"
				set output to ""
				set searchResults to notes whose body contains "%s" or name contains "%s"
				repeat with eachNote in searchResults
					set noteId to id of eachNote
					set noteName to name of eachNote
					try
						set noteFolder to name of folder of eachNote
					on error
						set noteFolder to "Notes"
					end try
					set noteCreated to creation date of eachNote
					set noteModified to modification date of eachNote
					set output to output & noteId & "|||" & noteName & "|||" & noteFolder & "|||" & (noteCreated as string) & "|||" & (noteModified as string) & "||||"
				end repeat
				return output
			end tell
		`, escapeAppleScriptString(keyword), escapeAppleScriptString(keyword))
	}

	result, err := runAppleScript(script)
	if err != nil {
		return nil, err
	}

	if result == "" {
		return []Note{}, nil
	}

	notes := strings.Split(result, "||||")
	var noteList []Note

	for _, n := range notes {
		n = strings.TrimSpace(n)
		if n == "" {
			continue
		}
		parts := strings.Split(n, "|||")
		if len(parts) >= 5 {
			note := Note{
				ID:        parts[0],
				Name:      parts[1],
				Container: parts[2],
			}
			if t, err := parseAppleDate(parts[3]); err == nil {
				note.CreationDate = t
			}
			if t, err := parseAppleDate(parts[4]); err == nil {
				note.ModificationDate = t
			}
			noteList = append(noteList, note)
		}
	}

	return noteList, nil
}

// ListFolders returns all folders
func (c *NotesClient) ListFolders() ([]Folder, error) {
	script := `
		tell application "Notes"
			set output to ""
			repeat with eachFolder in folders
				set folderId to id of eachFolder
				set folderName to name of eachFolder
				set output to output & folderId & "|||" & folderName & "||||"
			end repeat
			return output
		end tell
	`

	result, err := runAppleScript(script)
	if err != nil {
		return nil, err
	}

	if result == "" {
		return []Folder{}, nil
	}

	folders := strings.Split(result, "||||")
	var folderList []Folder

	for _, f := range folders {
		f = strings.TrimSpace(f)
		if f == "" {
			continue
		}
		parts := strings.Split(f, "|||")
		if len(parts) >= 2 {
			folderList = append(folderList, Folder{
				ID:   parts[0],
				Name: parts[1],
			})
		}
	}

	return folderList, nil
}

// CreateFolder creates a new folder
func (c *NotesClient) CreateFolder(name string) (*Folder, error) {
	script := fmt.Sprintf(`
		tell application "Notes"
			set newFolder to make new folder with properties {name:"%s"}
			return id of newFolder
		end tell
	`, escapeAppleScriptString(name))

	id, err := runAppleScript(script)
	if err != nil {
		return nil, err
	}

	return &Folder{
		ID:   id,
		Name: name,
	}, nil
}

// ExportNote exports a note's body to a string
func (c *NotesClient) ExportNote(identifier string) (string, string, error) {
	note, err := c.ShowNote(identifier)
	if err != nil {
		return "", "", err
	}
	return note.Name, note.Body, nil
}

// FindNotesByName finds all notes with the given name
func (c *NotesClient) FindNotesByName(name string) ([]Note, error) {
	script := fmt.Sprintf(`
		tell application "Notes"
			set output to ""
			set foundNotes to notes whose name is "%s"
			repeat with eachNote in foundNotes
				set noteId to id of eachNote
				set noteName to name of eachNote
				try
					set noteFolder to name of folder of eachNote
				on error
					set noteFolder to "Notes"
				end try
				set noteCreated to creation date of eachNote
				set noteModified to modification date of eachNote
				set output to output & noteId & "|||" & noteName & "|||" & noteFolder & "|||" & (noteCreated as string) & "|||" & (noteModified as string) & "||||"
			end repeat
			return output
		end tell
	`, escapeAppleScriptString(name))

	result, err := runAppleScript(script)
	if err != nil {
		return nil, err
	}

	if result == "" {
		return []Note{}, nil
	}

	notes := strings.Split(result, "||||")
	var noteList []Note

	for _, n := range notes {
		n = strings.TrimSpace(n)
		if n == "" {
			continue
		}
		parts := strings.Split(n, "|||")
		if len(parts) >= 5 {
			note := Note{
				ID:        parts[0],
				Name:      parts[1],
				Container: parts[2],
			}
			if t, err := parseAppleDate(parts[3]); err == nil {
				note.CreationDate = t
			}
			if t, err := parseAppleDate(parts[4]); err == nil {
				note.ModificationDate = t
			}
			noteList = append(noteList, note)
		}
	}

	return noteList, nil
}

// MoveNote moves a note to a target folder
// Returns: source folder name, target folder name, error
func (c *NotesClient) MoveNote(identifier, targetFolder string) (string, string, error) {
	var script string

	// Check if identifier is an ID or name
	if strings.HasPrefix(identifier, "x-coredata://") {
		script = fmt.Sprintf(`
			tell application "Notes"
				try
					set foundNote to note id "%s"
					try
						set sourceFolder to name of folder of foundNote
					on error
						set sourceFolder to "Notes"
					end try
					move foundNote to folder "%s"
					return sourceFolder & "|||" & "%s"
				on error errMsg
					return "ERROR|||" & errMsg
				end try
			end tell
		`, escapeAppleScriptString(identifier), escapeAppleScriptString(targetFolder), escapeAppleScriptString(targetFolder))
	} else {
		script = fmt.Sprintf(`
			tell application "Notes"
				try
					set foundNote to note "%s"
					try
						set sourceFolder to name of folder of foundNote
					on error
						set sourceFolder to "Notes"
					end try
					move foundNote to folder "%s"
					return sourceFolder & "|||" & "%s"
				on error errMsg
					return "ERROR|||" & errMsg
				end try
			end tell
		`, escapeAppleScriptString(identifier), escapeAppleScriptString(targetFolder), escapeAppleScriptString(targetFolder))
	}

	result, err := runAppleScript(script)
	if err != nil {
		return "", "", err
	}

	parts := strings.Split(result, "|||")
	if len(parts) >= 2 {
		if parts[0] == "ERROR" {
			return "", "", fmt.Errorf("%s", parts[1])
		}
		return parts[0], parts[1], nil
	}

	return "", "", fmt.Errorf("unexpected response from AppleScript")
}

// Helper functions

// escapeAppleScriptString escapes special characters for AppleScript
func escapeAppleScriptString(s string) string {
	s = strings.ReplaceAll(s, "\\", "\\\\")
	s = strings.ReplaceAll(s, "\"", "\\\"")
	s = strings.ReplaceAll(s, "'", "\\'")
	return s
}

// parseAppleDate parses an AppleScript date string
func parseAppleDate(s string) (time.Time, error) {
	s = strings.TrimSpace(s)
	// AppleScript date format: "Wednesday, March 4, 2026 at 12:00:00 PM"
	layout := "Monday, January 2, 2006 at 3:04:05 PM"
	t, err := time.Parse(layout, s)
	if err != nil {
		// Try alternate format
		layout2 := "Monday, January 2, 2006 at 15:04:05"
		t, err = time.Parse(layout2, s)
	}
	return t, err
}

// textToHTML converts plain text to HTML format for Apple Notes
func textToHTML(text string) string {
	// Convert newlines to <br> tags for Apple Notes
	text = strings.ReplaceAll(text, "\n", "<br>")
	// Wrap in <p> tags
	text = "<p>" + text + "</p>"
	return text
}
