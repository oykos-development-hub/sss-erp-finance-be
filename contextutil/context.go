package contextutil

import "context"

type contextKey string

const UserIDKey contextKey = "userID"

// SetUserIDInContext sets the user ID in the context
func SetUserIDInContext(ctx context.Context, userID int) context.Context {
	return context.WithValue(ctx, UserIDKey, userID)
}

// GetUserIDFromContext retrieves the user ID from the context
func GetUserIDFromContext(ctx context.Context) (int, bool) {
	userID, ok := ctx.Value(UserIDKey).(int)
	return userID, ok
}
