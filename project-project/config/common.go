package config

import "time"

const (
	MaxLifetime         = time.Hour
	AESKey              = "abcdefgehjhijkmlkjjwwoew"
	DefaultAuthorize    = ""
	DefaultProjectCover = "https://pic.sucaibar.com/pic/201306/13/db1d26c115.jpg"
	CtxTimeOut          = 3 * time.Second
)

// Deleted
const (
	NoDeleted = iota
	Deleted
)

// Archive
const (
	NoArchive = iota
	Archive
)

// AccessControlType
const (
	Open = iota
	Private
	Custom
)

// TaskBoardTheme
const (
	Default = "default"
	Simple  = "simple"
)

//Collected

const (
	NotCollected = iota
	Collected
	Collect  = "collect"
	DCollect = "cancel"
)

// IsOwner
const (
	NoOwner = iota
	Owner
)

// IsExecutor
const (
	NoExecutor = iota
	Executor
)

const (
	NoCanRead = iota
	CanRead
)

// Done 是否完成 0没有完成
const (
	UnDone = iota
	Done
)

// Comment
const (
	NoComment = iota
	Comment
)

const (
	Status   = 1
	NoStatus = 0
)

const (
	Admin           = 0
	Member          = 0
	IsDefaultAdmin  = 0
	IsDefaultMember = 1
)
