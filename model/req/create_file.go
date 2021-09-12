package req

type CreateFile struct {
	// 默认 "refuse"
	CheckNameMode string
	DriveId       string
	Name          string
	ParentFileId  string
	Type          string
}

func init() {

}
