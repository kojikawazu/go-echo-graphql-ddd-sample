package usecase_auth

import (
	pkg_logger "backend/internal/pkg/logger"
	repository_auth "backend/internal/repository/auth"
	"errors"
	"regexp"
)

// 認証ユースケース(IF)
type IAuthUsecase interface {
	// ログイン
	Login(email string, password string) (string, error)
}

// 認証ユースケース(Impl)
type AuthUsecase struct {
	Logger         *pkg_logger.AppLogger
	authRepository repository_auth.IAuthRepository
}

// 認証ユースケースのインスタンス化
func NewAuthUsecase(l *pkg_logger.AppLogger, ar repository_auth.IAuthRepository) IAuthUsecase {
	return &AuthUsecase{
		Logger:         l,
		authRepository: ar,
	}
}

// ログイン
func (u *AuthUsecase) Login(email string, password string) (string, error) {
	u.Logger.InfoLog.Println("Login called")

	// バリデーション
	if email == "" || password == "" {
		u.Logger.ErrorLog.Println("Invalid email or password")
		return "", errors.New("invalid email or password")
	}
	// Emailの形式チェック
	matched, err := regexp.MatchString(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`, email)
	if err != nil || !matched {
		u.Logger.ErrorLog.Println("Invalid email format")
		return "", errors.New("invalid email format")
	}

	// 認証リポジトリからログイン(repository層)
	id, err := u.authRepository.Login(email, password)
	if err != nil {
		u.Logger.ErrorLog.Printf("Failed to login: %v", err)
		return "", errors.New("failed to login")
	}

	u.Logger.InfoLog.Println("Login successful. 1 user found")
	return id, nil
}
