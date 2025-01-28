package util

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

// OnePasswordConfig holds the configuration for 1Password Connect
type OnePasswordConfig struct {
	ConnectHost string
	ConnectToken string
	Vault string
}

// SecretRef represents a reference to a 1Password secret
type SecretRef struct {
	ItemID string
	Field  string
}

// FetchOnePasswordSecret retrieves a secret from 1Password Connect and stores it safely
func FetchOnePasswordSecret(config OnePasswordConfig, ref SecretRef, targetPath string) error {
	client := &http.Client{}
	url := fmt.Sprintf("%s/v1/vaults/%s/items/%s", config.ConnectHost, config.Vault, ref.ItemID)
	
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}
	
	req.Header.Add("Authorization", "Bearer "+config.ConnectToken)
	
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error fetching secret: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("error response from 1Password: %s - %s", resp.Status, string(body))
	}
	
	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("error decoding response: %w", err)
	}
	
	// Extract the field value from the response
	fields, ok := result["fields"].([]interface{})
	if !ok {
		return fmt.Errorf("invalid response format: fields not found")
	}
	
	var secretValue string
	for _, field := range fields {
		f, ok := field.(map[string]interface{})
		if !ok {
			continue
		}
		if f["label"] == ref.Field {
			secretValue = f["value"].(string)
			break
		}
	}
	
	if secretValue == "" {
		return fmt.Errorf("field %s not found in item", ref.Field)
	}
	
	// Create directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(targetPath), 0700); err != nil {
		return fmt.Errorf("error creating directory: %w", err)
	}
	
	// Write secret to file with restricted permissions
	if err := os.WriteFile(targetPath, []byte(secretValue), 0600); err != nil {
		return fmt.Errorf("error writing secret to file: %w", err)
	}
	
	return nil
}

// CleanupOnePasswordSecret safely removes the secret file
func CleanupOnePasswordSecret(path string) error {
	err := os.Remove(path)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("error removing secret file: %w", err)
	}
	return nil
}
