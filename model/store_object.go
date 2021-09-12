package model

import "time"

type StoredObject struct {
	IsFolder       bool
	LastModified   time.Time
	CreationDate   time.Time
	ContentLength  int64
	MineType       string
	IsNullResource bool
}
