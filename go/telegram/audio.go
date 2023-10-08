package telegram

type Audio struct {
	FileId       string     `json:"file_id"`
	FileUniqueId string     `json:"file_unique_id"`
	Duration     int16      `json:"duration"`
	Performer    string     `json:"performer"`
	Title        string     `json:"title"`
	FileName     string     `json:"file_name"`
	MimeType     string     `json:"mime_type"`
	FileSize     int32      `json:"file_size"`
	Thumbnail    *PhotoSize `json:"thumbnail"`
}
