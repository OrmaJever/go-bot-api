package telegram

type PhotoSize struct {
	FileId       string `json:"file_id"`
	FileUniqueId string `json:"file_unique_id"`
	Width        int16  `json:"width"`
	Height       int16  `json:"height"`
	FileSize     int32  `json:"file_size"`
}
