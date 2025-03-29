package repository_todo

import (
	domain_todo "backend/internal/domain/todo"
)

// Todoリポジトリ(IF)
type ITodoRepository interface {
	// 全てのTodoを取得
	GetAllTodos() ([]domain_todo.Todo, error)
	// 特定のTodoを取得
	GetTodoById(id string) (domain_todo.Todo, error)
	// 特定のユーザーのTodoを取得
	GetTodoByUserId(userId string) ([]domain_todo.Todo, error)
	// 新しいTodoを作成
	CreateTodo(todo domain_todo.Todo) (domain_todo.Todo, error)
	// 特定のTodoを更新
	UpdateTodo(todo domain_todo.Todo) (domain_todo.Todo, error)
	// 特定のTodoを削除
	DeleteTodo(id string) error
}
