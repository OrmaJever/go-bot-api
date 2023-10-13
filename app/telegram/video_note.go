package telegram

type VideoNote struct {
	FileId       string     `json:"file_id"`
	FileUniqueId string     `json:"file_unique_id"`
	Length       int16      `json:"length"`
	Duration     int16      `json:"duration"`
	Thumbnail    *PhotoSize `json:"thumbnail"`
	FileSize     int32      `json:"file_size"`
}
