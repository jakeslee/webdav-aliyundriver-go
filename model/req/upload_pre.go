package req

type UploadPre struct {

	//  = "refuse";
	CheckNameMode string
	// ;
	ContentHash string
	//  = "none";
	ContentHashName string
	// ;
	DriveId string
	// ;
	name string
	// ;
	ParentFileId string
	// ;
	ProofCode string
	//  = "v1";
	ProofVersion string
	// ;
	size         int64
	PartInfoList []PartInfo
	//  = "file";
	Type string
}

type PartInfo struct {
	PartNumber int32
	UploadUrl  string
}
