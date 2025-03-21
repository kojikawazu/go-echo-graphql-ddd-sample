package main

import (
	"backend/config"
	infrastructure_auth "backend/internal/infrastructure/auth"
	infrastructure_todo "backend/internal/infrastructure/todo"
	infrastructure_user "backend/internal/infrastructure/user"
	interfaces_auth "backend/internal/interfaces/auth"
	interfaces_graphql "backend/internal/interfaces/graphql"
	pkg_logger "backend/internal/pkg/logger"
	pkg_supabase "backend/internal/pkg/supabase"
	"backend/internal/router"
	usecase_auth "backend/internal/usecase/auth"
	usecase_todo "backend/internal/usecase/todo"
	usecase_user "backend/internal/usecase/user"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/labstack/echo/v4"
)

// main関数のセットアップ
func setUp(e *echo.Echo, ac *config.AppConfig, l *pkg_logger.AppLogger, sc *pkg_supabase.SupabaseClient) {
	// Supabaseの接続
	err := sc.InitSupabase(l)
	if err != nil {
		l.ErrorLog.Fatalf("Failed to initialize Supabase: %v", err)
	}
	// テストクエリ
	err = sc.TestQuery(l)
	if err != nil {
		l.ErrorLog.Fatalf("Failed to test query: %v", err)
	}

	// DI
	// repository
	userRepository := infrastructure_user.NewUserRepository(l, sc)
	todoRepository := infrastructure_todo.NewTodoRepository(l, sc)
	authRepository := infrastructure_auth.NewAuthRepository(l, sc)
	// usecase
	userUsecase := usecase_user.NewUserUsecase(l, userRepository)
	todoUsecase := usecase_todo.NewTodoUsecase(l, todoRepository)
	authUsecase := usecase_auth.NewAuthUsecase(l, authRepository)
	// handler
	authHandler := interfaces_auth.NewAuthHandler(ac, l)
	// graphql
	graphqlHandler := interfaces_graphql.NewGraphQLHandler(l, userUsecase, todoUsecase, authUsecase, authHandler)

	// router
	router.SetUpRouter(e, l, ac, graphqlHandler, authHandler)
}

// アプリケーションのメイン関数
func main() {
	// 環境変数の読み込み
	appConfig := config.NewAppConfig()
	appConfig.SetUpEnv()

	// ログ設定
	logger := pkg_logger.NewAppLogger()
	logger.SetUpLogger()

	// Supabaseの初期化
	supabaseClient := pkg_supabase.NewSupabaseClient()

	// Echoの設定
	e := echo.New()

	// セットアップ
	setUp(e, appConfig, logger, supabaseClient)

	// シグナルハンドラーの設定
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// 終了ゴルーチン
	go func() {
		<-quit
		logger.InfoLog.Println("Shutting down server...")

		// Echoサーバーのシャットダウン
		if err := e.Close(); err != nil {
			logger.ErrorLog.Printf("Echo shutdown failed: %v", err)
		}

		// Supabaseコネクションプールのクローズ
		supabaseClient.ClosePool(logger)
	}()

	// サーバーの起動
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	if err := e.Start(":" + port); err != nil && err != http.ErrServerClosed {
		logger.ErrorLog.Fatalf("Echo server failed: %v", err)
	}
}
