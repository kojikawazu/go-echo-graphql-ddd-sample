package domain_todo

import "time"

// Todo情報
type Todo struct {
	ID          string    `json:"id"          db:"id"`          // UUID型
	Description string    `json:"description" db:"description"` // タスクの説明
	Completed   bool      `json:"completed"   db:"completed"`   // タスクが完了しているかどうか
	UserId      string    `json:"user_id"     db:"user_id"`     // ユーザーID
	CreatedAt   time.Time `json:"created_at" db:"created_at"`   // タイムスタンプ
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`   // タイムスタンプ
}
