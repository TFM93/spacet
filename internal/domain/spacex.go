package domain

import "context"

type (

	// SpaceXAPI is an interface for fetching spacexapi data
	SpaceXAPIQueries interface {
		// GetUpcomingLaunches fetches the upcoming launches from spacex api
		GetUpcomingLaunches(ctx context.Context) ([]*Launch, error)
		// GetLaunchPads fetches the available launchpads from spacex api
		GetLaunchPads(ctx context.Context) (ret []*LaunchPad, _ error)
	}
)
