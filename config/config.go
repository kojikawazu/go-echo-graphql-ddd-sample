package config

import (
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

// アプリケーションの設定
type AppConfig struct {
	TestAPI   string
	UserID    string
	UserRole  string
	JWTSecret string
}

// アプリケーションの設定のインスタンス化
func NewAppConfig() *AppConfig {
	return &AppConfig{}
}

// プロジェクトのルートディレクトリを特定する関数
func (c *AppConfig) getProjectRoot() string {
	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get current working directory: %v", err)
	}

	for {
		if _, err := os.Stat(filepath.Join(currentDir, "Makefile")); err == nil {
			return currentDir
		}

		parentDir := filepath.Dir(currentDir)
		if parentDir == currentDir {
			log.Fatal("Could not find project root (Makefile or go.mod not found)")
		}
		currentDir = parentDir
	}
}

// 環境変数の読み込み
func (c *AppConfig) SetUpEnv() {
	mode := os.Getenv("TEST_MODE")

	// テストによって環境変数ファイルを変える
	var envFilePath string
	if mode == "true" {
		envFilePath = ".env.test"
	} else {
		envFilePath = ".env"
	}
	// テストによってパスを変える
	projectRoot := c.getProjectRoot()
	absPath := filepath.Join(projectRoot, envFilePath)

	// 環境変数の読み込み
	err := godotenv.Load(absPath)
	if err != nil {
		log.Println("No " + absPath + " file found")
	}

	c.TestAPI = os.Getenv("TEST_API")
	c.UserID = os.Getenv("USER_ID")
	c.UserRole = os.Getenv("ROLE_USER")
	c.JWTSecret = os.Getenv("JWT_SECRET")
}
