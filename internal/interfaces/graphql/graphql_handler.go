package interfaces_graphql

import (
	domain_todo "backend/internal/domain/todo"
	interfaces_auth "backend/internal/interfaces/auth"
	pkg_logger "backend/internal/pkg/logger"
	pkg_timer "backend/internal/pkg/timer"
	usecase_auth "backend/internal/usecase/auth"
	usecase_todo "backend/internal/usecase/todo"
	usecase_user "backend/internal/usecase/user"
	"errors"

	"github.com/graphql-go/graphql"
)

// GraphQLハンドラ(Impl)
type GraphQLHandler struct {
	Logger      *pkg_logger.AppLogger
	timer       *pkg_timer.TimerPkg
	userUsecase usecase_user.IUserUsecase
	todoUsecase usecase_todo.ITodoUsecase
	authUsecase usecase_auth.IAuthUsecase
	authHandler *interfaces_auth.AuthHandler
}

// GraphQLハンドラのインスタンス化
func NewGraphQLHandler(l *pkg_logger.AppLogger, uu usecase_user.IUserUsecase, tu usecase_todo.ITodoUsecase, au usecase_auth.IAuthUsecase, ah *interfaces_auth.AuthHandler) *GraphQLHandler {
	return &GraphQLHandler{
		Logger:      l,
		userUsecase: uu,
		todoUsecase: tu,
		authUsecase: au,
		authHandler: ah,
		timer:       pkg_timer.NewTimerPkg(),
	}
}

// ルートクエリを構築
func (h *GraphQLHandler) BuildRootQuery() *graphql.Object {
	rootQuery := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"users": &graphql.Field{
				Type: graphql.NewList(userType),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					h.Logger.InfoLog.Println("Fetching users...")
					h.timer.Start()

					userId, ok := p.Context.Value(h.authHandler.AppConfig.UserID).(string)
					if !ok || userId == "" {
						h.Logger.ErrorLog.Println("unauthorized")
						h.Logger.PrintDuration("Fetching users", h.timer.GetDuration())
						return nil, errors.New("unauthorized")
					}

					users, err := h.userUsecase.GetAllUsers()
					if err != nil {
						h.Logger.ErrorLog.Printf("Failed to get all users: %v", err)
						h.Logger.PrintDuration("Fetching users", h.timer.GetDuration())
						return nil, err
					}

					result := make([]map[string]interface{}, 0, len(users))
					for _, u := range users {
						result = append(result, map[string]interface{}{
							"id":       u.ID,
							"username": u.Username,
							"email":    u.Email,
						})
					}

					h.Logger.InfoLog.Printf("Fetched %d users", len(result))
					h.Logger.PrintDuration("Fetching users", h.timer.GetDuration())
					return result, nil
				},
			},
			"todos": &graphql.Field{
				Type: graphql.NewList(todoType),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					h.Logger.InfoLog.Println("Fetching todos...")
					h.timer.Start()

					userId, ok := p.Context.Value(h.authHandler.AppConfig.UserID).(string)
					if !ok || userId == "" {
						h.Logger.ErrorLog.Println("unauthorized")
						h.Logger.PrintDuration("Fetching todos", h.timer.GetDuration())
						return nil, errors.New("unauthorized")
					}

					todos, err := h.todoUsecase.GetAllTodos()
					if err != nil {
						h.Logger.ErrorLog.Printf("Failed to get all todos: %v", err)
						h.Logger.PrintDuration("Fetching todos", h.timer.GetDuration())
						return nil, err
					}

					result := make([]map[string]interface{}, 0, len(todos))
					for _, t := range todos {
						result = append(result, map[string]interface{}{
							"id":          t.ID,
							"description": t.Description,
							"completed":   t.Completed,
						})
					}

					h.Logger.InfoLog.Printf("Fetched %d todos", len(result))
					h.Logger.PrintDuration("Fetching todos", h.timer.GetDuration())
					return result, nil
				},
			},
			"todo": &graphql.Field{
				Type: todoType,
				Args: graphql.FieldConfigArgument{"id": &graphql.ArgumentConfig{Type: graphql.String}},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					h.Logger.InfoLog.Println("Fetching todo by id...")
					h.timer.Start()

					userId, ok := p.Context.Value(h.authHandler.AppConfig.UserID).(string)
					if !ok || userId == "" {
						h.Logger.ErrorLog.Println("unauthorized")
						h.Logger.PrintDuration("Fetching todo by id", h.timer.GetDuration())
						return nil, errors.New("unauthorized")
					}

					id := p.Args["id"].(string)
					h.Logger.InfoLog.Printf("Fetching todo by id: %s", id)
					todo, err := h.todoUsecase.GetTodoById(id)
					if err != nil {
						switch err.Error() {
						case "id is empty":
							h.Logger.ErrorLog.Printf("Todo not found: %v", err)
							h.Logger.PrintDuration("Fetching todo by id", h.timer.GetDuration())
							return nil, err
						default:
							h.Logger.ErrorLog.Printf("Failed to get todo by id: %v", err)
							h.Logger.PrintDuration("Fetching todo by id", h.timer.GetDuration())
							return nil, err
						}
					}

					result := map[string]interface{}{
						"id":          todo.ID,
						"description": todo.Description,
						"completed":   todo.Completed,
					}

					h.Logger.InfoLog.Printf("Fetched todo: %v", result != nil)
					h.Logger.PrintDuration("Fetching todo by id", h.timer.GetDuration())
					return result, nil
				},
			},
			"todoByUserId": &graphql.Field{
				Type: graphql.NewList(todoType),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					h.Logger.InfoLog.Println("Fetching todo by user id...")
					h.timer.Start()

					userId, ok := p.Context.Value(h.authHandler.AppConfig.UserID).(string)
					if !ok || userId == "" {
						h.Logger.ErrorLog.Println("unauthorized")
						h.Logger.PrintDuration("Fetching todo by user id", h.timer.GetDuration())
						return nil, errors.New("unauthorized")
					}

					todos, err := h.todoUsecase.GetTodoByUserId(userId)
					if err != nil {
						switch err.Error() {
						case "user_id is empty":
							h.Logger.ErrorLog.Printf("User id is empty: %v", err)
							h.Logger.PrintDuration("Fetching todo by user id", h.timer.GetDuration())
							return nil, err
						default:
							h.Logger.ErrorLog.Printf("Failed to get todo by user id: %v", err)
							h.Logger.PrintDuration("Fetching todo by user id", h.timer.GetDuration())
							return nil, err
						}
					}

					result := make([]map[string]interface{}, 0, len(todos))
					for _, t := range todos {
						result = append(result, map[string]interface{}{
							"id":          t.ID,
							"description": t.Description,
							"completed":   t.Completed,
						})
					}

					h.Logger.InfoLog.Printf("Fetched %d todos", len(result))
					h.Logger.PrintDuration("Fetching todo by user id", h.timer.GetDuration())
					return result, nil
				},
			},
		},
	})

	return rootQuery
}

// ルートミューテーションを構築
func (h *GraphQLHandler) BuildRootMutation() *graphql.Object {
	rootMutation := graphql.NewObject(graphql.ObjectConfig{
		Name: "Mutation",
		Fields: graphql.Fields{
			"createTodo": &graphql.Field{
				Type: todoType,
				Args: graphql.FieldConfigArgument{
					"description": &graphql.ArgumentConfig{Type: graphql.String},
					"completed":   &graphql.ArgumentConfig{Type: graphql.Boolean},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					h.Logger.InfoLog.Println("Creating todo...")
					h.timer.Start()

					userId, ok := p.Context.Value(h.authHandler.AppConfig.UserID).(string)
					if !ok || userId == "" {
						h.Logger.ErrorLog.Println("unauthorized")
						h.Logger.PrintDuration("Creating todo", h.timer.GetDuration())
						return nil, errors.New("unauthorized")
					}

					description := p.Args["description"].(string)
					completed := p.Args["completed"].(bool)

					todo := domain_todo.Todo{
						Description: description,
						Completed:   completed,
						UserId:      userId,
					}

					createdTodo, err := h.todoUsecase.CreateTodo(todo)
					if err != nil {
						switch err.Error() {
						case "description is empty":
							h.Logger.ErrorLog.Printf("Description is empty: %v", err)
							return nil, err
						case "user_id is empty":
							h.Logger.ErrorLog.Printf("User id is empty: %v", err)
							h.Logger.PrintDuration("Creating todo", h.timer.GetDuration())
							return nil, err
						default:
							h.Logger.ErrorLog.Printf("Failed to create todo: %v", err)
							h.Logger.PrintDuration("Creating todo", h.timer.GetDuration())
							return nil, err
						}
					}

					return map[string]interface{}{
						"id":          createdTodo.ID,
						"description": createdTodo.Description,
						"completed":   createdTodo.Completed,
						"userId":      createdTodo.UserId,
					}, nil
				},
			},
			"updateTodo": &graphql.Field{
				Type: todoType,
				Args: graphql.FieldConfigArgument{
					"id":          &graphql.ArgumentConfig{Type: graphql.String},
					"description": &graphql.ArgumentConfig{Type: graphql.String},
					"completed":   &graphql.ArgumentConfig{Type: graphql.Boolean},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					h.Logger.InfoLog.Println("Updating todo...")
					h.timer.Start()

					userId, ok := p.Context.Value(h.authHandler.AppConfig.UserID).(string)
					if !ok || userId == "" {
						h.Logger.ErrorLog.Println("unauthorized")
						h.Logger.PrintDuration("Updating todo", h.timer.GetDuration())
						return nil, errors.New("unauthorized")
					}

					id := p.Args["id"].(string)
					description := p.Args["description"].(string)
					completed := p.Args["completed"].(bool)

					h.Logger.InfoLog.Printf("Updating todo by id: %s", id)
					todo, err := h.todoUsecase.UpdateTodo(domain_todo.Todo{
						ID:          id,
						Description: description,
						Completed:   completed,
						UserId:      userId,
					})

					if err != nil {
						switch err.Error() {
						case "id is empty":
							h.Logger.ErrorLog.Printf("Todo not found: %v", err)
							h.Logger.PrintDuration("Updating todo", h.timer.GetDuration())
							return nil, err
						case "description is empty":
							h.Logger.ErrorLog.Printf("Description is empty: %v", err)
							h.Logger.PrintDuration("Updating todo", h.timer.GetDuration())
							return nil, err
						case "user_id is empty":
							h.Logger.ErrorLog.Printf("User id is empty: %v", err)
							h.Logger.PrintDuration("Updating todo", h.timer.GetDuration())
							return nil, err
						default:
							h.Logger.ErrorLog.Printf("Failed to get todo by id: %v", err)
							h.Logger.PrintDuration("Updating todo", h.timer.GetDuration())
							return nil, err
						}
					}

					result := map[string]interface{}{
						"id":          todo.ID,
						"description": todo.Description,
						"completed":   todo.Completed,
						"userId":      todo.UserId,
					}

					h.Logger.InfoLog.Printf("Updated todo: %v", result != nil)
					h.Logger.PrintDuration("Updating todo", h.timer.GetDuration())
					return result, nil
				},
			},
			"deleteTodo": &graphql.Field{
				Type: deleteTodoPayload,
				Args: graphql.FieldConfigArgument{"id": &graphql.ArgumentConfig{Type: graphql.String}},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					h.Logger.InfoLog.Println("Deleting todo...")
					h.timer.Start()

					userId, ok := p.Context.Value(h.authHandler.AppConfig.UserID).(string)
					if !ok || userId == "" {
						h.Logger.ErrorLog.Println("unauthorized")
						h.Logger.PrintDuration("Deleting todo", h.timer.GetDuration())
						return nil, errors.New("unauthorized")
					}

					id := p.Args["id"].(string)
					err := h.todoUsecase.DeleteTodo(id)
					if err != nil {
						switch err.Error() {
						case "id is empty":
							h.Logger.ErrorLog.Printf("Todo not found: %v", err)
							h.Logger.PrintDuration("Deleting todo", h.timer.GetDuration())
							return nil, err
						default:
							h.Logger.ErrorLog.Printf("Failed to delete todo: %v", err)
							h.Logger.PrintDuration("Deleting todo", h.timer.GetDuration())
							return nil, err
						}
					}

					h.Logger.InfoLog.Println("Todo deleted successfully")
					h.Logger.PrintDuration("Deleting todo", h.timer.GetDuration())
					return map[string]interface{}{
						"success": true,
						"message": "Todo deleted successfully",
					}, nil
				},
			},
			"login": &graphql.Field{
				Type: loginPayload,
				Args: graphql.FieldConfigArgument{
					"email":    &graphql.ArgumentConfig{Type: graphql.String},
					"password": &graphql.ArgumentConfig{Type: graphql.String},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					h.Logger.InfoLog.Println("Logging in...")
					h.timer.Start()

					email := p.Args["email"].(string)
					password := p.Args["password"].(string)

					token, err := h.authUsecase.Login(email, password)
					if err != nil {
						switch err.Error() {
						case "invalid email or password":
							h.Logger.ErrorLog.Printf("Invalid email or password: %v", err)
							h.Logger.PrintDuration("Logging in", h.timer.GetDuration())
							return nil, err
						case "invalid email format":
							h.Logger.ErrorLog.Printf("Invalid email format: %v", err)
							h.Logger.PrintDuration("Logging in", h.timer.GetDuration())
							return nil, err
						default:
							h.Logger.ErrorLog.Printf("Failed to login: %v", err)
							h.Logger.PrintDuration("Logging in", h.timer.GetDuration())
							return nil, err
						}
					}

					// JWTトークンを生成
					tokenString, err := h.authHandler.GenerateToken(token)
					if err != nil {
						h.Logger.ErrorLog.Printf("Failed to generate token: %v", err)
						h.Logger.PrintDuration("Logging in", h.timer.GetDuration())
						return nil, err
					}

					h.Logger.InfoLog.Printf("Logged in: %v", tokenString != "")
					h.Logger.PrintDuration("Logging in", h.timer.GetDuration())
					return map[string]interface{}{
						"token": tokenString,
					}, nil
				},
			},
		},
	})

	return rootMutation
}

// スキーマを構築
func (h *GraphQLHandler) GetSchema() graphql.Schema {
	schema, _ := graphql.NewSchema(graphql.SchemaConfig{
		Query:    h.BuildRootQuery(),
		Mutation: h.BuildRootMutation(),
	})

	return schema
}
