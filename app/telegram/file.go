package telegram

type File struct {
	FileId       string `json:"file_id"`
	FileUniqueId string `json:"file_unique_id"`
	FileSize     int32  `json:"file_size"`
	FilePath     string `json:"file_path"`
}
