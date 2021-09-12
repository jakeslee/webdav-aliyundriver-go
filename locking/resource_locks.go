package locking

import (
	"github.com/sirupsen/logrus"
	"strings"
	"sync"
	"time"
	"webdav-aliyundriver/model"
)

type ResourceLocks struct {
	// 100000
	CleanupLimit int32
	// 0
	CleanupCounter int32
	Locks          map[string]LockedObject
	LocksById      map[string]LockedObject
	TempLocks      map[string]LockedObject
	TempLocksById  map[string]LockedObject
	Root           LockedObject
	TempRoot       LockedObject
	// true
	Temporary bool
}

func Build() ResourceLocks {
	return ResourceLocks{
		Locks:         map[string]LockedObject{},
		LocksById:     map[string]LockedObject{},
		TempLocks:     map[string]LockedObject{},
		TempLocksById: map[string]LockedObject{},
		Temporary:     true,
		CleanupLimit:  100000,
	}
}

func (r ResourceLocks) Lock(transaction model.Transaction, path string, owner string,
	exclusive bool, depth int32, timeout int32, temporary bool) bool {
	var lo LockedObject
	if temporary {
		lo = r.GenerateTempLockedObjects(transaction, path)
		lo.Type = "read"
	} else {
		lo = r.GenerateLockedObjects(transaction, path)
		lo.Type = "write"
	}
	if lo.CheckLocks(exclusive, depth) {

		lo.Exclusive = exclusive
		lo.LockDepth = depth
		lo.ExpireAt = time.Now().Unix() + int64(1000*timeout)
		if lo.Parent != nil {
			lo.Parent.ExpireAt = lo.ExpireAt
			if lo.Parent == &r.Root {
				rootLo := r.LockedObjectByPath(transaction, r.Root.Path)
				rootLo.ExpireAt = lo.ExpireAt
			} else if lo.Parent == &r.TempRoot {
				tempRootLo := r.TempLockedObjectByPath(transaction, r.TempRoot.Path)
				tempRootLo.ExpireAt = lo.ExpireAt
			}
		}
		if lo.addLockedObjectOwner(owner) {
			return true
		} else {
			logrus.Errorf("Couldn't set owner \"" + owner + "\" to resource at '" + path + "'")
			return false
		}
	} else {
		// can not lock
		logrus.Errorf("Lock resource at " + path + " failed because" + "\na parent or child resource is currently locked")
		return false
	}
}

func (r ResourceLocks) Unlock(transaction model.Transaction, id string, owner string) bool {
	if lock, ok := r.LocksById[id]; ok {
		path := lock.Path
		if lo, ok := r.LocksById[id]; ok {
			lo.RemoveLockedObjectOwner(owner)
			if lo.Children == nil && lo.Owner == nil {
				lo.RemoveLockedObject()
			}

		} else {
			// there is no lock at that path. someone tried to unlock it
			// anyway. could point to a problem
			logrus.Trace("net.sf.webdav.locking.ResourceLocks.unlock(): no lock for path " + path)
			return false
		}

		if r.CleanupCounter > r.CleanupLimit {
			r.CleanupCounter = 0
			r.CleanLockedObjects(transaction, r.Root, !r.Temporary)
		}
	}
	r.CheckTimeouts(transaction, !r.Temporary)
}

func (r ResourceLocks) UnlockTemporaryLockedObjects(transaction model.Transaction, path string, owner string) {
	panic("implement me")
}

func (r ResourceLocks) CheckTimeouts(transaction model.Transaction, temporary bool) {
	panic("implement me")
}

func (r ResourceLocks) ExclusiveLock(transaction model.Transaction, path string, owner string, depth, timeout int32) bool {
	panic("implement me")
}

func (r ResourceLocks) SharedLock(transaction model.Transaction, path string, owner string, depth int32, timeout int32) bool {
	panic("implement me")
}

func (r ResourceLocks) LockedObjectByID(transaction model.Transaction, id string) LockedObject {
	panic("implement me")
}

func (r ResourceLocks) LockedObjectByPath(transaction model.Transaction, path string) LockedObject {
	panic("implement me")
}

func (r ResourceLocks) TempLockedObjectByID(transaction model.Transaction, id string) LockedObject {
	panic("implement me")
}

func (r ResourceLocks) TempLockedObjectByPath(transaction model.Transaction, path string) LockedObject {
	panic("implement me")
}

//GenerateTempLockedObjects 为路径及其父资源生成临时LockedObjects
//文件夹。如果LockedObjects已经存在，不创建新的LockedObjects

func (r ResourceLocks) GenerateLockedObjects(transaction model.Transaction, path string) LockedObject {

	if object, ok := r.Locks[path]; ok {
		// there is already a LockedObject on the specified path

		return object

	} else {
		returnObject := CreateLockedObject(r, path, !r.Temporary)
		parentPath := ParentPath(path)
		if len(parentPath) > 0 {
			parentLockedObject := r.GenerateLockedObjects(transaction, parentPath)
			parentLockedObject.Children = append(parentLockedObject.Children, returnObject)
			returnObject.Parent = &parentLockedObject
		}
		return returnObject
	}

}
func (r ResourceLocks) GenerateTempLockedObjects(transaction model.Transaction, path string) LockedObject {
	if lo, ok := r.TempLocks[path]; ok {
		return lo
	} else {
		lockedObject := CreateLockedObject(r, path, r.Temporary)
		parentPath := ParentPath(path)
		if len(parentPath) > 0 {
			parentLockedObject := r.GenerateTempLockedObjects(transaction, parentPath)
			parentLockedObject.Children = append(parentLockedObject.Children, lockedObject)
			lockedObject.Parent = &parentLockedObject
		}
		return lockedObject
	}
}

func (r ResourceLocks) CleanLockedObjects(transaction model.Transaction, lo LockedObject, temporary bool)bool {

	if lo.Children == nil {
		if lo.Owner == nil {
			if temporary {
				lo.RemoveTempLockedObject()
			} else {
				lo.RemoveLockedObject()
			}

			return true
		} else {
			return false
		}
	} else {
		 canDelete := true
		 limit := len(lo.Children)
		for i := 0; i < limit; i++ {
			if !r.CleanLockedObjects(transaction, lo.Children[i], temporary) {
				canDelete = false
			} else {
				// because the deleting shifts the array
				i--
				limit--
			}
		}

		if canDelete {
			if lo.Owner == nil {
				if temporary {
					lo.RemoveTempLockedObject()
				} else {
					lo.RemoveLockedObject()
				}
				return true
			} else {
				return false
			}
		} else {
			return false
		}
	}
}

// ParentPath 通过删除最后一个'/'及其之后的所有内容，从给定路径创建父路径
func ParentPath(path string) string {
	slash := strings.LastIndex(path, "/")
	if slash == -1 {
		return ""
	} else {
		if slash == 0 {
			// return "root" if parent path is empty string
			return "/"
		} else {
			return path[:slash]
		}
	}
}
