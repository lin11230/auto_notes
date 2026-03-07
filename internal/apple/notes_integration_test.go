//go:build integration && darwin

package apple

import (
	"os/exec"
	"strings"
	"testing"
	"time"
)

func requireAppleNotesIntegration(t *testing.T) {
	t.Helper()

	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	if _, err := exec.LookPath("osascript"); err != nil {
		t.Skip("Skipping integration test: osascript is not available")
	}

	if _, err := runAppleScript(`tell application "Notes" to count of every folder`); err != nil {
		t.Skipf("Skipping integration test: Apple Notes is not accessible: %v", err)
	}
}

// TestListFoldersIntegration tests listing folders (requires macOS with Notes.app)
func TestListFoldersIntegration(t *testing.T) {
	requireAppleNotesIntegration(t)

	client := NewNotesClient()
	folders, err := client.ListFolders()
	if err != nil {
		t.Fatalf("ListFolders() error: %v", err)
	}

	if len(folders) == 0 {
		t.Log("Warning: No folders found. This might be expected if Notes.app is empty.")
	}

	for i, folder := range folders {
		if folder.ID == "" {
			t.Errorf("Folder %d has empty ID", i)
		}
		if folder.Name == "" {
			t.Errorf("Folder %d has empty Name", i)
		}
	}
}

// TestCreateAndDeleteNoteIntegration tests creating and deleting a note
func TestCreateAndDeleteNoteIntegration(t *testing.T) {
	requireAppleNotesIntegration(t)

	client := NewNotesClient()

	testTitle := "Test Note " + time.Now().Format("20060102150405")
	testBody := "This is a test note created by automated tests."

	note, err := client.CreateNote(testTitle, testBody, "")
	if err != nil {
		t.Fatalf("CreateNote() error: %v", err)
	}

	if note.ID == "" {
		t.Error("Created note has empty ID")
	}
	if note.Name != testTitle {
		t.Errorf("Created note name = %q, want %q", note.Name, testTitle)
	}

	err = client.DeleteNote(note.ID, false)
	if err != nil {
		t.Errorf("DeleteNote() error: %v", err)
	}
}

// TestShowNoteIntegration tests showing a note by ID
func TestShowNoteIntegration(t *testing.T) {
	requireAppleNotesIntegration(t)

	client := NewNotesClient()

	testTitle := "Test Show Note " + time.Now().Format("20060102150405")
	testBody := "Test body for show operation."

	createdNote, err := client.CreateNote(testTitle, testBody, "")
	if err != nil {
		t.Fatalf("CreateNote() error: %v", err)
	}
	defer client.DeleteNote(createdNote.ID, false)

	note, err := client.ShowNote(createdNote.ID)
	if err != nil {
		t.Fatalf("ShowNote() error: %v", err)
	}

	if note.ID != createdNote.ID {
		t.Errorf("ShowNote() ID = %q, want %q", note.ID, createdNote.ID)
	}
	if note.Name != testTitle {
		t.Errorf("ShowNote() Name = %q, want %q", note.Name, testTitle)
	}
	if !strings.Contains(note.Body, testBody) {
		t.Errorf("ShowNote() Body does not contain expected text")
	}
}

// TestSearchNotesIntegration tests searching notes
func TestSearchNotesIntegration(t *testing.T) {
	requireAppleNotesIntegration(t)

	client := NewNotesClient()

	uniqueKeyword := "TESTUNIQUE" + time.Now().Format("20060102150405")
	testTitle := "Test Search " + uniqueKeyword
	testBody := "This note contains the keyword: " + uniqueKeyword

	createdNote, err := client.CreateNote(testTitle, testBody, "")
	if err != nil {
		t.Fatalf("CreateNote() error: %v", err)
	}
	defer client.DeleteNote(createdNote.ID, false)

	notes, err := client.SearchNotes(uniqueKeyword, "")
	if err != nil {
		t.Fatalf("SearchNotes() error: %v", err)
	}

	if len(notes) == 0 {
		t.Error("SearchNotes() returned no results, expected at least 1")
	}

	found := false
	for _, note := range notes {
		if note.ID == createdNote.ID {
			found = true
			break
		}
	}
	if !found {
		t.Error("SearchNotes() did not find the created test note")
	}
}

// TestFindNotesByNameIntegration tests finding notes by name
func TestFindNotesByNameIntegration(t *testing.T) {
	requireAppleNotesIntegration(t)

	client := NewNotesClient()

	uniqueName := "TestFindByName" + time.Now().Format("20060102150405")
	testBody := "Test body for find by name."

	createdNote, err := client.CreateNote(uniqueName, testBody, "")
	if err != nil {
		t.Fatalf("CreateNote() error: %v", err)
	}
	defer client.DeleteNote(createdNote.ID, false)

	notes, err := client.FindNotesByName(uniqueName)
	if err != nil {
		t.Fatalf("FindNotesByName() error: %v", err)
	}

	if len(notes) != 1 {
		t.Errorf("FindNotesByName() returned %d notes, want 1", len(notes))
	}

	if len(notes) > 0 && notes[0].Name != uniqueName {
		t.Errorf("FindNotesByName() returned note with name %q, want %q", notes[0].Name, uniqueName)
	}
}

// TestMoveNoteIntegration tests moving a note between folders
func TestMoveNoteIntegration(t *testing.T) {
	requireAppleNotesIntegration(t)

	client := NewNotesClient()

	folders, err := client.ListFolders()
	if err != nil {
		t.Fatalf("ListFolders() error: %v", err)
	}

	if len(folders) < 2 {
		t.Skip("Need at least 2 folders to test move operation")
	}

	testTitle := "Test Move Note " + time.Now().Format("20060102150405")
	testBody := "This note will be moved."

	createdNote, err := client.CreateNote(testTitle, testBody, "")
	if err != nil {
		t.Fatalf("CreateNote() error: %v", err)
	}
	defer client.DeleteNote(createdNote.ID, false)

	time.Sleep(1 * time.Second)

	note, err := client.ShowNote(createdNote.ID)
	if err != nil {
		t.Fatalf("ShowNote() error before move: %v", err)
	}

	sourceFolder := note.Container
	t.Logf("Note created in folder: %s", sourceFolder)

	var targetFolderName string
	for _, folder := range folders {
		if folder.Name != sourceFolder {
			targetFolderName = folder.Name
			break
		}
	}

	if targetFolderName == "" {
		t.Skip("Could not find a different folder to move to")
	}

	returnedSourceFolder, returnedTargetFolder, err := client.MoveNote(createdNote.ID, targetFolderName)
	if err != nil {
		t.Fatalf("MoveNote() error: %v", err)
	}

	t.Logf("MoveNote returned: source=%s, target=%s", returnedSourceFolder, returnedTargetFolder)

	if returnedTargetFolder != targetFolderName {
		t.Errorf("MoveNote() targetFolder = %q, want %q", returnedTargetFolder, targetFolderName)
	}

	time.Sleep(1 * time.Second)

	note, err = client.ShowNote(createdNote.ID)
	if err != nil {
		t.Fatalf("ShowNote() error after move: %v", err)
	}

	if note.Container != targetFolderName {
		t.Errorf("After move, note.Container = %q, want %q", note.Container, targetFolderName)
	} else {
		t.Logf("Successfully moved note from %s to %s", sourceFolder, targetFolderName)
	}
}
