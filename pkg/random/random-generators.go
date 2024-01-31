package random

import "github.com/google/uuid"

func GetUUID() uuid.UUID { return uuid.New() }
