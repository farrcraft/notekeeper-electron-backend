package handler

import (
	"../codes"
	"../notebook"
	messages "../proto"
	"../rpc"

	"github.com/golang/protobuf/proto"
	uuid "github.com/satori/go.uuid"
)

// CreateNotebook is the RPC method to create a new notebook
func CreateNotebook(server *rpc.Server, message []byte) (proto.Message, error) {
	response := &messages.IdResponse{
		Header: rpc.NewResponseHeader(),
	}

	request := messages.CreateNotebookRequest{}
	err := proto.Unmarshal(message, &request)
	if err != nil {
		server.Logger.Debug("Error unmarshaling create notebook request - ", err)
		rpc.SetRPCError(response.Header, codes.ErrorDecode)
		return response, nil
	}

	if !server.IsSignedIn() {
		rpc.SetRPCError(response.Header, codes.ErrorUnauthorized)
		return response, nil
	}

	var scope notebook.Scope
	var container notebook.ContainerType
	if request.Scope == "account" {
		scope = notebook.ScopeAccount
	} else if request.Scope == "user" {
		scope = notebook.ScopeUser
	} else {
		return response, nil
	}

	if request.Container == "collection" {
		container = notebook.ContainerTypeCollection
	} else if request.Container == "shelf" {
		container = notebook.ContainerTypeShelf
	} else {
		return response, nil
	}

	ownerID, err := uuid.FromString(request.OwnerId)
	if err != nil {
		server.Logger.Debug("Invalid notebook owner id - ", err)
		rpc.SetRPCError(response.Header, codes.ErrorDecode)
		return response, nil
	}

	containerID, err := uuid.FromString(request.ContainerId)
	if err != nil {
		server.Logger.Debug("Invalid notebook container id - ", err)
		rpc.SetRPCError(response.Header, codes.ErrorDecode)
		return response, nil
	}

	t := rpc.MessageToTitle(request.Name)
	notebook := notebook.New(t, scope, container, server.DBFactory, server.Logger)
	notebook.OwnerID = ownerID
	notebook.ContainerID = containerID

	err = notebook.Save(server.Account.ActiveUser.PassphraseKey)
	if err != nil {
		rpc.SetInternalError(response.Header, err)
		return response, nil
	}

	response.Id = notebook.ID.String()

	return response, nil
}

// GetNotebooks gets all of the notebooks
func GetNotebooks(server *rpc.Server, message []byte) (proto.Message, error) {
	response := &messages.GetNotebooksResponse{
		Header: rpc.NewResponseHeader(),
	}

	request := messages.GetNotebooksRequest{}
	err := proto.Unmarshal(message, &request)
	if err != nil {
		server.Logger.Debug("Error unmarshaling get notebooks request - ", err)
		rpc.SetRPCError(response.Header, codes.ErrorDecode)
		return response, nil
	}

	if !server.IsSignedIn() {
		rpc.SetRPCError(response.Header, codes.ErrorUnauthorized)
		return response, nil
	}

	var scope notebook.Scope
	var container notebook.ContainerType
	if request.Scope == "account" {
		scope = notebook.ScopeAccount
	} else if request.Scope == "user" {
		scope = notebook.ScopeUser
	} else {
		return response, nil
	}

	if request.Container == "collection" {
		container = notebook.ContainerTypeCollection
	} else if request.Container == "shelf" {
		container = notebook.ContainerTypeShelf
	} else {
		return response, nil
	}

	ownerID, err := uuid.FromString(request.OwnerId)
	if err != nil {
		server.Logger.Debug("Invalid notebook owner id - ", err)
		rpc.SetRPCError(response.Header, codes.ErrorDecode)
		return response, nil
	}

	containerID, err := uuid.FromString(request.ContainerId)
	if err != nil {
		server.Logger.Debug("Invalid notebook container id - ", err)
		rpc.SetRPCError(response.Header, codes.ErrorDecode)
		return response, nil
	}

	// create a new notebook instance to act as a proxy
	nb := notebook.New(nil, scope, container, server.DBFactory, server.Logger)
	nb.OwnerID = ownerID
	nb.ContainerID = containerID

	notebooks, err := nb.LoadAll(server.Account.ActiveUser.PassphraseKey)
	if err != nil {
		rpc.SetInternalError(response.Header, err)
		return response, nil
	}

	for _, n := range notebooks {
		m := &messages.Notebook{
			Id:          n.ID.String(),
			Scope:       request.Scope,
			Container:   request.Container,
			OwnerId:     request.OwnerId,
			ContainerId: request.ContainerId,
			Name:        rpc.TitleToMessage(n.Title),
			Locked:      n.Locked,
			Default:     n.Default,
			Created:     rpc.TimeToMessage(n.Created),
			Updated:     rpc.TimeToMessage(n.Updated),
		}
		response.Notebooks = append(response.Notebooks, m)
	}

	return response, nil
}

// SaveNotebook saves an existing notebook
func SaveNotebook(server *rpc.Server, message []byte) (proto.Message, error) {
	response := &messages.EmptyResponse{
		Header: rpc.NewResponseHeader(),
	}

	request := messages.SaveNotebookRequest{}
	err := proto.Unmarshal(message, &request)
	if err != nil {
		server.Logger.Debug("Error unmarshaling save notebook request - ", err)
		rpc.SetRPCError(response.Header, codes.ErrorDecode)
		return response, nil
	}

	if !server.IsSignedIn() {
		rpc.SetRPCError(response.Header, codes.ErrorUnauthorized)
		return response, nil
	}

	var scope notebook.Scope
	var container notebook.ContainerType
	if request.Scope == "account" {
		scope = notebook.ScopeAccount
	} else if request.Scope == "user" {
		scope = notebook.ScopeUser
	} else {
		return response, nil
	}

	if request.Container == "collection" {
		container = notebook.ContainerTypeCollection
	} else if request.Container == "shelf" {
		container = notebook.ContainerTypeShelf
	} else {
		return response, nil
	}

	id, err := uuid.FromString(request.Id)
	if err != nil {
		server.Logger.Debug("Invalid notebook id - ", err)
		rpc.SetRPCError(response.Header, codes.ErrorDecode)
		return response, nil
	}

	ownerID, err := uuid.FromString(request.OwnerId)
	if err != nil {
		server.Logger.Debug("Invalid notebook owner id - ", err)
		rpc.SetRPCError(response.Header, codes.ErrorDecode)
		return response, nil
	}

	containerID, err := uuid.FromString(request.ContainerId)
	if err != nil {
		server.Logger.Debug("Invalid notebook container id - ", err)
		rpc.SetRPCError(response.Header, codes.ErrorDecode)
		return response, nil
	}

	t := rpc.MessageToTitle(request.Name)
	notebook := notebook.New(t, scope, container, server.DBFactory, server.Logger)
	notebook.ID = id
	notebook.Default = request.Default
	notebook.Locked = request.Locked
	notebook.OwnerID = ownerID
	notebook.ContainerID = containerID

	err = notebook.Save(server.Account.ActiveUser.PassphraseKey)
	if err != nil {
		rpc.SetInternalError(response.Header, err)
		return response, nil
	}
	return response, nil
}

// DeleteNotebook deletes a notebook
func DeleteNotebook(server *rpc.Server, message []byte) (proto.Message, error) {
	response := &messages.EmptyResponse{
		Header: rpc.NewResponseHeader(),
	}

	request := messages.DeleteNotebookRequest{}
	err := proto.Unmarshal(message, &request)
	if err != nil {
		server.Logger.Debug("Error unmarshaling create notebook request - ", err)
		rpc.SetRPCError(response.Header, codes.ErrorDecode)
		return response, nil
	}

	if !server.IsSignedIn() {
		rpc.SetRPCError(response.Header, codes.ErrorUnauthorized)
		return response, nil
	}

	var scope notebook.Scope
	var container notebook.ContainerType
	if request.Scope == "account" {
		scope = notebook.ScopeAccount
	} else if request.Scope == "user" {
		scope = notebook.ScopeUser
	} else {
		return response, nil
	}

	if request.Container == "collection" {
		container = notebook.ContainerTypeCollection
	} else if request.Container == "shelf" {
		container = notebook.ContainerTypeShelf
	} else {
		return response, nil
	}

	id, err := uuid.FromString(request.Id)
	if err != nil {
		server.Logger.Debug("Invalid notebook id - ", err)
		rpc.SetRPCError(response.Header, codes.ErrorDecode)
		return response, nil
	}

	ownerID, err := uuid.FromString(request.OwnerId)
	if err != nil {
		server.Logger.Debug("Invalid notebook owner id - ", err)
		rpc.SetRPCError(response.Header, codes.ErrorDecode)
		return response, nil
	}

	containerID, err := uuid.FromString(request.ContainerId)
	if err != nil {
		server.Logger.Debug("Invalid notebook container id - ", err)
		rpc.SetRPCError(response.Header, codes.ErrorDecode)
		return response, nil
	}

	notebook := notebook.New(nil, scope, container, server.DBFactory, server.Logger)
	notebook.ID = id
	notebook.OwnerID = ownerID
	notebook.ContainerID = containerID

	err = notebook.Delete(server.Account.ActiveUser.PassphraseKey)
	if err != nil {
		rpc.SetInternalError(response.Header, err)
		return response, nil
	}
	return response, nil
}
