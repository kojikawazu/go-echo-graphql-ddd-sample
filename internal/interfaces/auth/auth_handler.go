package interfaces_auth

import (
	"context"
	"net/http"
	"strings"
	"time"

	"backend/config"
	pkg_logger "backend/internal/pkg/logger"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

// 認証ハンドラ(Impl)
type AuthHandler struct {
	AppConfig *config.AppConfig
	Logger    *pkg_logger.AppLogger
}

// 認証ハンドラのインスタンス化
func NewAuthHandler(ac *config.AppConfig, l *pkg_logger.AppLogger) *AuthHandler {
	return &AuthHandler{
		AppConfig: ac,
		Logger:    l,
	}
}

// JWTトークンを生成
func (h *AuthHandler) GenerateToken(id string) (string, error) {
	h.Logger.InfoLog.Println("Generating token...")

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":   id,
		"role": h.AppConfig.UserRole,
		"exp":  time.Now().Add(time.Hour * 24).Unix(),
	})

	// JWTトークンをシグネーション
	tokenString, err := token.SignedString([]byte(h.AppConfig.JWTSecret))
	if err != nil {
		h.Logger.ErrorLog.Printf("Failed to sign token: %v", err)
		return "", err
	}

	h.Logger.InfoLog.Println("Token generated successfully")
	return tokenString, nil
}

// 認可
func (h *AuthHandler) ParseAndAuthorizeToken(c echo.Context, requiredRole string) (context.Context, error) {
	authHeader := c.Request().Header.Get("Authorization")
	if authHeader == "" {
		h.Logger.ErrorLog.Println("Missing Authorization header")
		return nil, echo.NewHTTPError(http.StatusUnauthorized, "Missing Authorization header")
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			h.Logger.ErrorLog.Println("Invalid token signing method")
			return nil, echo.NewHTTPError(http.StatusUnauthorized, "Invalid token signing method")
		}
		return []byte("secret"), nil
	})
	if err != nil || !token.Valid {
		h.Logger.ErrorLog.Println("Invalid token")
		return nil, echo.NewHTTPError(http.StatusUnauthorized, "Invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		h.Logger.ErrorLog.Println("Invalid token claims")
		return nil, echo.NewHTTPError(http.StatusUnauthorized, "Invalid token claims")
	}

	role, ok := claims["role"].(string)
	if !ok || role != requiredRole {
		h.Logger.ErrorLog.Println("Insufficient permissions")
		return nil, echo.NewHTTPError(http.StatusForbidden, "Insufficient permissions")
	}

	id, ok := claims["id"].(string)
	if !ok {
		h.Logger.ErrorLog.Println("Invalid user ID in token")
		return nil, echo.NewHTTPError(http.StatusUnauthorized, "Invalid user ID in token")
	}

	// context に userID, role を追加
	ctx := context.WithValue(c.Request().Context(), h.AppConfig.UserID, id)
	ctx = context.WithValue(ctx, h.AppConfig.UserRole, role)

	h.Logger.InfoLog.Println("Authorization successful")
	return ctx, nil
}
