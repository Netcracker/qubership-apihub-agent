package service

import (
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/Netcracker/qubership-apihub-agent/view"
	"github.com/shaj13/libcache"
	_ "github.com/shaj13/libcache/lru"
	log "github.com/sirupsen/logrus"
)

type ServiceListCache interface {
	GetServicesList(namespace string, workspaceId string) ServicesListResult
	handleDiscoveryStart(namespace string, workspaceId string, requestedServices []string) (uint64, bool)
	isRunCurrent(namespace string, workspaceId string, runId uint64) bool
	addService(namespace string, workspaceId string, runId uint64, service view.Service) bool
	setResultStatus(namespace string, workspaceId string, runId uint64, status view.StatusEnum, details string)
	clearResultsForNamespace(namespace string, workspaceId string)
}

type ServicesListResult struct {
	Services          []view.Service
	Status            view.StatusEnum
	Details           string
	RequestedServices []string
}

type serviceCacheEntry struct {
	// runId identifies the discovery run that owns this entry. Writes from any
	// other run are dropped, so a superseded run cannot pollute a newer result.
	runId             uint64
	services          []view.Service
	status            view.StatusEnum
	details           string
	requestedServices []string
}

func NewServiceListCache(ttl time.Duration) ServiceListCache {
	cache := libcache.LRU.New(1000)
	cache.SetTTL(ttl)
	cache.RegisterOnExpired(func(key, _ interface{}) {
		cache.Delete(key)
	})
	return &serviceListCacheImpl{cache: cache}
}

type serviceListCacheImpl struct {
	cache      libcache.Cache
	cacheMutex sync.Mutex
	nextRunId  uint64
}

func (s *serviceListCacheImpl) GetServicesList(namespace string, workspaceId string) ServicesListResult {
	s.cacheMutex.Lock()
	defer s.cacheMutex.Unlock()

	id := getNamespaceWithWorkspaceId(namespace, workspaceId)

	val, exists := s.cache.Peek(id)
	if !exists {
		return ServicesListResult{Services: make([]view.Service, 0), Status: view.StatusNone}
	}

	entry := val.(*serviceCacheEntry)
	services := make([]view.Service, len(entry.services))
	copy(services, entry.services)

	return ServicesListResult{
		Services:          services,
		Status:            entry.status,
		Details:           entry.details,
		RequestedServices: entry.requestedServices,
	}
}

func (s *serviceListCacheImpl) handleDiscoveryStart(namespace string, workspaceId string, requestedServices []string) (uint64, bool) {
	s.cacheMutex.Lock()
	defer s.cacheMutex.Unlock()

	id := getNamespaceWithWorkspaceId(namespace, workspaceId)

	if val, exists := s.cache.Peek(id); exists {
		if val.(*serviceCacheEntry).status == view.StatusRunning {
			return 0, false
		}
	}

	s.nextRunId++
	runId := s.nextRunId

	s.cache.Store(id, &serviceCacheEntry{
		runId:             runId,
		services:          []view.Service{},
		status:            view.StatusRunning,
		requestedServices: requestedServices,
	})
	return runId, true
}

func (s *serviceListCacheImpl) clearResultsForNamespace(namespace string, workspaceId string) {
	s.cacheMutex.Lock()
	defer s.cacheMutex.Unlock()

	id := getNamespaceWithWorkspaceId(namespace, workspaceId)

	s.nextRunId++

	s.cache.Store(id, &serviceCacheEntry{
		runId:    s.nextRunId,
		services: []view.Service{},
		status:   view.StatusNone,
	})
}

func (s *serviceListCacheImpl) isRunCurrent(namespace string, workspaceId string, runId uint64) bool {
	s.cacheMutex.Lock()
	defer s.cacheMutex.Unlock()

	val, exists := s.cache.Peek(getNamespaceWithWorkspaceId(namespace, workspaceId))
	return exists && val.(*serviceCacheEntry).runId == runId
}

func (s *serviceListCacheImpl) addService(namespace string, workspaceId string, runId uint64, service view.Service) bool {
	s.cacheMutex.Lock()
	defer s.cacheMutex.Unlock()

	id := getNamespaceWithWorkspaceId(namespace, workspaceId)

	val, exists := s.cache.Peek(id)
	if !exists {
		log.Debugf("Discarding service %s for namespace %s and workspaceId %s: cache entry is gone", service.Id, namespace, workspaceId)
		return false
	}

	entry := val.(*serviceCacheEntry)
	if entry.runId != runId {
		log.Debugf("Discarding service %s for namespace %s and workspaceId %s: discovery run %d was superseded by %d", service.Id, namespace, workspaceId, runId, entry.runId)
		return false
	}

	entry.services = append(entry.services, service)

	sort.Slice(entry.services, func(i, j int) bool {
		return entry.services[i].Name < entry.services[j].Name
	})
	return true
}

func (s *serviceListCacheImpl) setResultStatus(namespace string, workspaceId string, runId uint64, status view.StatusEnum, details string) {
	s.cacheMutex.Lock()
	defer s.cacheMutex.Unlock()

	id := getNamespaceWithWorkspaceId(namespace, workspaceId)

	val, exists := s.cache.Peek(id)
	if !exists {
		log.Warnf("Trying to update missing entry cache status for namespace %s and workspaceId %s", namespace, workspaceId)
		return
	}

	entry := val.(*serviceCacheEntry)
	if entry.runId != runId {
		log.Debugf("Discarding status %s for namespace %s and workspaceId %s: discovery run %d was superseded by %d", status, namespace, workspaceId, runId, entry.runId)
		return
	}
	if entry.status == view.StatusRunning {
		entry.status = status
		entry.details = details
	}
}

const sep = "@||@"

func getNamespaceWithWorkspaceId(namespace string, workspaceId string) string {
	return fmt.Sprintf("%s%s%s", namespace, sep, workspaceId)
}
