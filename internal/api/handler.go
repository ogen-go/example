package api

import (
	"context"
	"fmt"

	"example/internal/oas"
)

// Compile-time check for Handler.
var _ oas.Handler = (*Handler)(nil)

type Handler struct {
	oas.UnimplementedHandler // automatically implement all methods
}

func (h Handler) GetPetById(ctx context.Context, params oas.GetPetByIdParams) (oas.GetPetByIdRes, error) {
	return &oas.Pet{
		ID:     oas.NewOptInt64(params.PetId),
		Name:   fmt.Sprintf("Pet %d", params.PetId),
		Status: oas.NewOptPetStatus(oas.PetStatusAvailable),
	}, nil
}
