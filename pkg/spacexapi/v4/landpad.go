package spacexapi

import (
	"context"
)

type LandPad struct {
	ID        string `json:"id"`
	Name      string
	Status    string
	Type      string
	Locality  string
	Region    string
	Latitude  float64
	Longitude float64
}

type LandPadFilters struct {
	Status *string
}

func (l *LandPadFilters) ToQuery() map[string]interface{} {
	query := make(map[string]interface{})
	if l.Status != nil {
		query["status"] = *l.Status
	}
	return query
}

func (c *client) GetAllLandPads(ctx context.Context, filters *LandPadFilters) ([]*LandPad, error) {
	if filters == nil {
		return getAll(ctx, c.GetLandPads, nil)
	}
	return getAll(ctx, c.GetLandPads, filters.ToQuery())
}

func (c *client) GetLandPads(ctx context.Context, filters *FiltersWithPagination) (hasMoreData bool, _ []*LandPad, _ error) {
	endpoint := c.baseURL + "/landpads/query"
	return get[LandPadFilters, *LandPad](c, ctx, endpoint, filters)
}
