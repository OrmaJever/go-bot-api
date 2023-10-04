package telegram

type Document struct {
	FileId       string     `json:"file_id"`
	FileUniqueId string     `json:"file_unique_id"`
	Thumbnail    *PhotoSize `json:"thumbnail"`
	FileName     string     `json:"file_name"`
	MimeType     string     `json:"mime_type"`
	FileSize     int32      `json:"file_size"`
}
