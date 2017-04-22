package rpc

// RegisterHandlers registers all of the RPC handlers
func (rpc *Server) RegisterHandlers(handlers map[string]Handler) {
	rpc.Handlers = handlers
}
