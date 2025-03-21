package usecase_todo

import (
	domain_todo "backend/internal/domain/todo"
	pkg_logger "backend/internal/pkg/logger"
	"errors"
)

// Todoユースケース(IF)
type ITodoUsecase interface {
	// 全てのTodoを取得
	GetAllTodos() ([]domain_todo.Todo, error)
	// idを指定してTodoを取得
	GetTodoById(id string) (domain_todo.Todo, error)
	// 特定のユーザーのTodoを取得
	GetTodoByUserId(userId string) ([]domain_todo.Todo, error)
	// 新しいTodoを作成
	CreateTodo(todo domain_todo.Todo) (domain_todo.Todo, error)
	// Todoを更新
	UpdateTodo(todo domain_todo.Todo) (domain_todo.Todo, error)
	// Todoを削除
	DeleteTodo(id string) error
}

// Todoユースケース(Impl)
type TodoUsecase struct {
	Logger         *pkg_logger.AppLogger
	todoRepository domain_todo.ITodoRepository
}

// Todoユースケースのインスタンス化
func NewTodoUsecase(l *pkg_logger.AppLogger, tr domain_todo.ITodoRepository) ITodoUsecase {
	return &TodoUsecase{
		Logger:         l,
		todoRepository: tr,
	}
}

// 全てのTodoを取得
func (u *TodoUsecase) GetAllTodos() ([]domain_todo.Todo, error) {
	u.Logger.InfoLog.Println("GetAllTodos called")

	// Todoリポジトリから全てのTodoを取得(repository層)
	todos, err := u.todoRepository.GetAllTodos()
	if err != nil {
		u.Logger.ErrorLog.Printf("Failed to get all todos: %v", err)
		return nil, err
	}

	u.Logger.InfoLog.Printf("Fetched %d todos", len(todos))
	return todos, nil
}

// idを指定してTodoを取得
func (u *TodoUsecase) GetTodoById(id string) (domain_todo.Todo, error) {
	u.Logger.InfoLog.Println("GetTodoById called")

	// バリデーション
	if id == "" {
		u.Logger.ErrorLog.Println("id is empty")
		return domain_todo.Todo{}, errors.New("id is empty")
	}

	// Todoリポジトリから指定されたidのTodoを取得(repository層)
	todo, err := u.todoRepository.GetTodoById(id)
	if err != nil {
		u.Logger.ErrorLog.Printf("Failed to get todo by id: %v", err)
		return domain_todo.Todo{}, err
	}

	u.Logger.InfoLog.Printf("Fetched todo: %v", todo)
	return todo, nil
}

// 特定のユーザーのTodoを取得
func (u *TodoUsecase) GetTodoByUserId(userId string) ([]domain_todo.Todo, error) {
	u.Logger.InfoLog.Println("GetTodoByUserId called")

	// バリデーション
	if userId == "" {
		u.Logger.ErrorLog.Println("user_id is empty")
		return nil, errors.New("user_id is empty")
	}

	// Todoリポジトリから特定のユーザーのTodoを取得(repository層)
	todos, err := u.todoRepository.GetTodoByUserId(userId)
	if err != nil {
		u.Logger.ErrorLog.Printf("Failed to get todo by user_id: %v", err)
		return nil, err
	}

	u.Logger.InfoLog.Printf("Fetched %d todos", len(todos))
	return todos, nil
}

// 新しいTodoを作成
func (u *TodoUsecase) CreateTodo(todo domain_todo.Todo) (domain_todo.Todo, error) {
	u.Logger.InfoLog.Println("CreateTodo called")

	// バリデーション
	if todo.Description == "" {
		u.Logger.ErrorLog.Println("description is empty")
		return domain_todo.Todo{}, errors.New("description is empty")
	}
	if todo.UserId == "" {
		u.Logger.ErrorLog.Println("user_id is empty")
		return domain_todo.Todo{}, errors.New("user_id is empty")
	}

	// Todoリポジトリから新しいTodoを作成(repository層)
	createdTodo, err := u.todoRepository.CreateTodo(todo)
	if err != nil {
		u.Logger.ErrorLog.Printf("Failed to create todo: %v", err)
		return domain_todo.Todo{}, err
	}

	u.Logger.InfoLog.Printf("Created todo: %v", createdTodo)
	return createdTodo, nil
}

// Todoを更新
func (u *TodoUsecase) UpdateTodo(todo domain_todo.Todo) (domain_todo.Todo, error) {
	u.Logger.InfoLog.Println("UpdateTodo called")

	// バリデーション
	if todo.ID == "" {
		u.Logger.ErrorLog.Println("id is empty")
		return domain_todo.Todo{}, errors.New("id is empty")
	}
	if todo.Description == "" {
		u.Logger.ErrorLog.Println("description is empty")
		return domain_todo.Todo{}, errors.New("description is empty")
	}
	if todo.UserId == "" {
		u.Logger.ErrorLog.Println("user_id is empty")
		return domain_todo.Todo{}, errors.New("user_id is empty")
	}

	// Todoリポジトリから指定されたidのTodoを更新(repository層)
	updatedTodo, err := u.todoRepository.UpdateTodo(todo)
	if err != nil {
		u.Logger.ErrorLog.Printf("Failed to update todo: %v", err)
		return domain_todo.Todo{}, err
	}

	u.Logger.InfoLog.Printf("Updated todo: %v", updatedTodo)
	return updatedTodo, nil
}

// Todoを削除
func (u *TodoUsecase) DeleteTodo(id string) error {
	u.Logger.InfoLog.Println("DeleteTodo called")

	// バリデーション
	if id == "" {
		u.Logger.ErrorLog.Println("id is empty")
		return errors.New("id is empty")
	}

	// Todoリポジトリから指定されたidのTodoを削除(repository層)
	err := u.todoRepository.DeleteTodo(id)
	if err != nil {
		u.Logger.ErrorLog.Printf("Failed to delete todo: %v", err)
		return err
	}

	u.Logger.InfoLog.Printf("Deleted todo: %v", id)
	return nil
}
