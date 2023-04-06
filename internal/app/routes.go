package app

// configureRouting() - метод для конфигурации раутера
func (ac *core) configureRouting() {
	ac.echo.Use(ac.recoverPanic)
	ac.echo.GET("/test/alive", ac.testAlive)
	ac.echo.GET("/test/auth", ac.testAlive, ac.authorize)

	ug := ac.echo.Group("/user")
	ug.POST("/signup", ac.signupUser)
	ug.POST("/login", ac.loginUser)
	ug.PUT("/update", ac.updateUser, ac.authorize)
	ug.PUT("/update/password", ac.updateUserPassword, ac.authorize)
	ug.DELETE("/delete", ac.deleteUser, ac.authorize)

	cg := ac.echo.Group("/country")
	cg.GET("/list", ac.getCountries)
}
