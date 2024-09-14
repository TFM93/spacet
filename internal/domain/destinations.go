package domain

type Destination string

func (d Destination) ToString() string {
	return string(d)
}

const (
	DestinationMars         Destination = "Mars"
	DestinationMoon         Destination = "Moon"
	DestinationPluto        Destination = "Pluto"
	DestinationAsteroidBelt Destination = "Asteroid Belt"
	DestinationEuropa       Destination = "Europa"
	DestinationTitan        Destination = "Titan"
	DestinationGanymede     Destination = "Ganymede"
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
