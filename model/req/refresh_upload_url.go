package req

type RefreshUploadUrl struct {
	DriveId      string
	PartInfoList []PartInfo
	FileId       string
	UploadId     string
}
