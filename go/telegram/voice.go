package telegram

type Voice struct {
	FileId       string `json:"file_id"`
	FileUniqueId string `json:"file_unique_id"`
	Duration     int16  `json:"duration"`
	MimeType     string `json:"mime_type"`
	FileSize     int32  `json:"file_size"`
}
