package spacex

import (
	"context"
	"spacet/internal/domain"
	"spacet/pkg/logger"
	"spacet/pkg/spacexapi/v4"
)

var _upcomingLaunchesOpt = true
var _upcomingLaunchesFilter = spacexapi.LaunchFilters{Upcoming: &_upcomingLaunchesOpt}
var _upcomingLaunchPadsFilter = spacexapi.LaunchPadFilters{}

type handler struct {
	logger logger.Interface
	client spacexapi.Client
}

// NewSpaceXQueries implements the domain.SpaceXAPIQueries interface
func NewSpaceXQueries(logger logger.Interface) domain.SpaceXAPIQueries {
	ur := &handler{logger: logger, client: spacexapi.New()}
	return ur
}

func (h *handler) GetUpcomingLaunches(ctx context.Context) (ret []*domain.Launch, _ error) {
	launches, err := h.client.GetAllLaunches(ctx, &_upcomingLaunchesFilter)
	if err != nil {
		return nil, err
	}
	// map to domain.Launch
	for _, l := range launches {
		ret = append(ret,
			&domain.Launch{
				ExternalID:  l.ID,
				Domain:      domain.SpaceXDomain,
				Name:        l.Name,
				DateUTC:     l.DateUTC,
				DateUnix:    l.DateUnix,
				LaunchPadID: l.LaunchPadID,
			},
		)
	}
	return ret, nil
}

func (h *handler) GetLaunchPads(ctx context.Context) (ret []*domain.LaunchPad, _ error) {
	launchPads, err := h.client.GetAllLaunchPads(ctx, &_upcomingLaunchPadsFilter)
	if err != nil {
		return nil, err
	}
	// map to domain.LaunchPad
	for _, l := range launchPads {
		ret = append(ret,
			&domain.LaunchPad{
				ID:       l.ID,
				Name:     l.Name,
				Locality: l.Locality,
				Region:   l.Region,
				Timezone: l.Timezone,
				Status:   l.Status,
			},
		)
	}
	return ret, nil
}
