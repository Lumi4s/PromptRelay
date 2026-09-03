package comfy

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"path/filepath"
)

func LoadWorkflows() (map[uint8][]byte, error) {
	log.Println("Loading workflows: started")
	dir := "./data/workflows"

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	if len(entries) == 0 {
		return nil, errors.New("Empty data/workflows catalogue")
	}

	var counter uint8
	var workflows map[uint8][]byte = make(map[uint8][]byte, 3)

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		path := filepath.Join(dir, entry.Name())

		data, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}

		workflows[counter] = data

		counter++
	}

	log.Printf("Loading workflows: loaded %v workflows", counter)
	return workflows, nil
}

/*
TODO: Workflow builder from prompt & seed & smth else
*/
func BuildWorkflow(prompt string, workflow []byte) ([]byte, error) {
	return workflow, nil
}

/*
Parser of workflows
*/
func ParseWorkflow(data []byte) error {
	var workflow map[string]any

	err := json.Unmarshal(data, &workflow)
	if err != nil {
		return err
	}

	return nil
}
