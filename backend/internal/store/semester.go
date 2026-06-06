package store

import (
	"curriculum/internal/model"
	"fmt"
	"time"
)

const semesterFile = "semesters.json"

func (s *Store) ListSemesters() ([]model.Semester, error) {
	var semesters []model.Semester
	if err := s.loadData(semesterFile, &semesters); err != nil {
		return nil, err
	}
	if semesters == nil {
		semesters = []model.Semester{}
	}
	return semesters, nil
}

func (s *Store) GetSemester(id string) (*model.Semester, error) {
	semesters, err := s.ListSemesters()
	if err != nil {
		return nil, err
	}
	for _, sem := range semesters {
		if sem.ID == id {
			semCopy := sem
			return &semCopy, nil
		}
	}
	return nil, fmt.Errorf("semester not found: %s", id)
}

func (s *Store) CreateSemester(sem *model.Semester) error {
	semesters, err := s.ListSemesters()
	if err != nil {
		return err
	}
	now := time.Now()
	sem.CreatedAt = now
	sem.UpdatedAt = now
	if sem.Status == "" {
		sem.Status = model.StatusDraft
	}
	if sem.Version == "" {
		sem.Version = "v1.0"
	}
	if sem.TotalWeeks == 0 {
		sem.TotalWeeks = 18
	}
	if sem.LessonsPerDay == 0 {
		sem.LessonsPerDay = 6
	}
	if sem.LessonMinutes == 0 {
		sem.LessonMinutes = 45
	}
	semesters = append(semesters, *sem)
	return s.saveData(semesterFile, semesters)
}

func (s *Store) UpdateSemester(sem *model.Semester) error {
	semesters, err := s.ListSemesters()
	if err != nil {
		return err
	}
	sem.UpdatedAt = time.Now()
	found := false
	for i, item := range semesters {
		if item.ID == sem.ID {
			semesters[i] = *sem
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("semester not found: %s", sem.ID)
	}
	return s.saveData(semesterFile, semesters)
}

func (s *Store) DeleteSemester(id string) error {
	semesters, err := s.ListSemesters()
	if err != nil {
		return err
	}
	newList := make([]model.Semester, 0, len(semesters))
	for _, sem := range semesters {
		if sem.ID != id {
			newList = append(newList, sem)
		}
	}
	return s.saveData(semesterFile, newList)
}

func (s *Store) AddUnit(semesterID string, unit *model.Unit) error {
	sem, err := s.GetSemester(semesterID)
	if err != nil {
		return err
	}
	now := time.Now()
	unit.SemesterID = semesterID
	unit.CreatedAt = now
	unit.UpdatedAt = now
	unit.OrderIndex = len(sem.Units) + 1
	sem.Units = append(sem.Units, *unit)
	sem.UpdatedAt = now
	return s.UpdateSemester(sem)
}

func (s *Store) UpdateUnit(semesterID string, unit *model.Unit) error {
	sem, err := s.GetSemester(semesterID)
	if err != nil {
		return err
	}
	unit.UpdatedAt = time.Now()
	found := false
	for i, u := range sem.Units {
		if u.ID == unit.ID {
			unit.CreatedAt = u.CreatedAt
			sem.Units[i] = *unit
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("unit not found: %s", unit.ID)
	}
	sem.UpdatedAt = time.Now()
	return s.UpdateSemester(sem)
}

func (s *Store) DeleteUnit(semesterID string, unitID string) error {
	sem, err := s.GetSemester(semesterID)
	if err != nil {
		return err
	}
	newUnits := make([]model.Unit, 0, len(sem.Units))
	for _, u := range sem.Units {
		if u.ID != unitID {
			newUnits = append(newUnits, u)
		}
	}
	sem.Units = newUnits
	sem.UpdatedAt = time.Now()
	return s.UpdateSemester(sem)
}

func (s *Store) AddLesson(semesterID string, unitID string, lesson *model.Lesson) error {
	sem, err := s.GetSemester(semesterID)
	if err != nil {
		return err
	}
	now := time.Now()
	lesson.UnitID = unitID
	lesson.CreatedAt = now
	lesson.UpdatedAt = now
	lesson.DurationMinutes = sem.LessonMinutes
	found := false
	for i, u := range sem.Units {
		if u.ID == unitID {
			lesson.OrderInUnit = len(u.Lessons) + 1
			sem.Units[i].Lessons = append(u.Lessons, *lesson)
			sem.Units[i].LessonCount = len(sem.Units[i].Lessons)
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("unit not found: %s", unitID)
	}
	sem.UpdatedAt = now
	return s.UpdateSemester(sem)
}

func (s *Store) UpdateLesson(semesterID string, unitID string, lesson *model.Lesson) error {
	sem, err := s.GetSemester(semesterID)
	if err != nil {
		return err
	}
	lesson.UpdatedAt = time.Now()
	found := false
	for i, u := range sem.Units {
		if u.ID == unitID {
			for j, l := range u.Lessons {
				if l.ID == lesson.ID {
					lesson.CreatedAt = l.CreatedAt
					sem.Units[i].Lessons[j] = *lesson
					found = true
					break
				}
			}
			break
		}
	}
	if !found {
		return fmt.Errorf("lesson not found: %s", lesson.ID)
	}
	sem.UpdatedAt = time.Now()
	return s.UpdateSemester(sem)
}

func (s *Store) DeleteLesson(semesterID string, unitID string, lessonID string) error {
	sem, err := s.GetSemester(semesterID)
	if err != nil {
		return err
	}
	found := false
	for i, u := range sem.Units {
		if u.ID == unitID {
			newLessons := make([]model.Lesson, 0, len(u.Lessons))
			for _, l := range u.Lessons {
				if l.ID != lessonID {
					newLessons = append(newLessons, l)
				}
			}
			sem.Units[i].Lessons = newLessons
			sem.Units[i].LessonCount = len(newLessons)
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("unit not found: %s", unitID)
	}
	sem.UpdatedAt = time.Now()
	return s.UpdateSemester(sem)
}

func (s *Store) GetUnit(semesterID string, unitID string) (*model.Unit, error) {
	sem, err := s.GetSemester(semesterID)
	if err != nil {
		return nil, err
	}
	for _, u := range sem.Units {
		if u.ID == unitID {
			uc := u
			return &uc, nil
		}
	}
	return nil, fmt.Errorf("unit not found: %s", unitID)
}

func (s *Store) GetLesson(semesterID string, unitID string, lessonID string) (*model.Lesson, error) {
	sem, err := s.GetSemester(semesterID)
	if err != nil {
		return nil, err
	}
	for _, u := range sem.Units {
		if u.ID == unitID {
			for _, l := range u.Lessons {
				if l.ID == lessonID {
					lc := l
					return &lc, nil
				}
			}
		}
	}
	return nil, fmt.Errorf("lesson not found: %s", lessonID)
}

func (s *Store) GetAllLessons(semesterID string) ([]model.Lesson, error) {
	sem, err := s.GetSemester(semesterID)
	if err != nil {
		return nil, err
	}
	var allLessons []model.Lesson
	for _, u := range sem.Units {
		allLessons = append(allLessons, u.Lessons...)
	}
	return allLessons, nil
}
