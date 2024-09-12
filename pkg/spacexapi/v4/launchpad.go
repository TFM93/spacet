package spacexapi

import (
	"context"
)

type LaunchPad struct {
	ID string `json:"id"`
}

type LaunchPadFilters struct {
	Status *string
}

func (l *LaunchPadFilters) ToQuery() map[string]interface{} {
	query := make(map[string]interface{})
	if l.Status != nil {
		query["status"] = *l.Status
	}
	return query
}

func (c *client) GetAllLaunchPads(ctx context.Context, filters *LaunchPadFilters) ([]*LaunchPad, error) {
	if filters == nil {
		return getAll(ctx, c.GetLaunchPads, nil)
	}
	return getAll(ctx, c.GetLaunchPads, filters.ToQuery())
}

func (c *client) GetLaunchPads(ctx context.Context, filters *FiltersWithPagination) (hasMoreData bool, _ []*LaunchPad, _ error) {
	endpoint := c.baseURL + "/landpads/query"
	return get[LaunchPadFilters, *LaunchPad](c, ctx, endpoint, filters)
}
