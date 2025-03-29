package usecase_user

import (
	domain_user "backend/internal/domain/user"
	pkg_logger "backend/internal/pkg/logger"
	repository_user "backend/internal/repository/user"
)

// ユーザーユースケース(IF)
type IUserUsecase interface {
	// 全てのユーザーを取得
	GetAllUsers() ([]domain_user.Users, error)
}

// ユーザーユースケース(Impl)
type UserUsecase struct {
	Logger         *pkg_logger.AppLogger
	userRepository repository_user.IUserRepository
}

// ユーザーユースケースのインスタンス化
func NewUserUsecase(l *pkg_logger.AppLogger, u repository_user.IUserRepository) IUserUsecase {
	return &UserUsecase{
		Logger:         l,
		userRepository: u,
	}
}

// 全てのユーザーを取得
func (u *UserUsecase) GetAllUsers() ([]domain_user.Users, error) {
	u.Logger.InfoLog.Println("GetAllUsers called")

	// ユーザーリポジトリから全てのユーザーを取得(repository層)
	users, err := u.userRepository.GetAllUsers()
	if err != nil {
		u.Logger.ErrorLog.Printf("Failed to get all users: %v", err)
		return nil, err
	}

	u.Logger.InfoLog.Printf("Fetched %d users", len(users))
	return users, nil
}
