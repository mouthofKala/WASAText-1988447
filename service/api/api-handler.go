package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {
	// Register routes
	rt.router.GET("/context", rt.wrap(rt.getContextReply))
	rt.router.GET("/users/:userID", rt.wrap(rt.getUserProfile))
	rt.router.GET("/users", rt.wrap(rt.searchUsers))
	rt.router.PUT("/profile/username", rt.wrap(rt.setMyUsername))
	rt.router.PUT("/profile/photo", rt.wrap(rt.setMyPhoto))
	rt.router.POST("/session", rt.wrap(rt.doLogin))
	rt.router.GET("/chats", rt.wrap(rt.getMyConversations))
	rt.router.POST("/chats/groups", rt.wrap(rt.makeGroup))
	rt.router.POST("/chats", rt.wrap(rt.makeChat))
	rt.router.POST("/chats/:chatID/members", rt.wrap(rt.addToGroup))
	rt.router.PUT("/profile/:chatID/photo", rt.wrap(rt.setGroupPhoto))
	rt.router.PUT("/profile/:chatID/groupname", rt.wrap(rt.setGroupName))
	rt.router.GET("/chats/:chatID/messages", rt.wrap(rt.getConversation))
	rt.router.POST("/chats/:chatID/messages", rt.wrap(rt.sendMessage))
	rt.router.DELETE("/chats/:chatID/messages/:messageID", rt.wrap(rt.deleteMessage))
	rt.router.POST("/chats/:chatID/messages/:messageID/forward", rt.wrap(rt.forwardMessage))
	rt.router.POST("/chats/:chatID/messages/:messageID/reactions", rt.wrap(rt.commentMessage))
	rt.router.DELETE("/chats/:chatID/messages/:messageID/reactions/:reactionID", rt.wrap(rt.uncommentMessage))
	rt.router.DELETE("/chats/:chatID/members/me", rt.wrap(rt.leaveGroup))
	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	return rt.router
}
