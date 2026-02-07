package routes

import (
	"po-backend/controllers"
	"po-backend/middlewares"
	"po-backend/repositories"
	"po-backend/services"

	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

func InitializeRoutes(e *echo.Echo, db *gorm.DB) {
	api := e.Group("/api/v1")

	userRepo := repositories.NewUserRepository(db)
	postRepo := repositories.NewPostRepository(db)
	commentRepo := repositories.NewCommentRepository(db)
	likeRepo := repositories.NewLikeRepository(db)
	followRepo := repositories.NewFollowRepository(db)
	notiRepo := repositories.NewNotificationRepository(db)
	bookmarkRepo := repositories.NewBookmarkRepository(db)

	userService := services.NewUserService(userRepo)
	postService := services.NewPostService(postRepo)
	commentService := services.NewCommentService(commentRepo)
	likeService := services.NewLikeService(likeRepo)
	followService := services.NewFollowService(followRepo)
	notiService := services.NewNotificationService(notiRepo)
	bookmarkService := services.NewBookmarkService(bookmarkRepo)

	userCtrl := controllers.NewUserController(userService)
	postCtrl := controllers.NewPostController(postService, followService)
	commentCtrl := controllers.NewCommentController(commentService, postService, notiService, userService)
	likeCtrl := controllers.NewLikeController(likeService, postService, commentService, notiService, userService)
	followCtrl := controllers.NewFollowController(followService, notiService, userService)
	notiCtrl := controllers.NewNotificationController(notiService)
	bookmarkCtrl := controllers.NewBookmarkController(bookmarkService)

	api.GET("/users", userCtrl.GetUsers)
	api.GET("/users/:id", userCtrl.GetUserByID)
	api.POST("/users", userCtrl.RegisterUser)
	api.POST("/login", userCtrl.LoginUser)
	api.GET("/verify", userCtrl.VerifyToken, middlewares.IsAuthenticated)
	api.GET("/search", userCtrl.SearchUsers)

	api.POST("/follow/:id", followCtrl.FollowUser, middlewares.IsAuthenticated)
	api.DELETE("/unfollow/:id", followCtrl.UnfollowUser, middlewares.IsAuthenticated)

	content := api.Group("/content")

	content.GET("/posts", postCtrl.GetAllPosts)
	content.GET("/posts/:id", postCtrl.GetPostByID)
	content.POST("/posts", postCtrl.CreatePost, middlewares.IsAuthenticated)
	content.PUT("/posts/:id", postCtrl.UpdatePost, middlewares.IsAuthenticated, middlewares.IsPostOwner(postService))
	content.DELETE("/posts/:id", postCtrl.DeletePost, middlewares.IsAuthenticated, middlewares.IsPostOwner(postService))

	content.POST("/comments", commentCtrl.CreateComment, middlewares.IsAuthenticated)
	content.PUT("/comments/:id", commentCtrl.UpdateComment, middlewares.IsAuthenticated, middlewares.IsCommentOwner(commentService))
	content.DELETE("/comments/:id", commentCtrl.DeleteComment, middlewares.IsAuthenticated, middlewares.IsCommentOwner(commentService))

	content.POST("/like/posts/:id", likeCtrl.LikePost, middlewares.IsAuthenticated)
	content.DELETE("/unlike/posts/:id", likeCtrl.UnlikePost, middlewares.IsAuthenticated)
	content.POST("/like/comments/:id", likeCtrl.LikeComment, middlewares.IsAuthenticated)
	content.DELETE("/unlike/comments/:id", likeCtrl.UnlikeComment, middlewares.IsAuthenticated)
	content.GET("/likes/posts/:id", likeCtrl.GetPostLikers)
	content.GET("/likes/comments/:id", likeCtrl.GetCommentLikers)

	content.GET("/following/posts", postCtrl.GetFollowingPosts, middlewares.IsAuthenticated)

	content.GET("/notis", notiCtrl.GetNotifications, middlewares.IsAuthenticated)
	content.PUT("/notis/read", notiCtrl.MarkAllRead, middlewares.IsAuthenticated)
	content.PUT("/notis/read/:id", notiCtrl.MarkOneRead, middlewares.IsAuthenticated)

	content.GET("/bookmarks", bookmarkCtrl.GetBookmarks, middlewares.IsAuthenticated)
	content.POST("/bookmarks/:id", bookmarkCtrl.CreateBookmark, middlewares.IsAuthenticated)
	content.DELETE("/bookmarks/:id", bookmarkCtrl.DeleteBookmark, middlewares.IsAuthenticated)

	api.GET("/ws/subscribe", controllers.HandleWebSocket)
}
