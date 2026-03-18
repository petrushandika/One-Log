package masking

import (
	"encoding/json"
	"regexp"
	"strings"
)

// List of commonly known sensitive JSON keys (case-insensitive checks)
var sensitiveKeys = []string{
	"password", "token", "secret", "credit_card", "cc_number", "cvv", "api_key", "apikey",
	"auth_token", "access_token", "refresh_token", "ssn", "social_security",
}

// Regex to capture potential email addresses for masking
var emailRegex = regexp.MustCompile(`[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`)

// MaskSensitiveData takes a generic unmarshalled map, masks sensitive PII values,
// and returns the masked map. This modifies the map.
func MaskSensitiveData(data map[string]interface{}) map[string]interface{} {
	if data == nil {
		return nil
	}
	return maskInternal(data).(map[string]interface{})
}

func maskInternal(node interface{}) interface{} {
	switch v := node.(type) {
	case map[string]interface{}:
		for key, val := range v {
			if isSensitiveKey(key) {
				v[key] = "***MASKED***"
			} else {
				v[key] = maskInternal(val)
			}
		}
		return v
	case []interface{}:
		for i, val := range v {
			v[i] = maskInternal(val)
		}
		return v
	case string:
		// Also mask loose emails in strings
		if emailRegex.MatchString(v) {
			return emailRegex.ReplaceAllString(v, "***EMAIL_MASKED***")
		}
		return v
	default:
		return v
	}
}

func isSensitiveKey(key string) bool {
	lowerKey := strings.ToLower(key)
	for _, sensitive := range sensitiveKeys {
		if strings.Contains(lowerKey, sensitive) {
			return true
		}
	}
	return false
}

// MaskJSONContext helps unmarshal, mask, and re-marshal db-bound JSON byte arrays
func MaskJSONContext(rawJSON []byte) ([]byte, error) {
	if len(rawJSON) == 0 {
		return rawJSON, nil
	}

	var data map[string]interface{}
	if err := json.Unmarshal(rawJSON, &data); err != nil {
		// Just return original if it's not a valid map
		return rawJSON, nil
	}

	maskedData := MaskSensitiveData(data)
	
	maskedJSON, err := json.Marshal(maskedData)
	if err != nil {
		return rawJSON, nil
	}

	return maskedJSON, nil
}
