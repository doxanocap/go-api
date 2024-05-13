package rest

func (r *REST) InitRoutes() {
	router := r.router
	router.GET("/metrics", r.middlewares.GinMetricsHandler())
	router.Use(r.middlewares.ErrorHandler())

	{
		v1 := router.Group("/v1")
		core := v1.Group("/core")
		{
			ws := core.Group("/ws")
			{
				ws.GET("/pool/:user_id/:workspace_id", r.websocket.Pool)
			}

			workspace := core.Group("/workspace")
			{
				workspace.Use(r.middlewares.VerifySession())
				workspace.POST("", r.workspace.Create)
				workspace.GET("/:id", r.workspace.GetByID)

				workspace.GET("/by-user", r.workspace.GetAllByUserID)
			}

		}
	}
}
