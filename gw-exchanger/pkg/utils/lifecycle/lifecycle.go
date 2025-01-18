package lifecycle

import "context"

type Lifecycle interface {
	Start(context context.Context) error
	Stop(context context.Context) error
}
