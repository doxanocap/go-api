package manager

import (
	"context"
	"go.uber.org/fx"
)

func Run(
	lc fx.Lifecycle,
	manager *Manager,
) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) (err error) {
			manager.Server().REST().Run()
			return
		},
		OnStop: func(ctx context.Context) (err error) {
			if err = manager.Server().REST().Shutdown(ctx); err != nil {
				return err
			}
			//if err = manager.db.Close(); err != nil {
			//	return err
			//}
			//if err = manager.cacheProvider.Close(); err != nil {
			//	return err
			//}
			return nil
		},
	})
}
