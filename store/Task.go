package store

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func (t Task) NextDate(now time.Time) (string, error) {
	err := t.CheckSave()
	if err != nil {
		return "", err
	}
	//при добавлении таски допускается пустой репит, а при некстдата - нет
	if t.Repeat == "" {
		return "", errors.New("поле повтор пустое")
	}
	newD, _ := time.Parse("20060102", t.Date)
	s := strings.Split(t.Repeat, " ")
	for {
		if s[0] == "y" {
			newD = newD.AddDate(1, 0, 0)
		} else if s[0] == "d" && len(s) > 1 {
			days, _ := strconv.Atoi(s[1])
			newD = newD.AddDate(0, 0, days)
		}
		if !newD.Before(now) {
			return newD.Format("20060102"), nil
		}
	}
}

func (t Task) CheckSave() error {
	if t.Title == "" {
		return errors.New("заголовок не указан")
	}
	if len(t.Date) != 0 {
		if _, err := time.Parse("20060102", t.Date); err != nil {
			return err
		}
		s := strings.Split(t.Repeat, " ")
		if s[0] != "y" && s[0] != "d" && s[0] != "" {
			return errors.New("неверный формат повторения: разрешается y, d или пустое значение")
		}
		if s[0] == "d" {
			if len(s) < 2 {
				return errors.New("не указано кол-во дней")
			}
			days, err := strconv.Atoi(s[1])
			if err != nil {
				return err
			}
			if days > 400 {
				return errors.New("максимум 400 дней")
			}
		}
	}
	return nil
}
