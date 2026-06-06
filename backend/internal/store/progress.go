package store

import (
	"curriculum/internal/model"
	"fmt"
	"time"
)

const classPlanFile = "class_plans.json"
const progressFile = "progress_records.json"

func (s *Store) ListClassPlans() ([]model.ClassPlan, error) {
	var plans []model.ClassPlan
	if err := s.loadData(classPlanFile, &plans); err != nil {
		return nil, err
	}
	if plans == nil {
		plans = []model.ClassPlan{}
	}
	return plans, nil
}

func (s *Store) GetClassPlan(id string) (*model.ClassPlan, error) {
	plans, err := s.ListClassPlans()
	if err != nil {
		return nil, err
	}
	for _, p := range plans {
		if p.ID == id {
			pc := p
			return &pc, nil
		}
	}
	return nil, fmt.Errorf("class plan not found: %s", id)
}

func (s *Store) GetClassPlansByClass(classID string) ([]model.ClassPlan, error) {
	plans, err := s.ListClassPlans()
	if err != nil {
		return nil, err
	}
	var result []model.ClassPlan
	for _, p := range plans {
		if p.ClassID == classID {
			result = append(result, p)
		}
	}
	return result, nil
}

func (s *Store) GetClassPlansByClassAndSubject(classID, subject string) (*model.ClassPlan, error) {
	plans, err := s.ListClassPlans()
	if err != nil {
		return nil, err
	}
	for _, p := range plans {
		if p.ClassID == classID && p.Subject == subject {
			pc := p
			return &pc, nil
		}
	}
	return nil, fmt.Errorf("class plan not found for class %s subject %s", classID, subject)
}

func (s *Store) CreateClassPlan(cp *model.ClassPlan) error {
	plans, err := s.ListClassPlans()
	if err != nil {
		return err
	}
	now := time.Now()
	cp.CreatedAt = now
	cp.UpdatedAt = now
	plans = append(plans, *cp)
	return s.saveData(classPlanFile, plans)
}

func (s *Store) UpdateClassPlan(cp *model.ClassPlan) error {
	plans, err := s.ListClassPlans()
	if err != nil {
		return err
	}
	cp.UpdatedAt = time.Now()
	found := false
	for i, item := range plans {
		if item.ID == cp.ID {
			cp.CreatedAt = item.CreatedAt
			plans[i] = *cp
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("class plan not found: %s", cp.ID)
	}
	return s.saveData(classPlanFile, plans)
}

func (s *Store) DeleteClassPlan(id string) error {
	plans, err := s.ListClassPlans()
	if err != nil {
		return err
	}
	newList := make([]model.ClassPlan, 0, len(plans))
	for _, p := range plans {
		if p.ID != id {
			newList = append(newList, p)
		}
	}
	return s.saveData(classPlanFile, newList)
}

func (s *Store) ListProgressRecords(classPlanID string) ([]model.ProgressRecord, error) {
	var allRecords []model.ProgressRecord
	if err := s.loadData(progressFile, &allRecords); err != nil {
		return nil, err
	}
	if allRecords == nil {
		return []model.ProgressRecord{}, nil
	}
	var result []model.ProgressRecord
	for _, r := range allRecords {
		if r.ClassPlanID == classPlanID {
			result = append(result, r)
		}
	}
	return result, nil
}

func (s *Store) GetAllProgressRecords() ([]model.ProgressRecord, error) {
	var allRecords []model.ProgressRecord
	if err := s.loadData(progressFile, &allRecords); err != nil {
		return nil, err
	}
	if allRecords == nil {
		allRecords = []model.ProgressRecord{}
	}
	return allRecords, nil
}

func (s *Store) AddProgressRecord(record *model.ProgressRecord) error {
	allRecords, err := s.GetAllProgressRecords()
	if err != nil {
		return err
	}
	now := time.Now()
	record.ID = fmt.Sprintf("rec_%d", now.UnixNano())
	record.CreatedAt = now
	if record.ActualDate.IsZero() {
		record.ActualDate = now
	}
	allRecords = append(allRecords, *record)
	return s.saveData(progressFile, allRecords)
}

func (s *Store) GetLastProgressRecord(classPlanID string) (*model.ProgressRecord, error) {
	records, err := s.ListProgressRecords(classPlanID)
	if err != nil {
		return nil, err
	}
	if len(records) == 0 {
		return nil, nil
	}
	last := records[len(records)-1]
	return &last, nil
}

func (s *Store) GetCompletedLessonIDs(classPlanID string) ([]string, error) {
	records, err := s.ListProgressRecords(classPlanID)
	if err != nil {
		return nil, err
	}
	var ids []string
	seen := make(map[string]bool)
	for _, r := range records {
		if r.RecordType == model.RecordNormal || r.RecordType == model.RecordSub {
			if !seen[r.LessonID] {
				ids = append(ids, r.LessonID)
				seen[r.LessonID] = true
			}
		}
	}
	return ids, nil
}
