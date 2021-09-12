package res

import "webdav-aliyundriver/model/req"

type UploadPre struct {
	FileId       string
	FileName     string
	Location     string
	RapidUpload  bool
	Type         string
	UploadId     string
	PartInfoList []req.PartInfo
}
