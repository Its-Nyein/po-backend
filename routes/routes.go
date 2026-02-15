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
	hashtagRepo := repositories.NewHashtagRepository(db)
	storyRepo := repositories.NewStoryRepository(db)
	convRepo := repositories.NewConversationRepository(db)
	msgRepo := repositories.NewMessageRepository(db)

	userService := services.NewUserService(userRepo)
	postService := services.NewPostService(postRepo)
	commentService := services.NewCommentService(commentRepo)
	likeService := services.NewLikeService(likeRepo)
	followService := services.NewFollowService(followRepo)
	notiService := services.NewNotificationService(notiRepo)
	bookmarkService := services.NewBookmarkService(bookmarkRepo)
	hashtagService := services.NewHashtagService(hashtagRepo)
	storyService := services.NewStoryService(storyRepo, followRepo)
	convService := services.NewConversationService(convRepo, msgRepo, followRepo)

	userCtrl := controllers.NewUserController(userService)
	postCtrl := controllers.NewPostController(postService, followService, hashtagService, notiService, userService)
	commentCtrl := controllers.NewCommentController(commentService, postService, notiService, userService)
	likeCtrl := controllers.NewLikeController(likeService, postService, commentService, notiService, userService)
	followCtrl := controllers.NewFollowController(followService, notiService, userService)
	notiCtrl := controllers.NewNotificationController(notiService)
	bookmarkCtrl := controllers.NewBookmarkController(bookmarkService)
	hashtagCtrl := controllers.NewHashtagController(hashtagService)
	storyCtrl := controllers.NewStoryController(storyService, followService)
	convCtrl := controllers.NewConversationController(convService)

	api.GET("/users", userCtrl.GetUsers)
	api.GET("/users/:id", userCtrl.GetUserByID)
	api.GET("/users/username/:username", userCtrl.GetUserByUsername)
	api.POST("/users", userCtrl.RegisterUser)
	api.POST("/login", userCtrl.LoginUser)
	api.GET("/verify", userCtrl.VerifyToken, middlewares.IsAuthenticated)
	api.GET("/search", userCtrl.SearchUsers)
	api.PUT("/users/profile", userCtrl.UpdateProfile, middlewares.IsAuthenticated)
	api.PUT("/users/password", userCtrl.ChangePassword, middlewares.IsAuthenticated)
	api.DELETE("/users/account", userCtrl.DeleteAccount, middlewares.IsAuthenticated)

	api.GET("/following/users", followCtrl.GetFollowingUsers, middlewares.IsAuthenticated)
	api.GET("/users/:id/followers", followCtrl.GetFollowersByUserID)
	api.GET("/users/:id/following", followCtrl.GetFollowingByUserID)
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

	content.GET("/hashtags/:tag/posts", hashtagCtrl.GetPostsByHashtag, middlewares.IsAuthenticated)
	content.GET("/hashtags/trending", hashtagCtrl.GetTrending, middlewares.IsAuthenticated)

	content.GET("/following/posts", postCtrl.GetFollowingPosts, middlewares.IsAuthenticated)

	content.GET("/notis", notiCtrl.GetNotifications, middlewares.IsAuthenticated)
	content.PUT("/notis/read", notiCtrl.MarkAllRead, middlewares.IsAuthenticated)
	content.PUT("/notis/read/:id", notiCtrl.MarkOneRead, middlewares.IsAuthenticated)

	content.GET("/bookmarks", bookmarkCtrl.GetBookmarks, middlewares.IsAuthenticated)
	content.POST("/bookmarks/:id", bookmarkCtrl.CreateBookmark, middlewares.IsAuthenticated)
	content.DELETE("/bookmarks/:id", bookmarkCtrl.DeleteBookmark, middlewares.IsAuthenticated)

	content.POST("/stories", storyCtrl.CreateStory, middlewares.IsAuthenticated)
	content.DELETE("/stories/:id", storyCtrl.DeleteStory, middlewares.IsAuthenticated, middlewares.IsStoryOwner(storyService))
	content.GET("/stories/feed", storyCtrl.GetFeedStories, middlewares.IsAuthenticated)
	content.GET("/stories/user/:id", storyCtrl.GetUserStories, middlewares.IsAuthenticated)
	content.POST("/stories/:id/view", storyCtrl.RecordView, middlewares.IsAuthenticated)
	content.GET("/stories/:id/viewers", storyCtrl.GetViewers, middlewares.IsAuthenticated, middlewares.IsStoryOwner(storyService))

	content.GET("/conversations", convCtrl.GetConversations, middlewares.IsAuthenticated)
	content.POST("/conversations", convCtrl.CreateConversation, middlewares.IsAuthenticated)
	content.GET("/conversations/unread", convCtrl.GetUnreadCount, middlewares.IsAuthenticated)
	content.GET("/conversations/can-message/:id", convCtrl.CheckMutualFollow, middlewares.IsAuthenticated)
	content.GET("/conversations/:id/messages", convCtrl.GetMessages, middlewares.IsAuthenticated)
	content.POST("/conversations/:id/messages", convCtrl.SendMessage, middlewares.IsAuthenticated)
	content.PUT("/conversations/:id/read", convCtrl.MarkRead, middlewares.IsAuthenticated)

	api.GET("/ws/subscribe", controllers.HandleWebSocket)
}
