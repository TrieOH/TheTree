package jobs

import "context"

// JWKSChangedArgs announces that the published keys changed: the keys module
// created, rotated or retired one. It is not work for this service. On
// Cloudflare the Worker in front of the container intercepts it and drops
// its cached JWKS, so it never reaches the queue; in-process there is no
// such cache and the handler does nothing.
type JWKSChangedArgs struct{}

func (JWKSChangedArgs) Kind() string { return "jwks.changed" }

func workJWKSChanged(context.Context, JWKSChangedArgs) error { return nil }
