package handler

import (
	"notekeeper-electron-backend/codes"
	"notekeeper-electron-backend/notebook"
	messages "notekeeper-electron-backend/proto"
	"notekeeper-electron-backend/rpc"

	"github.com/golang/protobuf/proto"
	uuid "github.com/satori/go.uuid"
)

func createNotebook(server *rpc.Server, message []byte, scope string, context *rpc.RequestContext) (proto.Message, error) {
	response := &messages.IdResponse{
		Header: rpc.NewResponseHeader(),
	}

	request := messages.CreateNotebookRequest{}
	err := proto.Unmarshal(message, &request)
	if err != nil {
		server.Logger.Warn("Error unmarshaling create notebook request - ", err)
		rpc.SetRPCError(response.Header, codes.ErrorDecode)
		return response, nil
	}

	if !server.IsSignedIn() {
		rpc.SetRPCError(response.Header, codes.ErrorUnauthorized)
		return response, nil
	}

	var notebookScope notebook.Scope
	var container notebook.ContainerType
	if scope == "account" {
		notebookScope = notebook.ScopeAccount
	} else if scope == "user" {
		notebookScope = notebook.ScopeUser
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
		server.Logger.Warn("Invalid notebook owner id - ", err)
		rpc.SetRPCError(response.Header, codes.ErrorDecode)
		return response, nil
	}

	containerID, err := uuid.FromString(request.ContainerId)
	if err != nil {
		server.Logger.Warn("Invalid notebook container id - ", err)
		rpc.SetRPCError(response.Header, codes.ErrorDecode)
		return response, nil
	}

	t := rpc.MessageToTitle(request.Name)
	notebook, err := notebook.New(t, notebookScope, container, server.DBRegistry, server.Logger)
	if err != nil {
		server.Logger.Warn("Error creating notebook - ", err)
		rpc.SetRPCError(response.Header, codes.ErrorCreate)
		return response, nil
	}
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
func getNotebooks(server *rpc.Server, message []byte, scope string, context *rpc.RequestContext) (proto.Message, error) {
	response := &messages.GetNotebooksResponse{
		Header: rpc.NewResponseHeader(),
	}

	request := messages.GetNotebooksRequest{}
	err := proto.Unmarshal(message, &request)
	if err != nil {
		server.Logger.Warn("Error unmarshaling get notebooks request - ", err)
		rpc.SetRPCError(response.Header, codes.ErrorDecode)
		return response, nil
	}

	if !server.IsSignedIn() {
		rpc.SetRPCError(response.Header, codes.ErrorUnauthorized)
		return response, nil
	}

	var notebookScope notebook.Scope
	var container notebook.ContainerType
	if scope == "account" {
		notebookScope = notebook.ScopeAccount
	} else if scope == "user" {
		notebookScope = notebook.ScopeUser
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
		server.Logger.Warn("Invalid notebook owner id - ", err)
		rpc.SetRPCError(response.Header, codes.ErrorDecode)
		return response, nil
	}

	containerID, err := uuid.FromString(request.ContainerId)
	if err != nil {
		server.Logger.Warn("Invalid notebook container id - ", err)
		rpc.SetRPCError(response.Header, codes.ErrorDecode)
		return response, nil
	}

	// create a new notebook instance to act as a proxy
	nb, err := notebook.New(nil, notebookScope, container, server.DBRegistry, server.Logger)
	if err != nil {
		server.Logger.Warn("Error creating notebook - ", err)
		rpc.SetRPCError(response.Header, codes.ErrorCreate)
		return response, nil
	}
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

func saveNotebook(server *rpc.Server, message []byte, scope string, context *rpc.RequestContext) (proto.Message, error) {
	response := &messages.EmptyResponse{
		Header: rpc.NewResponseHeader(),
	}

	request := messages.SaveNotebookRequest{}
	err := proto.Unmarshal(message, &request)
	if err != nil {
		server.Logger.Warn("Error unmarshaling save notebook request - ", err)
		rpc.SetRPCError(response.Header, codes.ErrorDecode)
		return response, nil
	}

	if !server.IsSignedIn() {
		rpc.SetRPCError(response.Header, codes.ErrorUnauthorized)
		return response, nil
	}

	var notebookScope notebook.Scope
	var container notebook.ContainerType
	if scope == "account" {
		notebookScope = notebook.ScopeAccount
	} else if scope == "user" {
		notebookScope = notebook.ScopeUser
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
		server.Logger.Warn("Invalid notebook id - ", err)
		rpc.SetRPCError(response.Header, codes.ErrorDecode)
		return response, nil
	}

	ownerID, err := uuid.FromString(request.OwnerId)
	if err != nil {
		server.Logger.Warn("Invalid notebook owner id - ", err)
		rpc.SetRPCError(response.Header, codes.ErrorDecode)
		return response, nil
	}

	containerID, err := uuid.FromString(request.ContainerId)
	if err != nil {
		server.Logger.Warn("Invalid notebook container id - ", err)
		rpc.SetRPCError(response.Header, codes.ErrorDecode)
		return response, nil
	}

	t := rpc.MessageToTitle(request.Name)
	notebook, err := notebook.New(t, notebookScope, container, server.DBRegistry, server.Logger)
	if err != nil {
		server.Logger.Warn("Error creating notebook - ", err)
		rpc.SetRPCError(response.Header, codes.ErrorCreate)
		return response, nil
	}
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

func deleteNotebook(server *rpc.Server, message []byte, scope string, context *rpc.RequestContext) (proto.Message, error) {
	response := &messages.EmptyResponse{
		Header: rpc.NewResponseHeader(),
	}

	request := messages.DeleteNotebookRequest{}
	err := proto.Unmarshal(message, &request)
	if err != nil {
		server.Logger.Warn("Error unmarshaling create notebook request - ", err)
		rpc.SetRPCError(response.Header, codes.ErrorDecode)
		return response, nil
	}

	if !server.IsSignedIn() {
		rpc.SetRPCError(response.Header, codes.ErrorUnauthorized)
		return response, nil
	}

	var notebookScope notebook.Scope
	var container notebook.ContainerType
	if scope == "account" {
		notebookScope = notebook.ScopeAccount
	} else if scope == "user" {
		notebookScope = notebook.ScopeUser
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
		server.Logger.Warn("Invalid notebook id - ", err)
		rpc.SetRPCError(response.Header, codes.ErrorDecode)
		return response, nil
	}

	ownerID, err := uuid.FromString(request.OwnerId)
	if err != nil {
		server.Logger.Warn("Invalid notebook owner id - ", err)
		rpc.SetRPCError(response.Header, codes.ErrorDecode)
		return response, nil
	}

	containerID, err := uuid.FromString(request.ContainerId)
	if err != nil {
		server.Logger.Warn("Invalid notebook container id - ", err)
		rpc.SetRPCError(response.Header, codes.ErrorDecode)
		return response, nil
	}

	notebook, err := notebook.New(nil, notebookScope, container, server.DBRegistry, server.Logger)
	if err != nil {
		server.Logger.Warn("Error creating notebook - ", err)
		rpc.SetRPCError(response.Header, codes.ErrorCreate)
		return response, nil
	}
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
