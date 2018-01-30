package main

import (
	"errors"
	"net/http"
	"net/smtp"
	"os"
	"strings"

	"github.com/labstack/echo"
)

type Mail struct {
	To      []string `json:"to"`
	Subject string   `json:"subject"`
	Message string   `json:"message"`
}

func Send(c echo.Context) error {
	m := new(Mail)
	if err := c.Bind(m); err != nil {
		return err
	}

	if err := postMail(m); err == nil {
		return c.JSON(http.StatusCreated, `{status:OK}`)
	} else {
		return c.JSON(http.StatusCreated, `{error:`+err.Error()+`}`)
	}
}

func postMail(m *Mail) error {
	gmail := os.Getenv("GMAIL")
	password := os.Getenv("GMAILPW")
	if gmail == "" || password == "" {
		return errors.New("環境変数が正しく設定されていません。")
	}

	auth := smtp.PlainAuth(
		"",
		gmail,    // 送信に使うアカウント
		password, // アカウントのパスワード or アプリケーションパスワード
		"smtp.gmail.com",
	)

	return smtp.SendMail(
		"smtp.gmail.com:587",
		auth,
		gmail, // 送信元
		m.To,  // 送信先
		[]byte(
			"To: "+strings.Join(m.To, ",")+"\r\n"+
				"Subject:"+m.Subject+"\r\n"+
				"\r\n"+
				m.Message),
	)
}
