package controller

import (
	"net/http"

	"github.com/Netcracker/qubership-apihub-agent/exception"
	"github.com/Netcracker/qubership-apihub-agent/secctx"
	"github.com/Netcracker/qubership-apihub-agent/service"
	"github.com/Netcracker/qubership-apihub-agent/view"
	log "github.com/sirupsen/logrus"
)

type ServiceController interface {
	ListServices(w http.ResponseWriter, r *http.Request)
	StartDiscovery(w http.ResponseWriter, r *http.Request)
	ListServiceNames(w http.ResponseWriter, r *http.Request)
	ListServiceItems(w http.ResponseWriter, r *http.Request)
}

func NewServiceController(serviceListCache service.ServiceListCache,
	discoveryService service.DiscoveryService,
	listNamesService service.ListService) ServiceController {
	return serviceControllerImpl{
		serviceListCache: serviceListCache,
		discoveryService: discoveryService,
		listService:      listNamesService,
	}
}

type serviceControllerImpl struct {
	serviceListCache service.ServiceListCache
	discoveryService service.DiscoveryService
	listService      service.ListService
}

func (s serviceControllerImpl) ListServices(w http.ResponseWriter, r *http.Request) {
	namespace := getStringParam(r, "name")
	workspaceId := getStringParam(r, "workspaceId")
	//v1 support
	if workspaceId == "" {
		workspaceId = view.DefaultWorkspaceId
	}
	services, status, details := s.serviceListCache.GetServicesList(namespace, workspaceId)
	respondWithJson(w, http.StatusOK, view.ServiceListResponse{Services: services, Status: status, Debug: details})
}

func (s serviceControllerImpl) StartDiscovery(w http.ResponseWriter, r *http.Request) {
	namespace := getStringParam(r, "name")
	workspaceId := getStringParam(r, "workspaceId")
	//v1 support
	if workspaceId == "" {
		workspaceId = view.DefaultWorkspaceId
	}
	err := s.discoveryService.StartDiscovery(secctx.Create(r), namespace, workspaceId)
	if err != nil {
		log.Error("Failed to start discovery process: ", err.Error())
		if customError, ok := err.(*exception.CustomError); ok {
			RespondWithCustomError(w, customError)
		} else {
			RespondWithCustomError(w, &exception.CustomError{
				Status:  http.StatusInternalServerError,
				Message: "Failed to start discovery process",
				Debug:   err.Error()})
		}
		return
	}
	w.WriteHeader(http.StatusAccepted)
}

func (s serviceControllerImpl) ListServiceNames(w http.ResponseWriter, r *http.Request) {
	namespace := getStringParam(r, "name")

	result, err := s.listService.ListServiceNames(namespace)
	if err != nil {
		log.Error("Failed to list service names: ", err.Error())
		if customError, ok := err.(*exception.CustomError); ok {
			RespondWithCustomError(w, customError)
		} else {
			RespondWithCustomError(w, &exception.CustomError{
				Status:  http.StatusInternalServerError,
				Message: "Failed to list service names",
				Debug:   err.Error()})
		}
		return
	}
	respondWithJson(w, http.StatusOK, view.ServiceNamesResponse{ServiceNames: result})
}

func (s serviceControllerImpl) ListServiceItems(w http.ResponseWriter, r *http.Request) {
	namespace := getStringParam(r, "name")

	result, err := s.listService.ListServiceItems(namespace)
	if err != nil {
		log.Error("Failed to list service items: ", err.Error())
		if customError, ok := err.(*exception.CustomError); ok {
			RespondWithCustomError(w, customError)
		} else {
			RespondWithCustomError(w, &exception.CustomError{
				Status:  http.StatusInternalServerError,
				Message: "Failed to list service items",
				Debug:   err.Error()})
		}
		return
	}
	respondWithJson(w, http.StatusOK, view.ServiceItemsResponse{ServiceItems: result})
}
