package localstorage

import "io/fs"

const (
	OS_READ        fs.FileMode = 04
	OS_WRITE       fs.FileMode = 02
	OS_EX          fs.FileMode = 01
	OS_USER_SHIFT  fs.FileMode = 6
	OS_GROUP_SHIFT fs.FileMode = 3
	OS_OTH_SHIFT   fs.FileMode = 0

	OS_USER_R   fs.FileMode = OS_READ << OS_USER_SHIFT
	OS_USER_W   fs.FileMode = OS_WRITE << OS_USER_SHIFT
	OS_USER_X   fs.FileMode = OS_EX << OS_USER_SHIFT
	OS_USER_RW  fs.FileMode = OS_USER_R | OS_USER_W
	OS_USER_RWX fs.FileMode = OS_USER_RW | OS_USER_X

	OS_GROUP_R   fs.FileMode = OS_READ << OS_GROUP_SHIFT
	OS_GROUP_W   fs.FileMode = OS_WRITE << OS_GROUP_SHIFT
	OS_GROUP_X   fs.FileMode = OS_EX << OS_GROUP_SHIFT
	OS_GROUP_RW  fs.FileMode = OS_GROUP_R | OS_GROUP_W
	OS_GROUP_RX  fs.FileMode = OS_GROUP_R | OS_GROUP_X
	OS_GROUP_RWX fs.FileMode = OS_GROUP_RW | OS_GROUP_X

	OS_OTH_R   fs.FileMode = OS_READ << OS_OTH_SHIFT
	OS_OTH_W   fs.FileMode = OS_WRITE << OS_OTH_SHIFT
	OS_OTH_X   fs.FileMode = OS_EX << OS_OTH_SHIFT
	OS_OTH_RW  fs.FileMode = OS_OTH_R | OS_OTH_W
	OS_OTH_RX  fs.FileMode = OS_OTH_R | OS_OTH_X
	OS_OTH_RWX fs.FileMode = OS_OTH_RW | OS_OTH_X

	OS_ALL_R   fs.FileMode = OS_USER_R | OS_GROUP_R | OS_OTH_R
	OS_ALL_W   fs.FileMode = OS_USER_W | OS_GROUP_W | OS_OTH_W
	OS_ALL_X   fs.FileMode = OS_USER_X | OS_GROUP_X | OS_OTH_X
	OS_ALL_RW  fs.FileMode = OS_ALL_R | OS_ALL_W
	OS_ALL_RX  fs.FileMode = OS_ALL_R | OS_ALL_X
	OS_ALL_RWX fs.FileMode = OS_ALL_RW | OS_GROUP_X
)
