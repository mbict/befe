package jwtoken

import "context"

var JwtContextKey int

func FromContext(ctx context.Context) (*Token, bool) {
	token, ok := ctx.Value(&JwtContextKey).(*Token)
	return token, ok
}

func ToContext(ctx context.Context, token *Token) context.Context {
	return context.WithValue(ctx, &JwtContextKey, token)
}
