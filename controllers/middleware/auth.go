package middleware

import (
	"context"
	"log"
	"net/http"
	"strings"

	firebase "firebase.google.com/go"
	"github.com/gin-gonic/gin"
)

// FirebaseAuthMiddleware verifies the Firebase ID token and stores the UID in the context.
func FirebaseAuthMiddleware(app *firebase.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract the ID token from the Authorization header
		idToken := c.GetHeader("Authorization")
		if idToken == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
			c.Abort()
			return
		}

		// Verify the ID token
		client, err := app.Auth(context.Background())
		if err != nil {
			log.Printf("error getting Auth client: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			c.Abort()
			return
		}

		// Strip the "Bearer " prefix
		parts := strings.Split(idToken, " ")
		token, err := client.VerifyIDToken(context.Background(), parts[1])
		if err != nil {
			log.Printf("error verifying ID token: %v\n", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid ID token"})
			c.Abort()
			return
		}

		// Log the verified token (optional)
		log.Printf("Verified ID token: %v\n", token)

		// Store the UID in the Gin context for use in downstream handlers
		c.Set("firebaseUID", token.UID)

		// Proceed to the next handler
		c.Next()
	}
}
