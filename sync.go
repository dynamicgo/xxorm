package xxorm

import (
	"sync"

	"github.com/dynamicgo/slf4go"
	"github.com/go-xorm/xorm"
)

// SyncHandler sync handler prototype
type SyncHandler func() []interface{}

type syncRegisterImpl struct {
	slf4go.Logger
	sync.RWMutex
	handlers []SyncHandler
}

func (register *syncRegisterImpl) Register(handler SyncHandler) {
	register.Lock()
	defer register.Unlock()

	register.handlers = append(register.handlers, handler)
}

func (register *syncRegisterImpl) Sync(engine *xorm.Engine) error {
	register.RLock()
	defer register.RUnlock()

	var tables []interface{}

	for _, handler := range register.handlers {
		tables = append(tables, handler()...)
	}

	return engine.Sync2(tables...)
}

var register = &syncRegisterImpl{
	Logger: slf4go.Get("orm"),
}

// Register .
func Register(handler SyncHandler) {
	register.Register(handler)
}

// RegisterWithName .
func RegisterWithName(name string, handler SyncHandler) {
	register.DebugF("register orm module %s", name)
	register.Register(handler)
}

// Sync .
func Sync(engine *xorm.Engine) error {
	return register.Sync(engine)
}
