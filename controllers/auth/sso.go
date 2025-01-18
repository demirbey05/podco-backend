package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/markbates/goth/gothic"
)

func authEndpoint(c *gin.Context) {
	// Handle auth endpoint here
	provider := c.Param("provider")
	c.Request.URL.RawQuery += "&provider=" + provider

	gothic.BeginAuthHandler(c.Writer, c.Request)

}

func callbackEndpoint(c *gin.Context) {
	// Handle callback endpoint here
	user, err := gothic.CompleteUserAuth(c.Writer, c.Request)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, user)
}

func logoutEndpoint(c *gin.Context) {
	// Handle logout endpoint here
	gothic.Logout(c.Writer, c.Request)
	c.JSON(200, gin.H{"message": "Successfully logged out"})
}
