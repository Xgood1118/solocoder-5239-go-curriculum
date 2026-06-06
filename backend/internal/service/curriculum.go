package service

import (
	"curriculum/internal/model"
	"curriculum/internal/store"
	"fmt"
	"math"
)

type CurriculumService struct {
	store *store.Store
}

func NewCurriculumService(s *store.Store) *CurriculumService {
	return &CurriculumService{store: s}
}

func (svc *CurriculumService) SuggestUnitLessons(semesterID, unitID string) (*model.LessonSuggestion, error) {
	unit, err := svc.store.GetUnit(semesterID, unitID)
	if err != nil {
		return nil, err
	}

	kpIDs := unit.KnowledgePointIDs
	if len(kpIDs) == 0 {
		return &model.LessonSuggestion{
			MinLessons: 5,
			MaxLessons: 8,
			Reason:     "单元未配置知识点，按默认经验值建议 5-8 课时",
		}, nil
	}

	totalDifficulty := 0
	totalEstimated := 0
	keyCount := 0
	difficultCount := 0
	prereqCount := 0

	for _, kpID := range kpIDs {
		kp, err := svc.store.GetKnowledgePoint(kpID)
		if err != nil {
			continue
		}
		totalDifficulty += int(kp.Difficulty)
		totalEstimated += kp.EstimatedLessons
		if kp.IsKeyPoint {
			keyCount++
		}
		if kp.IsDifficult {
			difficultCount++
		}
		prereqCount += len(kp.Prerequisites)
	}

	avgDifficulty := float64(totalDifficulty) / float64(len(kpIDs))
	baseLessons := float64(totalEstimated)
	coverageFactor := 1.0

	if keyCount > 0 {
		coverageFactor += 0.1 * float64(keyCount) / float64(len(kpIDs))
	}
	if difficultCount > 0 {
		coverageFactor += 0.15 * float64(difficultCount) / float64(len(kpIDs))
	}
	if prereqCount > 0 {
		coverageFactor += 0.05
	}

	baseLessons *= coverageFactor

	minLessons := int(math.Floor(baseLessons * 0.85))
	maxLessons := int(math.Ceil(baseLessons * 1.2))

	if minLessons < 3 {
		minLessons = 3
	}
	if maxLessons > 12 {
		maxLessons = 12
	}
	if maxLessons < minLessons+1 {
		maxLessons = minLessons + 1
	}

	reason := fmt.Sprintf("基于 %d 个知识点计算：平均难度 %.1f/5，预估基础 %.1f 课时；", len(kpIDs), avgDifficulty, baseLessons)
	if keyCount > 0 {
		reason += fmt.Sprintf("含 %d 个重点，", keyCount)
	}
	if difficultCount > 0 {
		reason += fmt.Sprintf("含 %d 个难点，", difficultCount)
	}
	reason += fmt.Sprintf("建议 %d-%d 课时。", minLessons, maxLessons)

	return &model.LessonSuggestion{
		MinLessons: minLessons,
		MaxLessons: maxLessons,
		Reason:     reason,
	}, nil
}

func (svc *CurriculumService) CalculateProgressSummary(classPlanID string) (*model.ProgressSummary, error) {
	classPlan, err := svc.store.GetClassPlan(classPlanID)
	if err != nil {
		return nil, err
	}

	sem, err := svc.store.GetSemester(classPlan.SemesterID)
	if err != nil {
		return nil, err
	}

	allLessons, err := svc.store.GetAllLessons(classPlan.SemesterID)
	if err != nil {
		return nil, err
	}

	totalLessons := len(allLessons)
	completedIDs, err := svc.store.GetCompletedLessonIDs(classPlanID)
	if err != nil {
		return nil, err
	}
	completedLessons := len(completedIDs)

	currentWeek := svc.calculateCurrentWeek(classPlan)

	planLessonsByWeek := 0
	for _, lesson := range allLessons {
		if lesson.PlanWeek <= currentWeek {
			planLessonsByWeek++
		}
	}

	deviation := completedLessons - planLessonsByWeek

	deviationLevel := "normal"
	if deviation <= -5 {
		deviationLevel = "critical"
	} else if deviation <= -3 {
		deviationLevel = "warning"
	} else if deviation >= 3 {
		deviationLevel = "ahead"
	}

	kpMap := make(map[string]bool)
	completedKP := make(map[string]bool)

	for _, lesson := range allLessons {
		for _, kpID := range lesson.KnowledgePointIDs {
			kpMap[kpID] = true
		}
	}

	completedLessonMap := make(map[string]bool)
	for _, id := range completedIDs {
		completedLessonMap[id] = true
	}

	for _, lesson := range allLessons {
		if completedLessonMap[lesson.ID] {
			for _, kpID := range lesson.KnowledgePointIDs {
				completedKP[kpID] = true
			}
		}
	}

	totalKP := len(kpMap)
	var keyTotal, keyDone, diffTotal, diffDone, examTotal, examDone int

	for kpID := range kpMap {
		kp, err := svc.store.GetKnowledgePoint(kpID)
		if err != nil {
			continue
		}
		if kp.IsKeyPoint {
			keyTotal++
			if completedKP[kpID] {
				keyDone++
			}
		}
		if kp.IsDifficult {
			diffTotal++
			if completedKP[kpID] {
				diffDone++
			}
		}
		if kp.IsExamPoint {
			examTotal++
			if completedKP[kpID] {
				examDone++
			}
		}
	}

	keyRate := 0.0
	if keyTotal > 0 {
		keyRate = float64(keyDone) / float64(keyTotal)
	}
	diffRate := 0.0
	if diffTotal > 0 {
		diffRate = float64(diffDone) / float64(diffTotal)
	}
	examRate := 0.0
	if examTotal > 0 {
		examRate = float64(examDone) / float64(examTotal)
	}

	return &model.ProgressSummary{
		ClassPlanID:          classPlanID,
		TotalLessons:         totalLessons,
		CompletedLessons:     completedLessons,
		PlanLessonsByWeek:    planLessonsByWeek,
		Deviation:            deviation,
		DeviationLevel:       deviationLevel,
		CurrentWeek:          currentWeek,
		KeyCoverageRate:      keyRate,
		DifficultCoverageRate: diffRate,
		ExamCoverageRate:     examRate,
	}, nil
}

func (svc *CurriculumService) calculateCurrentWeek(classPlan *model.ClassPlan) int {
	now := classPlan.StartDate
	days := 0
	today := classPlan.EndDate
	for !now.After(today) && days < 18*7 {
		days++
		now = now.AddDate(0, 0, 1)
	}
	week := days / 7
	if week < 1 {
		week = 1
	}
	if week > 18 {
		week = 18
	}
	return week
}

func (svc *CurriculumService) CalculateCoverageReport(classPlanID string) (*model.CoverageReport, error) {
	classPlan, err := svc.store.GetClassPlan(classPlanID)
	if err != nil {
		return nil, err
	}

	sem, err := svc.store.GetSemester(classPlan.SemesterID)
	if err != nil {
		return nil, err
	}

	completedIDs, err := svc.store.GetCompletedLessonIDs(classPlanID)
	if err != nil {
		return nil, err
	}
	completedMap := make(map[string]bool)
	for _, id := range completedIDs {
		completedMap[id] = true
	}

	var keyPoints []model.CoverageItem
	var diffPoints []model.CoverageItem
	var examPoints []model.CoverageItem
	var unpassedUnits []string

	for _, unit := range sem.Units {
		unitKP := make(map[string]bool)
		unitCompletedKP := make(map[string]bool)

		for _, lesson := range unit.Lessons {
			for _, kpID := range lesson.KnowledgePointIDs {
				unitKP[kpID] = true
			}
			if completedMap[lesson.ID] {
				for _, kpID := range lesson.KnowledgePointIDs {
					unitCompletedKP[kpID] = true
				}
			}
		}

		var keyTotal, keyDone, diffTotal, diffDone, examTotal, examDone int
		unitHasContent := false

		for kpID := range unitKP {
			kp, err := svc.store.GetKnowledgePoint(kpID)
			if err != nil {
				continue
			}
			unitHasContent = true
			if kp.IsKeyPoint {
				keyTotal++
				if unitCompletedKP[kpID] {
					keyDone++
				}
			}
			if kp.IsDifficult {
				diffTotal++
				if unitCompletedKP[kpID] {
					diffDone++
				}
			}
			if kp.IsExamPoint {
				examTotal++
				if unitCompletedKP[kpID] {
					examDone++
				}
			}
		}

		if !unitHasContent {
			continue
		}

		keyRate := 0.0
		if keyTotal > 0 {
			keyRate = float64(keyDone) / float64(keyTotal)
		}
		keyPoints = append(keyPoints, model.CoverageItem{
			UnitID:   unit.ID,
			UnitName: unit.Title,
			Rate:     keyRate,
			IsPassed: keyRate >= 0.8,
		})

		diffRate := 0.0
		if diffTotal > 0 {
			diffRate = float64(diffDone) / float64(diffTotal)
		}
		diffPoints = append(diffPoints, model.CoverageItem{
			UnitID:   unit.ID,
			UnitName: unit.Title,
			Rate:     diffRate,
			IsPassed: diffRate >= 0.8,
		})

		examRate := 0.0
		if examTotal > 0 {
			examRate = float64(examDone) / float64(examTotal)
		}
		examPoints = append(examPoints, model.CoverageItem{
			UnitID:   unit.ID,
			UnitName: unit.Title,
			Rate:     examRate,
			IsPassed: examRate >= 0.8,
		})

		unitMinRate := math.Min(keyRate, math.Min(diffRate, examRate))
		if unitMinRate < 0.8 {
			unpassedUnits = append(unpassedUnits, unit.ID)
		}
	}

	return &model.CoverageReport{
		KeyPoints:     keyPoints,
		DifficultPoints: diffPoints,
		ExamPoints:    examPoints,
		UnpassedUnits: unpassedUnits,
	}, nil
}

func (svc *CurriculumService) GenerateGanttData(classPlanID string) (*model.GanttCompareData, error) {
	classPlan, err := svc.store.GetClassPlan(classPlanID)
	if err != nil {
		return nil, err
	}

	classInfo, err := svc.store.GetClass(classPlan.ClassID)
	if err != nil {
		return nil, err
	}

	sem, err := svc.store.GetSemester(classPlan.SemesterID)
	if err != nil {
		return nil, err
	}

	completedIDs, err := svc.store.GetCompletedLessonIDs(classPlanID)
	if err != nil {
		return nil, err
	}
	completedMap := make(map[string]bool)
	for _, id := range completedIDs {
		completedMap[id] = true
	}

	var bars []model.GanttBar
	for _, unit := range sem.Units {
		if len(unit.Lessons) == 0 {
			continue
		}
		firstLesson := unit.Lessons[0]
		lastLesson := unit.Lessons[len(unit.Lessons)-1]
		startWeek := firstLesson.PlanWeek
		endWeek := lastLesson.PlanWeek

		actualCompleted := 0
		lastCompletedWeek := 0
		for _, lesson := range unit.Lessons {
			if completedMap[lesson.ID] {
				actualCompleted++
				lastCompletedWeek = lesson.PlanWeek
			}
		}
		_ = actualCompleted

		if startWeek == 0 {
			startWeek = unit.OrderIndex
		}
		if endWeek == 0 {
			endWeek = unit.OrderIndex + 1
		}

		bars = append(bars, model.GanttBar{
			UnitID:    unit.ID,
			UnitTitle: unit.Title,
			StartWeek: startWeek,
			EndWeek:   endWeek,
		})
	}

	return &model.GanttCompareData{
		ClassID:   classInfo.ID,
		ClassName: classInfo.Name,
		Bars:      bars,
	}, nil
}

func (svc *CurriculumService) CompareGantt(classPlanIDs []string) ([]model.GanttCompareData, error) {
	var result []model.GanttCompareData
	for _, id := range classPlanIDs {
		data, err := svc.GenerateGanttData(id)
		if err != nil {
			return nil, err
		}
		result = append(result, *data)
	}
	return result, nil
}

func (svc *CurriculumService) GetExcelReport(grade, semesterID string) ([]model.ExcelReportRow, error) {
	classes, err := svc.store.ListClasses()
	if err != nil {
		return nil, err
	}

	var rows []model.ExcelReportRow
	for _, c := range classes {
		if c.Grade != grade {
			continue
		}

		plans, err := svc.store.GetClassPlansByClass(c.ID)
		if err != nil {
			continue
		}

		for _, plan := range plans {
			if semesterID != "" && plan.SemesterID != semesterID {
				continue
			}
			summary, err := svc.CalculateProgressSummary(plan.ID)
			if err != nil {
				continue
			}

			planProgress := 0.0
			if summary.TotalLessons > 0 {
				planProgress = float64(summary.PlanLessonsByWeek) / float64(summary.TotalLessons) * 100
			}
			actualProgress := 0.0
			if summary.TotalLessons > 0 {
				actualProgress = float64(summary.CompletedLessons) / float64(summary.TotalLessons) * 100
			}

			rows = append(rows, model.ExcelReportRow{
				ClassName:     c.Name,
				Subject:       plan.Subject,
				PlanProgress:  planProgress,
				ActualProgress: actualProgress,
				Deviation:     summary.Deviation,
				IsUpToStandard: summary.KeyCoverageRate >= 0.8 && summary.ExamCoverageRate >= 0.8,
			})
		}
	}

	return rows, nil
}

func (svc *CurriculumService) RecordProgress(classPlanID, lessonID, teacherID string, recordType model.RecordType, actualWeek int, isSubstitute bool) (*model.ProgressRecord, error) {
	classPlan, err := svc.store.GetClassPlan(classPlanID)
	if err != nil {
		return nil, err
	}

	if classPlan.IsLocked && classPlan.Status != model.StatusExecuting {
		return nil, fmt.Errorf("计划状态不允许登记进度")
	}

	lastRecord, err := svc.store.GetLastProgressRecord(classPlanID)
	if err != nil {
		return nil, err
	}
	_ = lastRecord

	record := &model.ProgressRecord{
		ClassPlanID:  classPlanID,
		LessonID:     lessonID,
		RecordType:   recordType,
		ActualWeek:   actualWeek,
		TeacherID:    teacherID,
		IsSubstitute: isSubstitute,
	}

	if err := svc.store.AddProgressRecord(record); err != nil {
		return nil, err
	}

	return record, nil
}
