package domain

type Destination string

func (d Destination) ToString() string {
	return string(d)
}

const (
	DestinationMars         Destination = "MARS"
	DestinationMoon         Destination = "MOON"
	DestinationPluto        Destination = "PLUTO"
	DestinationAsteroidBelt Destination = "ASTEROID_BELT"
	DestinationEuropa       Destination = "EUROPA"
	DestinationTitan        Destination = "TITAN"
	DestinationGanymede     Destination = "GANYMEDE"
)

var validDestinations = []Destination{
	DestinationMars,
	DestinationMoon,
	DestinationPluto,
	DestinationAsteroidBelt,
	DestinationEuropa,
	DestinationTitan,
	DestinationGanymede,
}

// IsValid checks if the destination is valid
func (d Destination) IsValid() bool {
	for _, v := range validDestinations {
		if d == v {
			return true
		}
	}
	return false
}
