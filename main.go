package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/kenzo0107/backlog"
)

func main() {
	// .envファイルから環境変数を読み込む
	if err := godotenv.Load(); err != nil {
		fmt.Printf("Error loading .env file: %s\n", err)
		return
	}

	apiKey := os.Getenv("API_KEY")
	baseURL := os.Getenv("BASE_URL")

	c := backlog.New(apiKey, baseURL)

	user, err := c.GetUserMySelf()
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
	fmt.Printf("user ID: %d, Name %s\n", user.ID, *user.Name)
}
