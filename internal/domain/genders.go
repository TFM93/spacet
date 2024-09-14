package domain

type Gender string

const (
	GenderUnspecified Gender = "UNSPECIFIED"
	GenderMale        Gender = "MALE"
	GenderFemale      Gender = "FEMALE"
	GenderNonBinary   Gender = "NON_BINARY"
	GenderOther       Gender = "OTHER"
)
