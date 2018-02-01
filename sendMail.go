package main

import (
	"bytes"
	"encoding/base64"
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

	var body bytes.Buffer
	body.WriteString(m.Message)

	var header bytes.Buffer
	header.WriteString("To: " + strings.Join(m.To, ",") + "\r\n")
	header.WriteString(encodeSubject(m.Subject))
	header.WriteString("MIME-Version: 1.0\r\n")
	header.WriteString("Content-Type: text/plain; charset=\"utf-8\"\r\n")
	header.WriteString("Content-Transfer-Encoding: base64\r\n")

	var message bytes.Buffer
	message = header
	message.WriteString("\r\n")
	message.WriteString(add76crlf(base64.StdEncoding.EncodeToString(body.Bytes())))

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
		[]byte(message.String()),
	)
}

// 76バイト毎にCRLFを挿入する
func add76crlf(msg string) string {
	var buffer bytes.Buffer
	for k, c := range strings.Split(msg, "") {
		buffer.WriteString(c)
		if k%76 == 75 {
			buffer.WriteString("\r\n")
		}
	}
	return buffer.String()
}

// UTF8文字列を指定文字数で分割
func utf8Split(utf8string string, length int) []string {
	resultString := []string{}
	var buffer bytes.Buffer
	for k, c := range strings.Split(utf8string, "") {
		buffer.WriteString(c)
		if k%length == length-1 {
			resultString = append(resultString, buffer.String())
			buffer.Reset()
		}
	}
	if buffer.Len() > 0 {
		resultString = append(resultString, buffer.String())
	}
	return resultString
}

// サブジェクトをMIMEエンコードする
func encodeSubject(subject string) string {
	var buffer bytes.Buffer
	buffer.WriteString("Subject:")
	for _, line := range utf8Split(subject, 13) {
		buffer.WriteString(" =?utf-8?B?")
		buffer.WriteString(base64.StdEncoding.EncodeToString([]byte(line)))
		buffer.WriteString("?=\r\n")
	}
	return buffer.String()
}
