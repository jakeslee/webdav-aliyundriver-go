package req

type Rename struct {
	// refuse
	CheckNameMode string
	DriveId       string
	Name          string
	FileId        string
}
