package managers

import (
	"mango/src/logger"
	"sync/atomic"
)

var entityManagerInstance EntityManager

func init() {
	entityManagerInstance = EntityManager{10}
}

type EntityManager struct {
	EntityCount int32
}

func (em *EntityManager) GenerateID() int32 {
	atomic.AddInt32(&em.EntityCount, 1)
	logger.Debug("EntityID generated: %d", em.EntityCount)
	return em.EntityCount
}

func GetEntityManager() *EntityManager {
	return &entityManagerInstance
}
