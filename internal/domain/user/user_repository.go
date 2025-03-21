package domain_user

// ユーザーリポジトリ(IF)
type IUserRepository interface {
	// 全ユーザー取得
	GetAllUsers() ([]Users, error)
}
