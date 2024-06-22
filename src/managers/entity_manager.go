package managers

import "mango/src/logger"

var entityManagerInstance EntityManager

func init() {
	entityManagerInstance = EntityManager{10}
}

type EntityManager struct {
	EntityCount int32
}

func (em *EntityManager) GenerateID() int32 {
	em.EntityCount++
	logger.Debug("EntityID generated: %d", em.EntityCount)
	return em.EntityCount
}

func GetEntityManager() *EntityManager {
	return &entityManagerInstance
}
