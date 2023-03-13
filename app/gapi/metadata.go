package gapi

import (
	"context"

	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

const (
	HTTP_AGENT_KEY = "grpcgateway-user-agent"
	GRPC_AGENT_KEY = "user-agent"
	HTTP_CLIENT_IP_KEY = "x-forwarded-for"
)

type metaData struct {
	UserAgent string
	ClientIP string
}

func (s *Server) extractMetadata(ctx context.Context) (*metaData) {
	mtda := new(metaData)
	
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return mtda
	}

	// http
	if userAgents := md.Get(HTTP_AGENT_KEY); len(userAgents) > 0 {
		mtda.UserAgent = userAgents[0]
	}

	if ips := md.Get(HTTP_CLIENT_IP_KEY); len(ips) > 0 {
		mtda.ClientIP = ips[0]
	}

	// grpc
	if userAgents := md.Get(GRPC_AGENT_KEY); len(userAgents) > 0 {
		mtda.UserAgent = userAgents[0]
	}

	if peer, ok := peer.FromContext(ctx); ok {
		mtda.ClientIP = peer.Addr.String()
	}

	return mtda
}