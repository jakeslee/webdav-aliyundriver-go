package locking

import "webdav-aliyundriver/model"

type IResourceLocks interface {

	/**
	 * Tries to lock the resource at "path".
	 *
	 * @param transaction
	 * @param path
	 *      what resource to lock
	 * @param owner
	 *      the owner of the lock
	 * @param exclusive
	 *      if the lock should be exclusive (or shared)
	 * @param depth
	 *      depth
	 * @param timeout
	 *      Lock Duration in seconds.
	 * @return true if the resource at path was successfully locked, false if an
	 *  existing lock prevented this
	 * @throws LockFailedException
	 */
	Lock(transaction model.Transaction, path string, owner string,
		exclusive bool, depth int32, timeout int32, temporary bool) bool

	/**
	 * Unlocks all resources at "path" (and all subfolders if existing)<p/> that
	 * have the same owner.
	 *
	 * @param transaction
	 * @param id
	 *      id to the resource to unlock
	 * @param owner
	 *      who wants to unlock
	 */

	Unlock(transaction model.Transaction, id string, owner string) bool

	/**
	 * Unlocks all resources at "path" (and all subfolders if existing)<p/> that
	 * have the same owner.
	 *
	 * @param transaction
	 * @param path
	 *      what resource to unlock
	 * @param owner
	 *      who wants to unlock
	 */

	UnlockTemporaryLockedObjects(transaction model.Transaction, path string, owner string)

	/**
	 * Deletes LockedObjects, where timeout has reached.
	 *
	 * @param transaction
	 * @param temporary
	 *      Check timeout on temporary or real locks
	 */

	CheckTimeouts(transaction model.Transaction, temporary bool)

	/**
	 * Tries to lock the resource at "path" exclusively.
	 *
	 * @param transaction
	 *      Transaction
	 * @param path
	 *      what resource to lock
	 * @param owner
	 *      the owner of the lock
	 * @param depth
	 *      depth
	 * @param timeout
	 *      Lock Duration in seconds.
	 * @return true if the resource at path was successfully locked, false if an
	 *  existing lock prevented this
	 * @throws LockFailedException
	 */

	ExclusiveLock(transaction model.Transaction, path string, owner string,
		depth, timeout int32) bool

	/**
	 * Tries to lock the resource at "path" shared.
	 *
	 * @param transaction
	 *      Transaction
	 * @param path
	 *      what resource to lock
	 * @param owner
	 *      the owner of the lock
	 * @param depth
	 *      depth
	 * @param timeout
	 *      Lock Duration in seconds.
	 * @return true if the resource at path was successfully locked, false if an
	 *  existing lock prevented this
	 * @throws LockFailedException
	 */

	SharedLock(transaction model.Transaction, path string, owner string,
		depth int32, timeout int32) bool

	/**
	 * Gets the LockedObject corresponding to specified id.
	 *
	 * @param transaction
	 * @param id
	 *      LockToken to requested resource
	 * @return LockedObject or null if no LockedObject on specified path exists
	 */

	LockedObjectByID(transaction model.Transaction, id string) LockedObject

	/**
	 * Gets the LockedObject on specified path.
	 *
	 * @param transaction
	 * @param path
	 *      Path to requested resource
	 * @return LockedObject or null if no LockedObject on specified path exists
	 */

	LockedObjectByPath(transaction model.Transaction, path string) LockedObject

	/**
	 * Gets the LockedObject corresponding to specified id (locktoken).
	 *
	 * @param transaction
	 * @param id
	 *      LockToken to requested resource
	 * @return LockedObject or null if no LockedObject on specified path exists
	 */
	TempLockedObjectByID(transaction model.Transaction, id string) LockedObject

	/**
	 * Gets the LockedObject on specified path.
	 *
	 * @param transaction
	 * @param path
	 *      Path to requested resource
	 * @return LockedObject or null if no LockedObject on specified path exists
	 */

	TempLockedObjectByPath(transaction model.Transaction, path string) LockedObject
}
