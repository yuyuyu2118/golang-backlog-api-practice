package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/joho/godotenv"
)

func main() {
	// GUIアプリケーションを初期化
	a := app.New()
	w := a.NewWindow("API Response Viewer")

	// .envファイルから環境変数を読み込む
	if err := godotenv.Load(); err != nil {
		w.SetContent(widget.NewLabel(fmt.Sprintf("Error loading .env file: %s", err)))
		w.ShowAndRun()
		return
	}

	apiKey := os.Getenv("API_KEY")
	baseURL := os.Getenv("BASE_URL")

	// APIエンドポイントの入力フィールド
	apiURLInput := widget.NewEntry()
	apiURLInput.SetPlaceHolder("Enter API endpoint (e.g., 'space')")

	// 結果表示用のラベル
	responseLabel := widget.NewLabel("")

	// リクエスト送信ボタン
	sendButton := widget.NewButton("Send Request", func() {
		apiURL := apiURLInput.Text
		url := fmt.Sprintf("%s/api/v2/%s?apiKey=%s", baseURL, apiURL, apiKey)

		client := &http.Client{
			Timeout: 10 * time.Second,
		}

		resp, err := client.Get(url)
		if err != nil {
			responseLabel.SetText(fmt.Sprintf("HTTP request failed: %s", err))
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			responseLabel.SetText(fmt.Sprintf("HTTP request failed with status code: %d", resp.StatusCode))
			return
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			responseLabel.SetText(fmt.Sprintf("Failed to read response body: %s", err))
			return
		}

		var prettyJSON bytes.Buffer
		if error := json.Indent(&prettyJSON, body, "", "    "); error != nil {
			responseLabel.SetText(fmt.Sprintf("Failed to format json: %s", error))
			return
		}

		responseLabel.SetText(prettyJSON.String())
	})

	// UIレイアウトの設定
	content := container.NewVBox(
		apiURLInput,
		sendButton,
		responseLabel,
	)

	w.SetContent(content)
	w.Resize(fyne.NewSize(800, 600))
	w.ShowAndRun()
}
