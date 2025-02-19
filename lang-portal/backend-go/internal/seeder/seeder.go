package seeder

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type WordData struct {
	German  string     `json:"german"`
	English string     `json:"english"`
	Parts   WordParts `json:"parts"`
}

type WordParts struct {
	Article string `json:"article"`
	Plural  string `json:"plural"`
}

type WordsFile struct {
	Words []WordData `json:"words"`
}

type GroupData struct {
	Name  string   `json:"name"`
	Words []string `json:"words"`
}

type GroupsFile struct {
	Groups []GroupData `json:"groups"`
}

func LoadSeedData(db *sql.DB, seedDir string) error {
	// Load and insert words
	wordsFile := filepath.Join(seedDir, "words.json")
	words, err := loadWordsFromJSON(wordsFile)
	if err != nil {
		return fmt.Errorf("failed to load words: %w", err)
	}

	wordIDs, err := insertWords(db, words)
	if err != nil {
		return fmt.Errorf("failed to insert words: %w", err)
	}

	// Load and insert groups
	groupsFile := filepath.Join(seedDir, "groups.json")
	groups, err := loadGroupsFromJSON(groupsFile)
	if err != nil {
		return fmt.Errorf("failed to load groups: %w", err)
	}

	if err := insertGroups(db, groups, words, wordIDs); err != nil {
		return fmt.Errorf("failed to insert groups: %w", err)
	}

	return nil
}

func loadWordsFromJSON(filename string) ([]WordData, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var wordsFile WordsFile
	if err := json.Unmarshal(data, &wordsFile); err != nil {
		return nil, err
	}

	return wordsFile.Words, nil
}

func loadGroupsFromJSON(filename string) ([]GroupData, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var groupsFile GroupsFile
	if err := json.Unmarshal(data, &groupsFile); err != nil {
		return nil, err
	}

	return groupsFile.Groups, nil
}

func insertWords(db *sql.DB, words []WordData) (map[string]int64, error) {
	wordIDs := make(map[string]int64)
	
	for _, word := range words {
		partsJSON, err := json.Marshal(word.Parts)
		if err != nil {
			return nil, err
		}

		result, err := db.Exec(
			"INSERT INTO words (german, english, parts) VALUES (?, ?, ?)",
			word.German,
			word.English,
			string(partsJSON),
		)
		if err != nil {
			return nil, err
		}

		id, err := result.LastInsertId()
		if err != nil {
			return nil, err
		}

		wordIDs[word.German] = id
	}

	return wordIDs, nil
}

func insertGroups(db *sql.DB, groups []GroupData, words []WordData, wordIDs map[string]int64) error {
	for _, group := range groups {
		// Insert group
		result, err := db.Exec("INSERT INTO groups (name) VALUES (?)", group.Name)
		if err != nil {
			return err
		}

		groupID, err := result.LastInsertId()
		if err != nil {
			return err
		}

		// Insert word-group relationships
		for _, wordGerman := range group.Words {
			wordID, ok := wordIDs[wordGerman]
			if !ok {
				return fmt.Errorf("word %q not found in words list", wordGerman)
			}

			_, err = db.Exec(
				"INSERT INTO words_groups (word_id, group_id) VALUES (?, ?)",
				wordID,
				groupID,
			)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
