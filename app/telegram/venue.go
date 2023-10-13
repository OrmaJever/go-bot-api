package telegram

type Venue struct {
	Location        *Location
	Title           string
	Address         string
	FoursquareId    string
	FoursquareType  string
	GooglePlaceId   string
	GooglePlaceType string
}
