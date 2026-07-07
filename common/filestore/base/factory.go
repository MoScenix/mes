package base

import "github.com/MoScenix/mes/common/filestore"

func NewStore() (Store, error) {
	return NewLocalStore(filestore.GetConf().ShareDir.ShareDir)
}
