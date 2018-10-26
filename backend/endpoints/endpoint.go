package endpoints

import (
	"github.com/go-kit/kit/endpoint"

	endpointGroup "github.com/dwarvesf/smithy/backend/endpoints/group"
	endpointPermission "github.com/dwarvesf/smithy/backend/endpoints/permission"
	endpointUser "github.com/dwarvesf/smithy/backend/endpoints/user"
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

	FindAccount   endpoint.Endpoint
	SendEmail     endpoint.Endpoint
	ConfirmCode   endpoint.Endpoint
	ResetPassword endpoint.Endpoint

	ViewAdd     endpoint.Endpoint
	ViewList    endpoint.Endpoint
	ViewDelete  endpoint.Endpoint
	ViewExecute endpoint.Endpoint

	GroupCreate  endpoint.Endpoint
	GroupFind    endpoint.Endpoint
	GroupFindAll endpoint.Endpoint
	GroupUpdate  endpoint.Endpoint
	GroupDelete  endpoint.Endpoint

	UserCreate  endpoint.Endpoint
	UserFind    endpoint.Endpoint
	UserFindAll endpoint.Endpoint
	UserUpdate  endpoint.Endpoint
	UserDelete  endpoint.Endpoint

	PermissionFindByGroup endpoint.Endpoint
	PermissionFindByUser  endpoint.Endpoint
	PermissionUpdate      endpoint.Endpoint

	AuthInformation endpoint.Endpoint
	LoginGoogle     endpoint.Endpoint
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
		FindAccount:     makeFindAccountEndpoint(s),
		SendEmail:       makeSendEmailEndpoint(s),
		ConfirmCode:     makeConfirmCodeEndpoint(s),
		ResetPassword:   makeResetPasswordEndpoint(s),

		ViewAdd:     makeAddViewEndpoint(s),
		ViewList:    makeListViewEndpoint(s),
		ViewDelete:  makeDeleteViewEndpoint(s),
		ViewExecute: makeExecuteViewEndpoint(s),

		GroupCreate:  endpointGroup.MakeCreateGroupEndpoint(s),
		GroupFind:    endpointGroup.MakeGroupFindEndpoint(s),
		GroupFindAll: endpointGroup.MakeGroupFindAllEndpoint(s),
		GroupUpdate:  endpointGroup.MakeUpdateGroupEndpoint(s),
		GroupDelete:  endpointGroup.MakeDeleteGroupEndpoint(s),

		UserCreate:  endpointUser.MakeCreateUserEndpoint(s),
		UserFind:    endpointUser.MakeUserFindEndpoint(s),
		UserFindAll: endpointUser.MakeUserFindAllEndpoint(s),
		UserUpdate:  endpointUser.MakeUpdateUserEndpoint(s),
		UserDelete:  endpointUser.MakeDeleteGroupEndpoint(s),

		PermissionFindByGroup: endpointPermission.MakePermissionFindByGroupEndpoint(s),
		PermissionFindByUser:  endpointPermission.MakePermissionFindByUserEndpoint(s),
		PermissionUpdate:      endpointPermission.MakeUpdatePermissionEndpoint(s),

		AuthInformation: makeAuthInformationEndpoint(s),
		LoginGoogle:     makeLoginGoogleEndpoint(s),
	}
}
