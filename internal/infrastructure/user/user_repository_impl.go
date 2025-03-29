package infrastructure_user

import (
	domain_user "backend/internal/domain/user"
	pkg_logger "backend/internal/pkg/logger"
	pkg_supabase "backend/internal/pkg/supabase"
	repository_user "backend/internal/repository/user"
)

// ユーザーリポジトリ(Impl)
type UserRepositoryImpl struct {
	Logger         *pkg_logger.AppLogger
	SupabaseClient *pkg_supabase.SupabaseClient
}

// ユーザーリポジトリのインスタンス化
func NewUserRepository(l *pkg_logger.AppLogger, sc *pkg_supabase.SupabaseClient) repository_user.IUserRepository {
	return &UserRepositoryImpl{
		Logger:         l,
		SupabaseClient: sc,
	}
}

// 全てのユーザーを取得
func (r *UserRepositoryImpl) GetAllUsers() ([]domain_user.Users, error) {
	r.Logger.InfoLog.Printf("Fetching users from Supabase.")

	query := `
        SELECT id, username, email, created_at, updated_at
        FROM users
    `

	// Supabaseからクエリを実行し、条件に一致するユーザーを取得
	rows, err := r.SupabaseClient.Pool.Query(r.SupabaseClient.Ctx, query)
	if err != nil {
		r.Logger.ErrorLog.Printf("Failed to fetch users: %v", err)
		return nil, err
	}

	// ユーザーのリストを作成
	users := []domain_user.Users{}
	for rows.Next() {
		var user domain_user.Users
		err = rows.Scan(
			&user.ID,
			&user.Username,
			&user.Email,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			r.Logger.ErrorLog.Printf("Failed to scan user: %v", err)
			return nil, err
		}
		users = append(users, user)
	}

	// ユーザーのリストを返す
	r.Logger.InfoLog.Printf("Fetched %d users successfully.", len(users))
	return users, nil
}
