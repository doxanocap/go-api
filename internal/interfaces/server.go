package interfaces

import "context"

type IServer interface {
	REST() IRESTServer
}

type IRESTServer interface {
	Run()
	Shutdown(ctx context.Context) error
}
