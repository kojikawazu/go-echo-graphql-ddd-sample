package infrastructure_auth

import (
	domain_auth "backend/internal/domain/auth"
	domain_user "backend/internal/domain/user"
	pkg_logger "backend/internal/pkg/logger"
	pkg_supabase "backend/internal/pkg/supabase"
)

// 認証リポジトリの実装(Impl)
type AuthRepositoryImpl struct {
	Logger         *pkg_logger.AppLogger
	SupabaseClient *pkg_supabase.SupabaseClient
}

// 認証リポジトリのインスタンス化
func NewAuthRepository(l *pkg_logger.AppLogger, sc *pkg_supabase.SupabaseClient) domain_auth.IAuthRepository {
	return &AuthRepositoryImpl{
		Logger:         l,
		SupabaseClient: sc,
	}
}

// ログイン
func (r *AuthRepositoryImpl) Login(email string, password string) (string, error) {
	r.Logger.InfoLog.Printf("Logging in with email: %s and password: %s", email, password)

	query := `
        SELECT id, username, email
        FROM users
        WHERE email = $1 AND password = $2
    `

	// Supabaseからクエリを実行し、条件に一致するユーザーを取得
	row := r.SupabaseClient.Pool.QueryRow(r.SupabaseClient.Ctx, query, email, password)

	user := domain_user.Users{}
	err := row.Scan(&user.ID, &user.Username, &user.Email)
	if err != nil {
		r.Logger.ErrorLog.Printf("Failed to fetch user: %v", err)
		return "", err
	}

	r.Logger.InfoLog.Println("Login successful. 1 user found")
	return user.ID, nil
}
