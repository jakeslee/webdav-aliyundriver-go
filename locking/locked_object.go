package locking

import (
	"webdav-aliyundriver/util"
)

type LockedObject struct {
	ResourceLock *ResourceLocks
	Path         string
	Id           string
	LockDepth    int32
	ExpireAt     int64
	Owner        []string
	Children     []LockedObject
	Parent       *LockedObject
	Exclusive    bool
	Type         string
}

func (o LockedObject) CheckLocks(exclusive bool, depth int32) bool {
	if o.CheckParents(exclusive) && o.CheckChildren(exclusive, depth) {
		return true
	}
	return false
}

//CheckParents 通过删除最后一个'/'和，从给定路径创建父路径
func (o LockedObject) CheckParents(exclusive bool) bool {
	if o.Path == "/" {
		return true
	} else {
		if o.Owner == nil {
			// no owner, checking parents
			return o.Parent != nil && o.Parent.CheckParents(exclusive)
		} else {
			// there already is a owner
			return !(o.Exclusive || exclusive) && o.Parent.CheckParents(exclusive)
		}
	}
}

//CheckChildren helper of checkLocks(). looks if the children are locked
func (o LockedObject) CheckChildren(exclusive bool, depth int32) bool {
	if o.Children == nil {
		// a file

		return o.Owner == nil || !(o.Exclusive || exclusive)
	} else {
		// a folder

		if o.Owner == nil {
			// no owner, checking children
			if depth != 0 {
				canLock := true
				limit := len(o.Children)
				for i := 0; i < limit; i++ {
					if !o.Children[i].CheckChildren(exclusive, depth-1) {
						canLock = false
					}
				}

				return canLock
			} else {
				// depth == 0 -> we don't care for children
				return true
			}
		} else {
			// there already is a owner
			return !(o.Exclusive || exclusive)
		}
	}
}

func (o LockedObject) addLockedObjectOwner(owner string) bool {
	if o.Owner == nil {
		o.Owner = []string{}
	} else {

		size := len(o.Owner)
		var newLockObjectOwner []string

		// check if the owner is already here (that should actually not
		// happen)
		for i := 0; i < size; i++ {
			if o.Owner[i] == owner {
				return false
			}
		}

		//  System.arraycopy(_owner, 0, newLockObjectOwner, 0, size);
		copy(o.Owner, newLockObjectOwner)
		o.Owner = newLockObjectOwner
	}

	o.Owner[len(o.Owner)-1] = owner
	return true
}

func (o LockedObject) RemoveTempLockedObject() {
	if &o != &o.ResourceLock.TempRoot {
		// removing from tree
		if o.Parent != nil && o.Parent.Children != nil {
			size := len(o.Parent.Children)
			for i := 0; i < size; i++ {
				if &o.Parent.Children[i] == &o {

					newChildren := []LockedObject{}
					for i2 := 0; i2 < (size - 1); i2++ {
						if i2 < i {
							newChildren[i2] = o.Parent.Children[i2]
						} else {
							newChildren[i2] = o.Parent.Children[i2+1]
						}
					}
					if len(newChildren) != 0 {
						o.Parent.Children = newChildren
					} else {
						o.Parent.Children = nil
					}
					break
				}
			}

			// removing from hashtable
			delete(o.ResourceLock.TempLocksById, o.Id)
			delete(o.ResourceLock.TempLocks, o.Path)

			// now the garbage collector has some work to do
		}
	}
}

func (o LockedObject) RemoveLockedObjectOwner(owner string) {
	if (this != _resourceLocks._tempRoot) {
		// removing from tree
		if (_parent != null && _parent._children != null) {
			int size = _parent._children.length;
			for (int i = 0; i < size; i++) {
				if (_parent._children[i].equals(this)) {
					LockedObject[] newChildren = new LockedObject[size - 1];
					for (int i2 = 0; i2 < (size - 1); i2++) {
						if (i2 < i) {
							newChildren[i2] = _parent._children[i2];
						} else {
							newChildren[i2] = _parent._children[i2 + 1];
						}
					}
					if (newChildren.length != 0) {
						_parent._children = newChildren;
					} else {
						_parent._children = null;
					}
					break;
				}
			}

			// removing from hashtable
			_resourceLocks._tempLocksByID.remove(getID());
			_resourceLocks._tempLocks.remove(getPath());

			// now the garbage collector has some work to do
		}
	}
}

func (o LockedObject) RemoveLockedObject() {


}

func CreateLockedObject(resLocks ResourceLocks, path string, temporary bool) LockedObject {
	lockedObject := LockedObject{
		Path:         path,
		Id:           util.NextIdStr(),
		ResourceLock: &resLocks,
	}
	if temporary {
		resLocks.Locks[path] = lockedObject
		resLocks.LocksById[lockedObject.Id] = lockedObject
	} else {
		resLocks.TempLocks[path] = lockedObject
		resLocks.TempLocksById[lockedObject.Id] = lockedObject
	}
	return lockedObject
}
