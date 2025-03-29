package repository_user

import (
	domain_user "backend/internal/domain/user"
)

// ユーザーリポジトリ(IF)
type IUserRepository interface {
	// 全ユーザー取得
	GetAllUsers() ([]domain_user.Users, error)
}
