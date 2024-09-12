package spacexapi

import (
	"context"
	"net/http"
	"time"
)

const (
	_defaultBaseURL                = "https://api.spacexdata.com/v4"
	_httpClientMaxIdleConns        = 6
	_httpClientMaxConnsPerHost     = 20
	_httpClientMaxIdleConnsPerHost = 2
	_httpClientDefaultTimeout      = 10 * time.Second
)

type Client interface {
	GetAllLandPads(ctx context.Context, filters *LandPadFilters) ([]*LandPad, error)
	GetAllLaunchPads(ctx context.Context, filters *LaunchPadFilters) ([]*LaunchPad, error)
	GetAllLaunches(ctx context.Context, filters *LaunchFilters) ([]*Launch, error)
	GetLandPads(ctx context.Context, filters *FiltersWithPagination) (hasMoreData bool, _ []*LandPad, _ error)
	GetLaunchPads(ctx context.Context, filters *FiltersWithPagination) (hasMoreData bool, _ []*LaunchPad, _ error)
	GetLaunches(ctx context.Context, filters *FiltersWithPagination) (hasMoreData bool, _ []*Launch, _ error)
}

type client struct {
	baseURL string
	http    *http.Client
}

// New creates a new spacex client
func New() Client {
	t := http.DefaultTransport.(*http.Transport).Clone()
	t.MaxIdleConns = _httpClientMaxIdleConns
	t.MaxIdleConnsPerHost = _httpClientMaxIdleConnsPerHost
	t.MaxConnsPerHost = _httpClientMaxConnsPerHost

	return &client{
		baseURL: _defaultBaseURL,
		http: &http.Client{
			Timeout:   _httpClientDefaultTimeout,
			Transport: t,
		}}
}
