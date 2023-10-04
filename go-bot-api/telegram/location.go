package telegram

type Location struct {
	Longitude            float64 `json:"longitude"`
	Latitude             float64 `json:"latitude"`
	HorizontalAccuracy   float64 `json:"horizontal_accuracy"`
	LivePeriod           int64   `json:"live_period"`
	Heading              int16   `json:"heading"`
	ProximityAlertRadius int16   `json:"proximity_alert_radius"`
}
