package main

func main() {
	server := *NewSolanum(5050)

	// server.GET("/posts", func(ctx *gin.Context) {
	// 	ctx.JSON(
	// 		http.StatusOK,
	// 		videoController.FindAll(),
	// 	)
	// })

	// server.POST("/posts", func(ctx *gin.Context) {
	// 	ctx.JSON(
	// 		http.StatusCreated,
	// 		videoController.Save(ctx),
	// 	)
	// })

	server.Run()
}
