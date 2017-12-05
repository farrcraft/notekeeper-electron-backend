package handler

import (
	"../codes"
	messages "../proto"
	"../rpc"
	"../tag"

	"github.com/golang/protobuf/proto"
	uuid "github.com/satori/go.uuid"
)

func strToTagScope(server *rpc.Server, s string, id uuid.UUID) (tag.Scope, bool) {
	var scope tag.Scope
	if s == "account" {
		if server.Account.ID != id {
			return scope, false
		}
		scope = tag.ScopeAccount
	} else if s == "user" {
		if server.Account.ActiveUser.ID != id {
			return scope, false
		}
		scope = tag.ScopeUser
	} else {
		return scope, false
	}
	return scope, true
}

// GetTags gets a set of tags
func GetTags(server *rpc.Server, message []byte) (proto.Message, error) {
	response := &messages.GetTagsResponse{
		Header: rpc.NewResponseHeader(),
	}

	request := messages.GetTagsRequest{}
	err := proto.Unmarshal(message, &request)
	if err != nil {
		server.Logger.Debug("Error unmarshaling get tags request - ", err)
		rpc.SetRPCError(response.Header, codes.ErrorDecode)
		return response, nil
	}

	if !server.IsSignedIn() {
		rpc.SetRPCError(response.Header, codes.ErrorUnauthorized)
		return response, nil
	}

	id, err := uuid.FromString(request.Id)
	if err != nil {
		server.Logger.Debug("Invalid id - ", err)
		rpc.SetRPCError(response.Header, codes.ErrorDecode)
		return response, nil
	}

	scope, ok := strToTagScope(server, request.Scope, id)
	if !ok {
		return response, nil
	}

	// create a new tag instance to act as a proxy
	t := tag.New(nil, scope, server.DBRegistry, server.Logger)
	tags, err := t.LoadAll(server.Account.ActiveUser.PassphraseKey)
	if err != nil {
		rpc.SetInternalError(response.Header, err)
		return response, nil
	}

	for _, t := range tags {
		m := &messages.Tag{
			Id:      t.ID.String(),
			Name:    rpc.TitleToMessage(t.Title),
			Scope:   request.Scope,
			Created: rpc.TimeToMessage(t.Created),
			Updated: rpc.TimeToMessage(t.Updated),
		}
		response.Tags = append(response.Tags, m)
	}

	return response, nil
}

// CreateTag creates a new tag
func CreateTag(server *rpc.Server, message []byte) (proto.Message, error) {
	response := &messages.IdResponse{
		Header: rpc.NewResponseHeader(),
	}

	request := messages.CreateTagRequest{}
	err := proto.Unmarshal(message, &request)
	if err != nil {
		server.Logger.Debug("Error unmarshaling create tag request - ", err)
		rpc.SetRPCError(response.Header, codes.ErrorDecode)
		return response, nil
	}

	if !server.IsSignedIn() {
		rpc.SetRPCError(response.Header, codes.ErrorUnauthorized)
		return response, nil
	}

	id, err := uuid.FromString(request.Id)
	if err != nil {
		server.Logger.Debug("Invalid id - ", err)
		rpc.SetRPCError(response.Header, codes.ErrorDecode)
		return response, nil
	}

	scope, ok := strToTagScope(server, request.Scope, id)
	if !ok {
		return response, nil
	}

	t := rpc.MessageToTitle(request.Name)
	newTag := tag.New(t, scope, server.DBRegistry, server.Logger)
	err = newTag.Save(server.Account.ActiveUser.PassphraseKey)
	if err != nil {
		rpc.SetInternalError(response.Header, err)
	} else {
		response.Id = newTag.ID.String()
	}

	return response, nil
}

// SaveTag saves an existing tag
func SaveTag(server *rpc.Server, message []byte) (proto.Message, error) {
	response := &messages.EmptyResponse{
		Header: rpc.NewResponseHeader(),
	}

	request := messages.SaveTagRequest{}
	err := proto.Unmarshal(message, &request)
	if err != nil {
		server.Logger.Debug("Error unmarshaling save tag request - ", err)
		rpc.SetRPCError(response.Header, codes.ErrorDecode)
		return response, nil
	}

	if !server.IsSignedIn() {
		rpc.SetRPCError(response.Header, codes.ErrorUnauthorized)
		return response, nil
	}

	id, err := uuid.FromString(request.OwnerId)
	if err != nil {
		server.Logger.Debug("Invalid id - ", err)
		rpc.SetRPCError(response.Header, codes.ErrorDecode)
		return response, nil
	}

	scope, ok := strToTagScope(server, request.Scope, id)
	if !ok {
		return response, nil
	}

	t := rpc.MessageToTitle(request.Name)
	newTag := tag.New(t, scope, server.DBRegistry, server.Logger)
	newTag.ID = id
	err = newTag.Save(server.Account.ActiveUser.PassphraseKey)
	if err != nil {
		rpc.SetInternalError(response.Header, err)
	}

	return response, nil
}

// DeleteTag deletes a tag
func DeleteTag(server *rpc.Server, message []byte) (proto.Message, error) {
	response := &messages.EmptyResponse{
		Header: rpc.NewResponseHeader(),
	}

	request := messages.DeleteTagRequest{}
	err := proto.Unmarshal(message, &request)
	if err != nil {
		server.Logger.Debug("Error unmarshaling delete tag request - ", err)
		rpc.SetRPCError(response.Header, codes.ErrorDecode)
		return response, nil
	}

	if !server.IsSignedIn() {
		rpc.SetRPCError(response.Header, codes.ErrorUnauthorized)
		return response, nil
	}

	id, err := uuid.FromString(request.OwnerId)
	if err != nil {
		server.Logger.Debug("Invalid id - ", err)
		rpc.SetRPCError(response.Header, codes.ErrorDecode)
		return response, nil
	}

	scope, ok := strToTagScope(server, request.Scope, id)
	if !ok {
		return response, nil
	}

	t := tag.New(nil, scope, server.DBRegistry, server.Logger)
	t.ID = id
	err = t.Delete(server.Account.ActiveUser.PassphraseKey)
	if err != nil {
		rpc.SetInternalError(response.Header, err)
	}

	return response, nil
}
