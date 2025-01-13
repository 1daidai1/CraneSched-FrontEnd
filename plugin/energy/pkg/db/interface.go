package db

import (
	"github.com/1daidai1/CraneFrontEnd/plugin/energy/pkg/types"
)

type DBInterface interface {
	SaveNodeEnergy(*types.NodeData) error
	SaveTaskEnergy(*types.TaskData) error
	Close() error
}
