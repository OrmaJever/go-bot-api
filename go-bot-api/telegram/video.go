package telegram

type Video struct {
	FileId       string     `json:"file_id"`
	FileUniqueId string     `json:"file_unique_id"`
	Width        int16      `json:"width"`
	Height       int16      `json:"height"`
	Duration     int16      `json:"duration"`
	Thumbnail    *PhotoSize `json:"thumbnail"`
	FileName     string     `json:"file_name"`
	MimeType     string     `json:"mime_type"`
	FileSize     int32      `json:"file_size"`
}
