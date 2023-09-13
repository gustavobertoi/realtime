package uuid

import (
	"regexp"

	"github.com/google/uuid"
)

func NewUUID() string {
	return uuid.New().String()
}

func IsValidUUID(uuid string) bool {
	uuidPattern := `^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`
	regex := regexp.MustCompile(uuidPattern)
	return regex.MatchString(uuid)
}
