package uuid_utils

import "github.com/google/uuid"

func GenerateTrxId() string {
	id := uuid.New()
	return id.String()
}
