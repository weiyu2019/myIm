package ctxdata

import "context"

func GetUId(ctx context.Context) string {
	if u, ok := ctx.Value(Identiy).(string); ok {
		return u
	}
	return ""
}
