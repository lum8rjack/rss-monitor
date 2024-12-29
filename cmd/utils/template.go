package utils

import (
	"bytes"
	"errors"
	"path/filepath"
	"text/template"
	"time"
)

type FormData struct {
	Date  string
	Posts []Post
}

type Post struct {
	Published string
	Title     string
	Link      string
}

func GenerateMessage(templateFile string, posts []Post) (string, error) {
	if templateFile == "" {
		return "", errors.New("must provide a template file")
	}

	now := time.Now()
	fd := FormData{
		Date:  now.Format("2006-01-02"),
		Posts: posts,
	}

	t, err := template.ParseFiles(templateFile)
	if err != nil {
		return "", err
	}

	name := filepath.Base(templateFile)

	var buf bytes.Buffer
	err = t.ExecuteTemplate(&buf, name, fd)
	return buf.String(), err
}
