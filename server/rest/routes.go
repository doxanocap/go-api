package rest

func (r *REST) InitRoutes() {
	router := r.router
	router.GET("/metrics", r.middlewares.GinMetricsHandler())
	router.Use(r.middlewares.ErrorHandler())

	{
		v1 := router.Group("/v1")
		v1.GET("/")
	}
}
