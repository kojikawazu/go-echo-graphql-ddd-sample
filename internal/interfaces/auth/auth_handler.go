package interfaces_auth

import (
	"net/http"
	"strings"
	"time"

	pkg_logger "backend/internal/pkg/logger"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

// 認証ハンドラ(Impl)
type AuthHandler struct {
	Logger *pkg_logger.AppLogger
}

// 認証ハンドラのインスタンス化
func NewAuthHandler(l *pkg_logger.AppLogger) *AuthHandler {
	return &AuthHandler{
		Logger: l,
	}
}

// JWTトークンを生成
func (h *AuthHandler) GenerateToken(id string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":   id,
		"role": "user",
		"exp":  time.Now().Add(time.Hour * 24).Unix(),
	})

	// JWTトークンをシグネーション
	tokenString, err := token.SignedString([]byte("secret"))
	if err != nil {
		h.Logger.ErrorLog.Printf("Failed to sign token: %v", err)
		return "", err
	}

	h.Logger.InfoLog.Println("Token generated successfully")
	return tokenString, nil
}

// 認証ミドルウェア
func (h *AuthHandler) AuthorizationMiddleware(next echo.HandlerFunc, requiredRole string) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			h.Logger.ErrorLog.Println("Missing Authorization header")
			return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Missing Authorization header"})
		}

		// "Bearer " を取り除く
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// JWT をパース
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				h.Logger.ErrorLog.Println("Invalid token signing method")
				return nil, echo.NewHTTPError(http.StatusUnauthorized, "Invalid token signing method")
			}
			return []byte("secret"), nil
		})

		if err != nil || !token.Valid {
			h.Logger.ErrorLog.Println("Invalid token")
			return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Invalid token"})
		}

		// クレームからユーザーIDとロールを取得
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			h.Logger.ErrorLog.Println("Invalid token claims")
			return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Invalid token claims"})
		}

		// ロールを確認（例: "admin", "user" など）
		role, ok := claims["role"].(string)
		if !ok || role != requiredRole {
			h.Logger.ErrorLog.Println("Insufficient permissions")
			return c.JSON(http.StatusForbidden, map[string]string{"message": "Insufficient permissions"})
		}

		// ユーザーIDをコンテキストに保存
		c.Set("userId", claims["id"])
		c.Set("role", role)

		h.Logger.InfoLog.Println("Authorization successful")
		return next(c)
	}
}
