package reqcontext

import "context"

type contextKey string

const (
	RequestIDKey contextKey = "request_id"
	IPAddressKey contextKey = "ip_address"
	UserAgentKey contextKey = "user_agent"
)

func SetRequestDetails(ctx context.Context, requestID, ip, userAgent string) context.Context {
	ctx = context.WithValue(ctx, RequestIDKey, requestID)
	ctx = context.WithValue(ctx, IPAddressKey, ip)
	ctx = context.WithValue(ctx, UserAgentKey, userAgent)
	return ctx
}

func GetRequestDetails(ctx context.Context) (requestID string, ip *string, userAgent *string) {
	if reqID, ok := ctx.Value(RequestIDKey).(string); ok {
		requestID = reqID
	}
	if ipVal, ok := ctx.Value(IPAddressKey).(string); ok && ipVal != "" {
		ip = &ipVal
	}
	if uaVal, ok := ctx.Value(UserAgentKey).(string); ok && uaVal != "" {
		userAgent = &uaVal
	}
	return
}
