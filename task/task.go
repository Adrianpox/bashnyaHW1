package task

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Task struct {
	ID        int       `json:"id"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"created_at"`
	IsDone    bool      `json:"is_done"`
}

var (
	ErrFileNotFound = errors.New("файл не найден")
	ErrEmptyText    = errors.New("текст задачи не может быть пустым")
	ErrInvalidOrder = errors.New("неверный номер задачи")
)

// Load читает файл и возвращает срез задач.
// Если файла нет — возвращает ErrFileNotFound.
func Load(path string) ([]Task, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, ErrFileNotFound
		}
		return nil, err
	}
	if len(data) == 0 {
		return []Task{}, nil
	}

	var tasks []Task
	if err := json.Unmarshal(data, &tasks); err != nil {
		return nil, fmt.Errorf("не удалось разобрать файл как JSON: %w", err)
	}
	return tasks, nil
}

func Save(path string, tasks []Task) error {
	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')

	dir := filepath.Dir(path)
	tmp, err := os.CreateTemp(dir, ".tasks-*.tmp")
	if err != nil {
		return err
	}
	tmpName := tmp.Name()

	renamed := false
	defer func() {
		if !renamed {
			os.Remove(tmpName)
		}
	}()

	if _, err := tmp.Write(data); err != nil {
		tmp.Close()
		return err
	}
	if err := tmp.Sync(); err != nil {
		tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}

	if err := os.Rename(tmpName, path); err != nil {
		return err
	}
	renamed = true
	return nil
}

func Add(path, text string) error {
	text = strings.TrimSpace(text)
	if text == "" {
		return ErrEmptyText
	}

	tasks, err := Load(path)
	if err != nil {
		return err
	}

	tasks = append(tasks, Task{
		ID:        nextID(tasks),
		Text:      text,
		CreatedAt: time.Now(),
		IsDone:    false,
	})
	return Save(path, tasks)
}

func MarkDone(path string, order int) error {
	tasks, err := Load(path)
	if err != nil {
		return err
	}
	if order < 1 || order > len(tasks) {
		return ErrInvalidOrder
	}
	tasks[order-1].IsDone = true
	return Save(path, tasks)
}

func nextID(tasks []Task) int {
	maxID := 0
	for _, t := range tasks {
		if t.ID > maxID {
			maxID = t.ID
		}
	}
	return maxID + 1
}
