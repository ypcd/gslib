package gid

import "sync/atomic"

type IGId interface {
	GenerateId() uint64
	GetId() uint64
}

type gidImp struct {
	Id atomic.Uint64
}

func NewGIdImp() IGId {
	return &gidImp{}
}

func (g *gidImp) GenerateId() uint64 {
	return g.Id.Add(1)
}

func (g *gidImp) GetId() uint64 {
	return g.Id.Load()
}
