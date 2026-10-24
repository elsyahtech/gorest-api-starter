package model

import (
	"fmt"
	"strings"
)

// convertPH returns the bind marker used by the configured SQL
// driver. The fallback marker matches SQLite and MySQL.
func convertPH(driver string, index int) string {
	switch driver {
	case "oracle":
		return fmt.Sprintf(":%d", index)
	case "sqlserver":
		return fmt.Sprintf("@p%d", index)
	case "postgres":
		return fmt.Sprintf("$%d", index)
	default:
		return "?"
	}
}

func ph(driver string, count int) string {
	placeholders := make([]string, count)

	for i := range placeholders {
		placeholders[i] = convertPH(driver, i+1)
	}

	return strings.Join(placeholders, ", ")
}

func convertBooleanToInt(value bool) int {
	if value {
		return 1
	}

	return 0
}
