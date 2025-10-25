package middleware

import (
	"net/http"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

const (
	
	SecureVerificationSessionKey = "secure_verified_at"
	
	SecureVerificationTimeout = 300 
)




func SecureVerificationRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		
		userId := c.GetInt("id")
		if userId == 0 {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Not logged in",
			})
			c.Abort()
			return
		}

		
		session := sessions.Default(c)
		verifiedAtRaw := session.Get(SecureVerificationSessionKey)

		if verifiedAtRaw == nil {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "Need security verification",
				"code":    "VERIFICATION_REQUIRED",
			})
			c.Abort()
			return
		}

		verifiedAt, ok := verifiedAtRaw.(int64)
		if !ok {
			
			session.Delete(SecureVerificationSessionKey)
			_ = session.Save()
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "Verification status is abnormal, please re-verify.",
				"code":    "VERIFICATION_INVALID",
			})
			c.Abort()
			return
		}

		
		elapsed := time.Now().Unix() - verifiedAt
		if elapsed >= SecureVerificationTimeout {
			
			session.Delete(SecureVerificationSessionKey)
			_ = session.Save()
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "Verification has expired, please re-verify.",
				"code":    "VERIFICATION_EXPIRED",
			})
			c.Abort()
			return
		}

		
		c.Next()
	}
}




func OptionalSecureVerification() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.GetInt("id")
		if userId == 0 {
			c.Set("secure_verified", false)
			c.Next()
			return
		}

		session := sessions.Default(c)
		verifiedAtRaw := session.Get(SecureVerificationSessionKey)

		if verifiedAtRaw == nil {
			c.Set("secure_verified", false)
			c.Next()
			return
		}

		verifiedAt, ok := verifiedAtRaw.(int64)
		if !ok {
			c.Set("secure_verified", false)
			c.Next()
			return
		}

		elapsed := time.Now().Unix() - verifiedAt
		if elapsed >= SecureVerificationTimeout {
			session.Delete(SecureVerificationSessionKey)
			_ = session.Save()
			c.Set("secure_verified", false)
			c.Next()
			return
		}

		c.Set("secure_verified", true)
		c.Set("secure_verified_at", verifiedAt)
		c.Next()
	}
}



func ClearSecureVerification(c *gin.Context) {
	session := sessions.Default(c)
	session.Delete(SecureVerificationSessionKey)
	_ = session.Save()
}
