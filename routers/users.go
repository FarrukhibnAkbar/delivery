package routers

func(r Router) UserRouters() {
	authGroup := r.router.Group("/api/auth")
	authGroup.POST("/sendcode", r.handler.SendCode)
	authGroup.POST("/registr", r.handler.Registration)
	
	authGroup.PUT("/profile", r.handler.UpdateProfile)
	authGroup.GET("/profile", r.handler.GetProfile)
	authGroup.POST("/location", r.handler.InsertUserLocation)
	authGroup.GET("/location", r.handler.GetUserLocation)

	productGroup := r.router.Group("/api")
	productGroup.POST("/favorite", r.handler.AddFavorite)

}