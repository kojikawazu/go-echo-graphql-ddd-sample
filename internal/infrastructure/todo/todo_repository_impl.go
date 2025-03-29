package infrastructure_todo

import (
	domain_todo "backend/internal/domain/todo"
	pkg_logger "backend/internal/pkg/logger"
	pkg_supabase "backend/internal/pkg/supabase"
	repository_todo "backend/internal/repository/todo"
)

// Todoリポジトリ(Impl)
type TodoRepositoryImpl struct {
	Logger         *pkg_logger.AppLogger
	SupabaseClient *pkg_supabase.SupabaseClient
}

// Todoリポジトリのインスタンス化
func NewTodoRepository(l *pkg_logger.AppLogger, sc *pkg_supabase.SupabaseClient) repository_todo.ITodoRepository {
	return &TodoRepositoryImpl{
		Logger:         l,
		SupabaseClient: sc,
	}
}

// 全てのTodoを取得
func (r *TodoRepositoryImpl) GetAllTodos() ([]domain_todo.Todo, error) {
	r.Logger.InfoLog.Println("GetAllTodos called")

	query := `
		SELECT id, description, completed, user_id, created_at, updated_at
		FROM todos
	`

	// Supabaseからクエリを実行し、条件に一致するユーザーを取得
	rows, err := r.SupabaseClient.Pool.Query(r.SupabaseClient.Ctx, query)
	if err != nil {
		r.Logger.ErrorLog.Printf("Failed to fetch todos: %v", err)
		return nil, err
	}

	// Todosのリストを作成
	todos := []domain_todo.Todo{}
	for rows.Next() {
		var todo domain_todo.Todo
		err = rows.Scan(
			&todo.ID,
			&todo.Description,
			&todo.Completed,
			&todo.UserId,
			&todo.CreatedAt,
			&todo.UpdatedAt,
		)
		if err != nil {
			r.Logger.ErrorLog.Printf("Failed to scan todo: %v", err)
			return nil, err
		}
		todos = append(todos, todo)
	}

	r.Logger.InfoLog.Printf("Fetched %d todos", len(todos))
	return todos, nil
}

// 特定のTodoを取得
func (r *TodoRepositoryImpl) GetTodoById(id string) (domain_todo.Todo, error) {
	r.Logger.InfoLog.Println("GetTodoById called")

	query := `
		SELECT id, description, completed, user_id, created_at, updated_at
		FROM todos
		WHERE id = $1
	`

	// Supabaseからクエリを実行し、条件に一致するユーザーを取得
	var todo domain_todo.Todo
	err := r.SupabaseClient.Pool.QueryRow(r.SupabaseClient.Ctx, query, id).
		Scan(&todo.ID,
			&todo.Description,
			&todo.Completed,
			&todo.UserId,
			&todo.CreatedAt,
			&todo.UpdatedAt,
		)
	if err != nil {
		r.Logger.ErrorLog.Printf("Failed to fetch todo: %v", err)
		return domain_todo.Todo{}, err
	}

	r.Logger.InfoLog.Printf("Fetched todo: %v", todo)
	return todo, nil
}

// 特定のユーザーのTodoを取得
func (r *TodoRepositoryImpl) GetTodoByUserId(userId string) ([]domain_todo.Todo, error) {
	r.Logger.InfoLog.Println("GetTodoByUserId called")

	query := `
		SELECT id, description, completed, user_id, created_at, updated_at
		FROM todos
		WHERE user_id = $1
	`

	// Supabaseからクエリを実行し、条件に一致するユーザーを取得
	rows, err := r.SupabaseClient.Pool.Query(r.SupabaseClient.Ctx, query, userId)
	if err != nil {
		r.Logger.ErrorLog.Printf("Failed to fetch todos: %v", err)
		return nil, err
	}

	// Todosのリストを作成
	todos := []domain_todo.Todo{}
	for rows.Next() {
		var todo domain_todo.Todo
		err = rows.Scan(
			&todo.ID,
			&todo.Description,
			&todo.Completed,
			&todo.UserId,
			&todo.CreatedAt,
			&todo.UpdatedAt,
		)
		if err != nil {
			r.Logger.ErrorLog.Printf("Failed to scan todo: %v", err)
			return nil, err
		}
		todos = append(todos, todo)
	}

	r.Logger.InfoLog.Printf("Fetched %d todos", len(todos))
	return todos, nil
}

// 新しいTodoを作成
func (r *TodoRepositoryImpl) CreateTodo(todo domain_todo.Todo) (domain_todo.Todo, error) {
	r.Logger.InfoLog.Println("CreateTodo called")

	query := `
		INSERT INTO todos (description, completed, user_id)
		VALUES ($1, $2, $3)
		RETURNING id, description, completed, user_id, created_at, updated_at
	`

	// トランザクション開始
	tx, err := r.SupabaseClient.Pool.Begin(r.SupabaseClient.Ctx)
	if err != nil {
		r.Logger.ErrorLog.Printf("Failed to begin transaction: %v", err)
		return domain_todo.Todo{}, err
	}
	defer func() {
		if err != nil {
			r.Logger.ErrorLog.Printf("Failed to rollback transaction: %v", err)
			tx.Rollback(r.SupabaseClient.Ctx)
		}
	}()

	// Supabaseからクエリを実行し、条件に一致するユーザーを取得
	err = tx.QueryRow(r.SupabaseClient.Ctx, query, todo.Description, todo.Completed, todo.UserId).
		Scan(&todo.ID,
			&todo.Description,
			&todo.Completed,
			&todo.UserId,
			&todo.CreatedAt,
			&todo.UpdatedAt,
		)
	if err != nil {
		r.Logger.ErrorLog.Printf("Failed to create todo: %v", err)
		return domain_todo.Todo{}, err
	}

	// トランザクションをコミット
	err = tx.Commit(r.SupabaseClient.Ctx)
	if err != nil {
		r.Logger.ErrorLog.Printf("Failed to commit transaction: %v", err)
		return domain_todo.Todo{}, err
	}

	// 正常系にし、ロールバックを防ぐ
	err = nil

	r.Logger.InfoLog.Printf("Created todo: %v", todo)
	return todo, nil
}

// 特定のTodoを更新
func (r *TodoRepositoryImpl) UpdateTodo(todo domain_todo.Todo) (domain_todo.Todo, error) {
	r.Logger.InfoLog.Println("UpdateTodo called")

	query := `
		UPDATE todos
		SET description = $1, completed = $2, user_id = $3, created_at = $4, updated_at = $5
		WHERE id = $6
		RETURNING id, description, completed, user_id, created_at, updated_at
	`

	// トランザクションを開始
	tx, err := r.SupabaseClient.Pool.Begin(r.SupabaseClient.Ctx)
	if err != nil {
		r.Logger.ErrorLog.Printf("Failed to begin transaction: %v", err)
		return domain_todo.Todo{}, err
	}
	defer func() {
		if err != nil {
			r.Logger.ErrorLog.Printf("Failed to rollback transaction: %v", err)
			tx.Rollback(r.SupabaseClient.Ctx)
		}
	}()

	// Supabaseからクエリを実行し、条件に一致するユーザーを取得
	err = tx.QueryRow(r.SupabaseClient.Ctx, query, todo.Description, todo.Completed, todo.UserId, todo.CreatedAt, todo.UpdatedAt, todo.ID).
		Scan(&todo.ID,
			&todo.Description,
			&todo.Completed,
			&todo.UserId,
			&todo.CreatedAt,
			&todo.UpdatedAt,
		)
	if err != nil {
		r.Logger.ErrorLog.Printf("Failed to update todo: %v", err)
		return domain_todo.Todo{}, err
	}

	// トランザクションをコミット
	err = tx.Commit(r.SupabaseClient.Ctx)
	if err != nil {
		r.Logger.ErrorLog.Printf("Failed to commit transaction: %v", err)
		return domain_todo.Todo{}, err
	}

	// 正常系にし、ロールバックを防ぐ
	err = nil

	r.Logger.InfoLog.Printf("Updated todo: %v", todo)
	return todo, nil
}

// 特定のTodoを削除
func (r *TodoRepositoryImpl) DeleteTodo(id string) error {
	r.Logger.InfoLog.Println("DeleteTodo called")

	query := `
		DELETE FROM todos
		WHERE id = $1
	`

	// トランザクションを開始
	tx, err := r.SupabaseClient.Pool.Begin(r.SupabaseClient.Ctx)
	if err != nil {
		r.Logger.ErrorLog.Printf("Failed to begin transaction: %v", err)
		return err
	}
	defer func() {
		if err != nil {
			r.Logger.ErrorLog.Printf("Failed to rollback transaction: %v", err)
			tx.Rollback(r.SupabaseClient.Ctx)
		}
	}()

	// Supabaseからクエリを実行し、条件に一致するユーザーを取得
	_, err = tx.Exec(r.SupabaseClient.Ctx, query, id)
	if err != nil {
		r.Logger.ErrorLog.Printf("Failed to delete todo: %v", err)
		return err
	}

	// トランザクションをコミット
	err = tx.Commit(r.SupabaseClient.Ctx)
	if err != nil {
		r.Logger.ErrorLog.Printf("Failed to commit transaction: %v", err)
		return err
	}

	// 正常系にし、ロールバックを防ぐ
	err = nil

	r.Logger.InfoLog.Printf("Deleted todo: %v", id)
	return nil
}
