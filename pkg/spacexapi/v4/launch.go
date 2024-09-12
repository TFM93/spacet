package spacexapi

import (
	"context"
)

type Launch struct {
	ID string `json:"id"`
}

type LaunchFilters struct {
	Upcoming *bool
}

func (l *LaunchFilters) ToQuery() map[string]interface{} {
	query := make(map[string]interface{})
	if l.Upcoming != nil {
		query["upcoming"] = *l.Upcoming
	}
	return query
}

func (c *client) GetAllLaunches(ctx context.Context, filters *LaunchFilters) ([]*Launch, error) {
	if filters == nil {
		return getAll(ctx, c.GetLaunches, nil)
	}
	return getAll(ctx, c.GetLaunches, filters.ToQuery())
}

func (c *client) GetLaunches(ctx context.Context, filters *FiltersWithPagination) (hasMoreData bool, _ []*Launch, _ error) {
	endpoint := c.baseURL + "/launches/query"
	return get[LaunchFilters, *Launch](c, ctx, endpoint, filters)
}
