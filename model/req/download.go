package req

type Download struct {
	DriveId string
	FileId  string
	// 默认 14400
	expireSec int32
}
