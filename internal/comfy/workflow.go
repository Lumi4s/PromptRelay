package comfy

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func LoadWorkflows() (map[string][]byte, error) {
	log.Println("Loading workflows: started")
	dir := "./data/workflows"

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	if len(entries) == 0 {
		return nil, errors.New("Empty data/workflows catalogue")
	}

	var workflows map[string][]byte = make(map[string][]byte, 3)

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		path := filepath.Join(dir, entry.Name())
		log.Printf("loaded workflow: %v (key is \"%v\")", path, entry.Name())

		data, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}

		workflows[entry.Name()] = data
	}

	log.Printf("Loading workflows: loaded %v workflows", len(workflows))
	return workflows, nil
}

/*
TODO: Workflow builder from prompt & seed & smth else
For now we have hardcoded variant for Krea2.json
..i cant really get how to dynamically edit prompt fields...
*/
func BuildWorkflow(prompt string, rawWorkflow []byte, workflowName string) ([]byte, error) {
	var mapWorkflow map[string]any
	if err := json.Unmarshal(rawWorkflow, &mapWorkflow); err != nil {
		return nil, fmt.Errorf("failed to unmarshal workflow: %w", err)
	}

	switch workflowName {
	case "Krea2.json":
		if err := insertPromptIntoKrea2(prompt, mapWorkflow); err != nil {
			return nil, fmt.Errorf("failed to insert prompt for %s: %w", workflowName, err)
		}
	case "Anima.json":
		if err := insertPromptIntoAnima(prompt, mapWorkflow); err != nil {
			return nil, fmt.Errorf("failed to insert prompt for %s: %w", workflowName, err)
		}
	default:
		return nil, fmt.Errorf("unknown workflow: %s", workflowName)
	}

	updatedBytes, err := json.Marshal(mapWorkflow)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal modified workflow: %w", err)
	}

	return updatedBytes, nil
}

func insertPromptIntoKrea2(prompt string, mapWorkflow map[string]any) error {
	promptField, ok := mapWorkflow["prompt"].(map[string]any)
	if !ok {
		return errors.New("key 'prompt' not found")
	}

	promptNode, ok := promptField["10374:8779"].(map[string]any)
	if !ok {
		return errors.New("key '10374:8779' not found")
	}

	inputs, ok := promptNode["inputs"].(map[string]any)
	if !ok {
		return errors.New("key 'inputs' not found")
	}

	inputs["value"] = prompt

	return nil
}

func insertPromptIntoAnima(prompt string, mapWorkflow map[string]any) error {
	promptField, ok := mapWorkflow["prompt"].(map[string]any)
	if !ok {
		return errors.New("key 'prompt' not found")
	}

	promptNode, ok := promptField["472"].(map[string]any)
	if !ok {
		return errors.New("key '472' not found")
	}

	inputs, ok := promptNode["inputs"].(map[string]any)
	if !ok {
		return errors.New("key 'inputs' not found")
	}

	inputs["value"] = prompt

	return nil
}

func ParseWorkflow(data []byte) error {
	var workflow map[string]any

	err := json.Unmarshal(data, &workflow)
	if err != nil {
		return err
	}

	return nil
}
