package auth

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
)

func InitAuth(r *gin.Engine) {
	goth.UseProviders(
		google.New(os.Getenv("GOOGLE_KEY"), os.Getenv("GOOGLE_SECRET"), os.Getenv("CALLBACK_URL")+"/auth/google/callback"),
	)
	gothic.Store = sessions.NewCookieStore(securecookie.GenerateRandomKey(128))
	// Initialize auth here

	r.GET("/auth/:provider", authEndpoint)
	r.GET("/auth/:provider/callback", callbackEndpoint)
	r.GET("/auth/logout", logoutEndpoint)

}
