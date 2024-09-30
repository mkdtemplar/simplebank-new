package gapi

import (
	"context"

	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

type Metadata struct {
	UserAgent string `json:"user_agent"`
	ClientIp  string `json:"client_ip"`
}

const (
	grpcGatwayUserAgentHeader = "grpcgateway-user-agent"
	userAgentHeader           = "user-agent"
	xForwardedForHeader       = "x-forwarded-for"
)

func (server *Server) extractMetadata(ctx context.Context) *Metadata {
	mtdt := &Metadata{}

	if md, ok := metadata.FromIncomingContext(ctx); ok {
		// log.Printf("md: %v\n", md)

		if userAgents := md.Get(grpcGatwayUserAgentHeader); len(userAgents) > 0 {
			mtdt.UserAgent = userAgents[0]
		}

		if userAgents := md.Get(userAgentHeader); len(userAgents) > 0 {
			mtdt.UserAgent = userAgents[0]
		}

		if clientIPs := md.Get(xForwardedForHeader); len(clientIPs) > 0 {
			mtdt.ClientIp = clientIPs[0]
		}

		if p, ok := peer.FromContext(ctx); ok {
			mtdt.ClientIp = p.Addr.String()
		}
	}

	return mtdt
}
