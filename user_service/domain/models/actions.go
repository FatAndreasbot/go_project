package models

import "context"

type Interceptor func(context.Context) (context.Context, error)

type contextValueKey string

// points to userID: uuid.UUID
const UserdataKey contextValueKey = "userdata-e65f9095-b3d6-492f-86f4-7801c78c7732"
