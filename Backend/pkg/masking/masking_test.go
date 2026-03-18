package masking

import (
	"encoding/json"
	"testing"
)

func TestMaskSensitiveData(t *testing.T) {
	rawJSON := `{
		"user_id": 123,
		"email": "user@example.com",
		"password": "supersecretpassword",
		"apiKey": "sk_live_123456789",
		"metadata": {
			"session_token": "abc123xyz",
			"items_viewed": ["item1", "item2"]
		},
		"billing": [
			{"cc_number": "1234-5678-9012-3456", "amount": 100},
			{"cc_number": "0000-0000-0000-0000", "amount": 50}
		],
		"notes": "Their email is private@email.com, do not share."
	}`

	var data map[string]interface{}
	err := json.Unmarshal([]byte(rawJSON), &data)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	maskedData := MaskSensitiveData(data)

	// Check shallow password
	if maskedData["password"] != "***MASKED***" {
		t.Errorf("Expected password to be masked, got: %v", maskedData["password"])
	}

	// Check shallow api key (case insensitive checks)
	if maskedData["apiKey"] != "***MASKED***" {
		t.Errorf("Expected apiKey to be masked, got: %v", maskedData["apiKey"])
	}

	// Check deep object masking (token)
	metadata := maskedData["metadata"].(map[string]interface{})
	if metadata["session_token"] != "***MASKED***" {
		t.Errorf("Expected metadata.session_token to be masked, got: %v", metadata["session_token"])
	}

	// Check arrays containing sensitive objects
	billing := maskedData["billing"].([]interface{})
	bill1 := billing[0].(map[string]interface{})
	if bill1["cc_number"] != "***MASKED***" {
		t.Errorf("Expected billing[0].cc_number to be masked, got: %v", bill1["cc_number"])
	}

	// Check loose string masking
	notes := maskedData["notes"].(string)
	if notes != "Their email is ***EMAIL_MASKED***, do not share." {
		t.Errorf("Expected loose string email to be matched and replaced, got: %v", notes)
	}

	// Unmasked properties should remain intact
	if maskedData["user_id"].(float64) != 123 {
		t.Errorf("Expected user_id to remain intact, got: %v", maskedData["user_id"])
	}
	if metadata["items_viewed"].([]interface{})[0].(string) != "item1" {
		t.Errorf("Expected inner arrays to remain intact, got: %v", metadata["items_viewed"])
	}
}

func TestMaskJSONContext(t *testing.T) {
	rawJSON := []byte(`{"secret_key": "topsecret", "public_id": 999}`)

	maskedBytes, err := MaskJSONContext(rawJSON)
	if err != nil {
		t.Fatalf("Failed to mask JSON context: %v", err)
	}

	var result map[string]interface{}
	json.Unmarshal(maskedBytes, &result)

	if result["secret_key"] != "***MASKED***" {
		t.Errorf("Expected secret_key to be masked but got: %v", result["secret_key"])
	}

	if result["public_id"].(float64) != 999 {
		t.Errorf("Expected public_id to be intact but got: %v", result["public_id"])
	}
}

func TestMaskSensitiveData_NilMap(t *testing.T) {
	result := MaskSensitiveData(nil)
	if result != nil {
		t.Errorf("Expected nil when input is nil, got: %v", result)
	}
}
