package redact

import "strings"

var sensitiveKeys = []string{
	"PASSWORD",
	"PASS",
	"SECRET",
	"TOKEN",
	"API_KEY",
	"APIKEY",
	"ACCESS_KEY",
	"PRIVATE_KEY",
	"CREDENTIAL",
	"AUTH",
	"DATABASE_URL",
	"CONNECTION_STRING",
}

func Value(key, value string) string {
	upperKey := strings.ToUpper(key)

	for _, sensitive := range sensitiveKeys {
		if strings.Contains(upperKey, sensitive) {
			return "<redacted>"
		}
	}

	return value
}
