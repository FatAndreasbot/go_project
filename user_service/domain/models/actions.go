package models

import "context"

type Interceptor func(context.Context) (context.Context, error)
