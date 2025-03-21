package interfaces_graphql

import (
	domain_todo "backend/internal/domain/todo"
	interfaces_auth "backend/internal/interfaces/auth"
	pkg_logger "backend/internal/pkg/logger"
	usecase_auth "backend/internal/usecase/auth"
	usecase_todo "backend/internal/usecase/todo"
	usecase_user "backend/internal/usecase/user"

	"github.com/graphql-go/graphql"
)

// GraphQLハンドラ(Impl)
type GraphQLHandler struct {
	Logger      *pkg_logger.AppLogger
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
	}
}

// ルートクエリを構築
func (r *GraphQLHandler) BuildRootQuery() *graphql.Object {
	rootQuery := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"users": &graphql.Field{
				Type: graphql.NewList(userType),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					r.Logger.InfoLog.Println("Fetching users...")

					users, err := r.userUsecase.GetAllUsers()
					if err != nil {
						r.Logger.ErrorLog.Printf("Failed to get all users: %v", err)
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

					r.Logger.InfoLog.Printf("Fetched %d users", len(result))
					return result, nil
				},
			},
			"todos": &graphql.Field{
				Type: graphql.NewList(todoType),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					r.Logger.InfoLog.Println("Fetching todos...")

					todos, err := r.todoUsecase.GetAllTodos()
					if err != nil {
						r.Logger.ErrorLog.Printf("Failed to get all todos: %v", err)
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

					r.Logger.InfoLog.Printf("Fetched %d todos", len(result))
					return result, nil
				},
			},
			"todo": &graphql.Field{
				Type: todoType,
				Args: graphql.FieldConfigArgument{"id": &graphql.ArgumentConfig{Type: graphql.String}},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					r.Logger.InfoLog.Println("Fetching todo by id...")

					id := p.Args["id"].(string)
					r.Logger.InfoLog.Printf("Fetching todo by id: %s", id)
					todo, err := r.todoUsecase.GetTodoById(id)
					if err != nil {
						switch err.Error() {
						case "id is empty":
							r.Logger.ErrorLog.Printf("Todo not found: %v", err)
							return nil, err
						default:
							r.Logger.ErrorLog.Printf("Failed to get todo by id: %v", err)
							return nil, err
						}
					}

					result := map[string]interface{}{
						"id":          todo.ID,
						"description": todo.Description,
						"completed":   todo.Completed,
					}

					r.Logger.InfoLog.Printf("Fetched todo: %v", result != nil)
					return result, nil
				},
			},
			"todoByUserId": &graphql.Field{
				Type: graphql.NewList(todoType),
				Args: graphql.FieldConfigArgument{"userId": &graphql.ArgumentConfig{Type: graphql.String}},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					r.Logger.InfoLog.Println("Fetching todo by user id...")

					userId := p.Args["userId"].(string)
					todos, err := r.todoUsecase.GetTodoByUserId(userId)
					if err != nil {
						switch err.Error() {
						case "user_id is empty":
							r.Logger.ErrorLog.Printf("User id is empty: %v", err)
							return nil, err
						default:
							r.Logger.ErrorLog.Printf("Failed to get todo by user id: %v", err)
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

					r.Logger.InfoLog.Printf("Fetched %d todos", len(result))
					return result, nil
				},
			},
		},
	})

	return rootQuery
}

// ルートミューテーションを構築
func (r *GraphQLHandler) BuildRootMutation() *graphql.Object {
	rootMutation := graphql.NewObject(graphql.ObjectConfig{
		Name: "Mutation",
		Fields: graphql.Fields{
			"createTodo": &graphql.Field{
				Type: todoType,
				Args: graphql.FieldConfigArgument{
					"description": &graphql.ArgumentConfig{Type: graphql.String},
					"completed":   &graphql.ArgumentConfig{Type: graphql.Boolean},
					"userId":      &graphql.ArgumentConfig{Type: graphql.String},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					r.Logger.InfoLog.Println("Creating todo...")

					description := p.Args["description"].(string)
					completed := p.Args["completed"].(bool)
					userId := p.Args["userId"].(string)

					todo := domain_todo.Todo{
						Description: description,
						Completed:   completed,
						UserId:      userId,
					}

					createdTodo, err := r.todoUsecase.CreateTodo(todo)
					if err != nil {
						switch err.Error() {
						case "description is empty":
							r.Logger.ErrorLog.Printf("Description is empty: %v", err)
							return nil, err
						case "user_id is empty":
							r.Logger.ErrorLog.Printf("User id is empty: %v", err)
							return nil, err
						default:
							r.Logger.ErrorLog.Printf("Failed to create todo: %v", err)
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
					"userId":      &graphql.ArgumentConfig{Type: graphql.String},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					r.Logger.InfoLog.Println("Updating todo...")

					id := p.Args["id"].(string)
					description := p.Args["description"].(string)
					completed := p.Args["completed"].(bool)
					userId := p.Args["userId"].(string)

					r.Logger.InfoLog.Printf("Updating todo by id: %s", id)
					todo, err := r.todoUsecase.UpdateTodo(domain_todo.Todo{
						ID:          id,
						Description: description,
						Completed:   completed,
						UserId:      userId,
					})

					if err != nil {
						switch err.Error() {
						case "id is empty":
							r.Logger.ErrorLog.Printf("Todo not found: %v", err)
							return nil, err
						case "description is empty":
							r.Logger.ErrorLog.Printf("Description is empty: %v", err)
							return nil, err
						case "user_id is empty":
							r.Logger.ErrorLog.Printf("User id is empty: %v", err)
							return nil, err
						default:
							r.Logger.ErrorLog.Printf("Failed to get todo by id: %v", err)
							return nil, err
						}
					}

					result := map[string]interface{}{
						"id":          todo.ID,
						"description": todo.Description,
						"completed":   todo.Completed,
						"userId":      todo.UserId,
					}

					r.Logger.InfoLog.Printf("Updated todo: %v", result != nil)
					return result, nil
				},
			},
			"deleteTodo": &graphql.Field{
				Type: graphql.Boolean,
				Args: graphql.FieldConfigArgument{"id": &graphql.ArgumentConfig{Type: graphql.String}},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					r.Logger.InfoLog.Println("Deleting todo...")

					id := p.Args["id"].(string)
					err := r.todoUsecase.DeleteTodo(id)
					if err != nil {
						switch err.Error() {
						case "id is empty":
							r.Logger.ErrorLog.Printf("Todo not found: %v", err)
							return nil, err
						default:
							r.Logger.ErrorLog.Printf("Failed to delete todo: %v", err)
							return nil, err
						}
					}

					return true, nil
				},
			},
			"login": &graphql.Field{
				Type: graphql.String,
				Args: graphql.FieldConfigArgument{
					"email":    &graphql.ArgumentConfig{Type: graphql.String},
					"password": &graphql.ArgumentConfig{Type: graphql.String},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					r.Logger.InfoLog.Println("Logging in...")

					email := p.Args["email"].(string)
					password := p.Args["password"].(string)

					token, err := r.authUsecase.Login(email, password)
					if err != nil {
						switch err.Error() {
						case "invalid email or password":
							r.Logger.ErrorLog.Printf("Invalid email or password: %v", err)
							return nil, err
						case "invalid email format":
							r.Logger.ErrorLog.Printf("Invalid email format: %v", err)
							return nil, err
						default:
							r.Logger.ErrorLog.Printf("Failed to login: %v", err)
							return nil, err
						}
					}

					// JWTトークンを生成
					tokenString, err := r.authHandler.GenerateToken(token)
					if err != nil {
						r.Logger.ErrorLog.Printf("Failed to generate token: %v", err)
						return nil, err
					}

					r.Logger.InfoLog.Printf("Logged in: %v", tokenString != "")
					return tokenString, nil
				},
			},
		},
	})

	return rootMutation
}

// スキーマを構築
func (r *GraphQLHandler) GetSchema() graphql.Schema {
	schema, _ := graphql.NewSchema(graphql.SchemaConfig{
		Query:    r.BuildRootQuery(),
		Mutation: r.BuildRootMutation(),
	})

	return schema
}
