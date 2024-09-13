package spacexapi

import (
	"context"
	"time"
)

type Launch struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	DateUTC     time.Time `json:"date_utc"`
	DateUnix    int64     `json:"date_unix"`
	LaunchPadID string    `json:"launchpad"`
	Upcoming    bool      `json:"upcoming"`
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
