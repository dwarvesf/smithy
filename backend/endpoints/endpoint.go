package endpoints

import (
	"github.com/go-kit/kit/endpoint"

	"github.com/dwarvesf/smithy/backend/service"
)

// Endpoints .
type Endpoints struct {
	AgentSync       endpoint.Endpoint
	AvailableModels endpoint.Endpoint
	AddHook         endpoint.Endpoint
	DBQuery         endpoint.Endpoint
	DBCreate        endpoint.Endpoint
	DBUpdate        endpoint.Endpoint
	DBDelete        endpoint.Endpoint
	ListVersion     endpoint.Endpoint
	RevertVersion   endpoint.Endpoint
	Login           endpoint.Endpoint
	ChangePassword  endpoint.Endpoint
	ViewAdd         endpoint.Endpoint
	ViewList        endpoint.Endpoint
	ViewDelete      endpoint.Endpoint
	ViewExecute     endpoint.Endpoint
<<<<<<< HEAD
	GroupList       endpoint.Endpoint
	GroupCreate     endpoint.Endpoint
	GroupDelete     endpoint.Endpoint
	GroupRead       endpoint.Endpoint
	GroupUpdate     endpoint.Endpoint
=======
	FindAccount     endpoint.Endpoint
	SendEmail       endpoint.Endpoint
	ConfirmCode     endpoint.Endpoint
	ResetPassword   endpoint.Endpoint
>>>>>>> send email & reset password
}

// MakeServerEndpoints returns an Endpoints struct
func MakeServerEndpoints(s service.Service) Endpoints {
	return Endpoints{
		AgentSync:       makeAgentSyncEndpoint(s),
		DBQuery:         makeDBQueryEndpoint(s),
		DBCreate:        makeDBCreateEndpoint(s),
		DBUpdate:        makeDBUpdateEndpoint(s),
		DBDelete:        makeDBDeleteEndpoint(s),
		AvailableModels: makeAvailableModelsEndpoint(s),
		AddHook:         makeAddHookEndpoint(s),
		ListVersion:     makeListVersionEndpoint(s),
		RevertVersion:   makeRevertVersionEndpoint(s),
		Login:           makeLoginEndpoint(s),
		ChangePassword:  makeChangePasswordEndpoint(s),
		ViewAdd:         makeAddViewEndpoint(s),
		ViewList:        makeListViewEndpoint(s),
		ViewDelete:      makeDeleteViewEndpoint(s),
		ViewExecute:     makeExecuteViewEndpoint(s),
<<<<<<< HEAD
		GroupList:       makeListGroupEndpoint(s),
		GroupCreate:     makeCreateGroupEndpoint(s),
		GroupDelete:     makeDeleteGroupEndpoint(s),
		GroupRead:       makeReadGroupEndpoint(s),
		GroupUpdate:     makeUpdateGroupEndpoint(s),
=======
		FindAccount:     makeFindAccountEndpoint(s),
		SendEmail:       makeSendEmailEndpoint(s),
		ConfirmCode:     makeConfirmCodeEndpoint(s),
		ResetPassword:   makeResetPasswordEndpoint(s),
>>>>>>> send email & reset password
	}
}
