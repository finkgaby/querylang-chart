package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"querylang-chart/graph/generated"
	"querylang-chart/graph/model"
	"querylang-chart/server"
)

// Deserialize is the resolver for the deserialize field.
func (r *queryResolver) Deserialize(ctx context.Context, queryLang string) (*model.DeserializedQuery, error) {
	dq := server.DeserializeQuery(queryLang)
	return &dq, nil
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
