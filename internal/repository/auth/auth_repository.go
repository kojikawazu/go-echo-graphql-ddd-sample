package repository_auth

// 認証リポジトリ(IF)
type IAuthRepository interface {
	// ログイン
	Login(email string, password string) (string, error)
}
