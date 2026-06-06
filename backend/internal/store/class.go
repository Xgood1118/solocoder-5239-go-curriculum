package store

import (
	"curriculum/internal/model"
	"fmt"
	"time"
)

const classFile = "classes.json"
const teacherFile = "teachers.json"

func (s *Store) ListClasses() ([]model.Class, error) {
	var classes []model.Class
	if err := s.loadData(classFile, &classes); err != nil {
		return nil, err
	}
	if classes == nil {
		classes = []model.Class{}
	}
	return classes, nil
}

func (s *Store) GetClass(id string) (*model.Class, error) {
	classes, err := s.ListClasses()
	if err != nil {
		return nil, err
	}
	for _, c := range classes {
		if c.ID == id {
			cc := c
			return &cc, nil
		}
	}
	return nil, fmt.Errorf("class not found: %s", id)
}

func (s *Store) CreateClass(c *model.Class) error {
	classes, err := s.ListClasses()
	if err != nil {
		return err
	}
	now := time.Now()
	c.CreatedAt = now
	c.UpdatedAt = now
	classes = append(classes, *c)
	return s.saveData(classFile, classes)
}

func (s *Store) UpdateClass(c *model.Class) error {
	classes, err := s.ListClasses()
	if err != nil {
		return err
	}
	c.UpdatedAt = time.Now()
	found := false
	for i, item := range classes {
		if item.ID == c.ID {
			classes[i] = *c
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("class not found: %s", c.ID)
	}
	return s.saveData(classFile, classes)
}

func (s *Store) DeleteClass(id string) error {
	classes, err := s.ListClasses()
	if err != nil {
		return err
	}
	newList := make([]model.Class, 0, len(classes))
	for _, c := range classes {
		if c.ID != id {
			newList = append(newList, c)
		}
	}
	return s.saveData(classFile, newList)
}

func (s *Store) ListTeachers() ([]model.Teacher, error) {
	var teachers []model.Teacher
	if err := s.loadData(teacherFile, &teachers); err != nil {
		return nil, err
	}
	if teachers == nil {
		teachers = []model.Teacher{}
	}
	return teachers, nil
}

func (s *Store) GetTeacher(id string) (*model.Teacher, error) {
	teachers, err := s.ListTeachers()
	if err != nil {
		return nil, err
	}
	for _, t := range teachers {
		if t.ID == id {
			tc := t
			return &tc, nil
		}
	}
	return nil, fmt.Errorf("teacher not found: %s", id)
}

func (s *Store) CreateTeacher(t *model.Teacher) error {
	teachers, err := s.ListTeachers()
	if err != nil {
		return err
	}
	now := time.Now()
	t.CreatedAt = now
	t.UpdatedAt = now
	teachers = append(teachers, *t)
	return s.saveData(teacherFile, teachers)
}

func (s *Store) UpdateTeacher(t *model.Teacher) error {
	teachers, err := s.ListTeachers()
	if err != nil {
		return err
	}
	t.UpdatedAt = time.Now()
	found := false
	for i, item := range teachers {
		if item.ID == t.ID {
			teachers[i] = *t
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("teacher not found: %s", t.ID)
	}
	return s.saveData(teacherFile, teachers)
}

func (s *Store) DeleteTeacher(id string) error {
	teachers, err := s.ListTeachers()
	if err != nil {
		return err
	}
	newList := make([]model.Teacher, 0, len(teachers))
	for _, t := range teachers {
		if t.ID != id {
			newList = append(newList, t)
		}
	}
	return s.saveData(teacherFile, newList)
}
