package controller

import (
	"fmt"
	"net/http"
	"time"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/model"
	passkeysvc "github.com/QuantumNous/new-api/service/passkey"
	"github.com/QuantumNous/new-api/setting/system_setting"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

const (
	
	SecureVerificationSessionKey = "secure_verified_at"
	
	SecureVerificationTimeout = 300 
)

type UniversalVerifyRequest struct {
	Method string `json:"method"` 
	Code   string `json:"code,omitempty"`
}

type VerificationStatusResponse struct {
	Verified  bool  `json:"verified"`
	ExpiresAt int64 `json:"expires_at,omitempty"`
}



func UniversalVerify(c *gin.Context) {
	userId := c.GetInt("id")
	if userId == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Not logged in",
		})
		return
	}

	var req UniversalVerifyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.ApiError(c, fmt.Errorf("Parameter error: %v", err))
		return
	}

	
	user := &model.User{Id: userId}
	if err := user.FillUserById(); err != nil {
		common.ApiError(c, fmt.Errorf("Failed to retrieve user information: %v", err))
		return
	}

	if user.Status != common.UserStatusEnabled {
		common.ApiError(c, fmt.Errorf("This user has been disabled."))
		return
	}

	
	twoFA, _ := model.GetTwoFAByUserId(userId)
	has2FA := twoFA != nil && twoFA.IsEnabled

	passkey, passkeyErr := model.GetPasskeyByUserID(userId)
	hasPasskey := passkeyErr == nil && passkey != nil

	if !has2FA && !hasPasskey {
		common.ApiError(c, fmt.Errorf("User has not enabled 2FA or Passkey."))
		return
	}

	
	var verified bool
	var verifyMethod string

	switch req.Method {
	case "2fa":
		if !has2FA {
			common.ApiError(c, fmt.Errorf("User has not enabled 2FA."))
			return
		}
		if req.Code == "" {
			common.ApiError(c, fmt.Errorf("The verification code cannot be empty."))
			return
		}
		verified = validateTwoFactorAuth(twoFA, req.Code)
		verifyMethod = "2FA"

	case "passkey":
		if !hasPasskey {
			common.ApiError(c, fmt.Errorf("User has not enabled Passkey"))
			return
		}
		
		
		
		verified = true 
		verifyMethod = "Passkey"

	default:
		common.ApiError(c, fmt.Errorf("Unsupported verification method: %s", req.Method))
		return
	}

	if !verified {
		common.ApiError(c, fmt.Errorf("Verification failed, please check the verification code."))
		return
	}

	
	session := sessions.Default(c)
	now := time.Now().Unix()
	session.Set(SecureVerificationSessionKey, now)
	if err := session.Save(); err != nil {
		common.ApiError(c, fmt.Errorf("Failed to save validation status: %v", err))
		return
	}

	
	model.RecordLog(userId, model.LogTypeSystem, fmt.Sprintf("General security verification successful (Verification method: %s)", verifyMethod))

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Verification successful",
		"data": gin.H{
			"verified":   true,
			"expires_at": now + SecureVerificationTimeout,
		},
	})
}


func GetVerificationStatus(c *gin.Context) {
	userId := c.GetInt("id")
	if userId == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Not logged in",
		})
		return
	}

	session := sessions.Default(c)
	verifiedAtRaw := session.Get(SecureVerificationSessionKey)

	if verifiedAtRaw == nil {
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "",
			"data": VerificationStatusResponse{
				Verified: false,
			},
		})
		return
	}

	verifiedAt, ok := verifiedAtRaw.(int64)
	if !ok {
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "",
			"data": VerificationStatusResponse{
				Verified: false,
			},
		})
		return
	}

	elapsed := time.Now().Unix() - verifiedAt
	if elapsed >= SecureVerificationTimeout {
		
		session.Delete(SecureVerificationSessionKey)
		_ = session.Save()
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "",
			"data": VerificationStatusResponse{
				Verified: false,
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
		"data": VerificationStatusResponse{
			Verified:  true,
			ExpiresAt: verifiedAt + SecureVerificationTimeout,
		},
	})
}



func CheckSecureVerification(c *gin.Context) bool {
	session := sessions.Default(c)
	verifiedAtRaw := session.Get(SecureVerificationSessionKey)

	if verifiedAtRaw == nil {
		return false
	}

	verifiedAt, ok := verifiedAtRaw.(int64)
	if !ok {
		return false
	}

	elapsed := time.Now().Unix() - verifiedAt
	if elapsed >= SecureVerificationTimeout {
		
		session.Delete(SecureVerificationSessionKey)
		_ = session.Save()
		return false
	}

	return true
}



func PasskeyVerifyAndSetSession(c *gin.Context) {
	session := sessions.Default(c)
	now := time.Now().Unix()
	session.Set(SecureVerificationSessionKey, now)
	_ = session.Save()
}



func PasskeyVerifyForSecure(c *gin.Context) {
	if !system_setting.GetPasskeySettings().Enabled {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "Admin has not enabled Passkey login.",
		})
		return
	}

	userId := c.GetInt("id")
	if userId == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Not logged in",
		})
		return
	}

	user := &model.User{Id: userId}
	if err := user.FillUserById(); err != nil {
		common.ApiError(c, fmt.Errorf("Failed to retrieve user information: %v", err))
		return
	}

	if user.Status != common.UserStatusEnabled {
		common.ApiError(c, fmt.Errorf("This user has been disabled."))
		return
	}

	credential, err := model.GetPasskeyByUserID(userId)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "This user has not yet linked a Passkey.",
		})
		return
	}

	wa, err := passkeysvc.BuildWebAuthn(c.Request)
	if err != nil {
		common.ApiError(c, err)
		return
	}

	waUser := passkeysvc.NewWebAuthnUser(user, credential)
	sessionData, err := passkeysvc.PopSessionData(c, passkeysvc.VerifySessionKey)
	if err != nil {
		common.ApiError(c, err)
		return
	}

	_, err = wa.FinishLogin(waUser, *sessionData, c.Request)
	if err != nil {
		common.ApiError(c, err)
		return
	}

	
	now := time.Now()
	credential.LastUsedAt = &now
	if err := model.UpsertPasskeyCredential(credential); err != nil {
		common.ApiError(c, err)
		return
	}

	
	PasskeyVerifyAndSetSession(c)

	
	model.RecordLog(userId, model.LogTypeSystem, "Passkey security verification successful")

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Passkey verification successful",
		"data": gin.H{
			"verified":   true,
			"expires_at": time.Now().Unix() + SecureVerificationTimeout,
		},
	})
}
