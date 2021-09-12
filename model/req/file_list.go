package req

type FileList struct {
	// ;
	DriveId string
	//  = false;
	all bool
	//  = "*";
	fields string
	//  = "image/resize,w_400/format,jpeg";
	ImageThumbnailProcess string
	//  = "image/resize,w_1920/format,jpeg";
	ImageUrlProcess string
	// ;
	ParentFileId string
	//  = "video/snapshot,t_0,f_jpg,ar_auto,w_300";
	VideoThumbnailProcess string
}
