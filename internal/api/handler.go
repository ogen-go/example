package api

import (
	"context"
	"fmt"

	"github.com/go-faster/sdk/zctx"
	"go.uber.org/zap"

	"example/internal/oas"
)

// Compile-time check for Handler.
var _ oas.Handler = (*Handler)(nil)

type Handler struct {
	oas.UnimplementedHandler // automatically implement all methods
}

func (h Handler) GetPetById(ctx context.Context, params oas.GetPetByIdParams) (oas.GetPetByIdRes, error) {
	zctx.From(ctx).Info("GetPetById", zap.Any("params", params))
	return &oas.Pet{
		ID:     oas.NewOptInt64(params.PetId),
		Name:   fmt.Sprintf("Pet %d", params.PetId),
		Status: oas.NewOptPetStatus(oas.PetStatusAvailable),
	}, nil
}
