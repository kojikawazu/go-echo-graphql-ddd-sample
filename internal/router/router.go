package router

import (
	"backend/config"
	interfaces_auth "backend/internal/interfaces/auth"
	interfaces_graphql "backend/internal/interfaces/graphql"
	pkg_logger "backend/internal/pkg/logger"
	"net/http"

	"github.com/graphql-go/graphql"
	"github.com/labstack/echo/v4"
)

// ルーティングの設定
func SetUpRouter(e *echo.Echo, l *pkg_logger.AppLogger, conf *config.AppConfig, gh *interfaces_graphql.GraphQLHandler, ah *interfaces_auth.AuthHandler) {
	l.InfoLog.Println("Setting up router...")

	// GraphQLのルーティング
	e.POST("/graphql", func(c echo.Context) error {
		// JSON ボディから `query` を取り出す
		var body struct {
			Query     string                 `json:"query"`
			Variables map[string]interface{} `json:"variables"`
		}
		err := c.Bind(&body)
		if err != nil || body.Query == "" {
			l.ErrorLog.Println("Invalid GraphQL query", err)
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid GraphQL query"})
		}

		// トークンを取得
		changedCtx, _ := ah.ParseAndAuthorizeToken(c, conf.UserRole)

		// GraphQLの実行
		result := graphql.Do(graphql.Params{
			Schema:         gh.GetSchema(),
			RequestString:  body.Query,
			Context:        changedCtx,
			VariableValues: body.Variables,
		})

		if len(result.Errors) > 0 {
			l.ErrorLog.Println("GraphQL errors", result.Errors)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "GraphQL errors"})
		}

		return c.JSON(http.StatusOK, result)
	})

	l.InfoLog.Println("Router setup complete")
}
