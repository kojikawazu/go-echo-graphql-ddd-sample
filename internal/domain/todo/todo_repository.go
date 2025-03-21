package domain_todo

// Todoリポジトリ(IF)
type ITodoRepository interface {
	// 全てのTodoを取得
	GetAllTodos() ([]Todo, error)
	// 特定のTodoを取得
	GetTodoById(id string) (Todo, error)
	// 特定のユーザーのTodoを取得
	GetTodoByUserId(userId string) ([]Todo, error)
	// 新しいTodoを作成
	CreateTodo(todo Todo) (Todo, error)
	// 特定のTodoを更新
	UpdateTodo(todo Todo) (Todo, error)
	// 特定のTodoを削除
	DeleteTodo(id string) error
}
