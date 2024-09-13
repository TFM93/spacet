package domain

import "context"

type (

	// SpaceXAPI is an interface for fetching spacexapi data
	SpaceXAPIQueries interface {
		GetUpcomingLaunches(ctx context.Context) ([]*Launch, error)
		GetLaunchPads(ctx context.Context) (ret []*LaunchPad, _ error)
	}
)
