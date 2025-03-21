package domain_user

import "time"

// ユーザー情報
type Users struct {
	ID        string    `json:"id"         db:"id"`         // UUID型
	Username  string    `json:"username"   db:"username"`   // ユーザー名
	Email     string    `json:"email"      db:"email"`      // メールアドレス
	Password  string    `json:"password"   db:"password"`   // パスワード
	CreatedAt time.Time `json:"created_at" db:"created_at"` // タイムスタンプ
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"` // タイムスタンプ
}
