package domain

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	UUID      uuid.UUID
	Name      string
	Email     string
	CreatedAt time.Time
}
