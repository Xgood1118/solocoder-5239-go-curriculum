package handler

import (
	"curriculum/internal/model"
	"encoding/csv"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func (h *Handler) ListClasses(c *gin.Context) {
	classes, err := h.store.ListClasses()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, classes)
}

func (h *Handler) GetClass(c *gin.Context) {
	id := c.Param("id")
	class, err := h.store.GetClass(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, class)
}

func (h *Handler) CreateClass(c *gin.Context) {
	var class model.Class
	if err := c.ShouldBindJSON(&class); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.store.CreateClass(&class); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, class)
}

func (h *Handler) UpdateClass(c *gin.Context) {
	id := c.Param("id")
	var class model.Class
	if err := c.ShouldBindJSON(&class); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	class.ID = id
	if err := h.store.UpdateClass(&class); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, class)
}

func (h *Handler) DeleteClass(c *gin.Context) {
	id := c.Param("id")
	if err := h.store.DeleteClass(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

func (h *Handler) ListTeachers(c *gin.Context) {
	teachers, err := h.store.ListTeachers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, teachers)
}

func (h *Handler) GetTeacher(c *gin.Context) {
	id := c.Param("id")
	teacher, err := h.store.GetTeacher(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, teacher)
}

func (h *Handler) CreateTeacher(c *gin.Context) {
	var teacher model.Teacher
	if err := c.ShouldBindJSON(&teacher); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.store.CreateTeacher(&teacher); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, teacher)
}

func (h *Handler) UpdateTeacher(c *gin.Context) {
	id := c.Param("id")
	var teacher model.Teacher
	if err := c.ShouldBindJSON(&teacher); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	teacher.ID = id
	if err := h.store.UpdateTeacher(&teacher); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, teacher)
}

func (h *Handler) DeleteTeacher(c *gin.Context) {
	id := c.Param("id")
	if err := h.store.DeleteTeacher(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

func (h *Handler) ListClassPlans(c *gin.Context) {
	classID := c.Query("classId")
	if classID != "" {
		plans, err := h.store.GetClassPlansByClass(classID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, plans)
		return
	}
	plans, err := h.store.ListClassPlans()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, plans)
}

func (h *Handler) GetClassPlan(c *gin.Context) {
	id := c.Param("classPlanId")
	plan, err := h.store.GetClassPlan(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, plan)
}

func (h *Handler) CreateClassPlan(c *gin.Context) {
	var plan model.ClassPlan
	if err := c.ShouldBindJSON(&plan); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.store.CreateClassPlan(&plan); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, plan)
}

func (h *Handler) UpdateClassPlan(c *gin.Context) {
	id := c.Param("classPlanId")
	var plan model.ClassPlan
	if err := c.ShouldBindJSON(&plan); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	plan.ID = id
	if err := h.store.UpdateClassPlan(&plan); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, plan)
}

func (h *Handler) DeleteClassPlan(c *gin.Context) {
	id := c.Param("classPlanId")
	if err := h.store.DeleteClassPlan(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

func (h *Handler) GetProgressSummary(c *gin.Context) {
	classPlanID := c.Param("classPlanId")
	summary, err := h.svc.CalculateProgressSummary(classPlanID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, summary)
}

func (h *Handler) ListProgressRecords(c *gin.Context) {
	classPlanID := c.Param("classPlanId")
	records, err := h.store.ListProgressRecords(classPlanID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, records)
}

type RecordProgressRequest struct {
	LessonID     string            `json:"lesson_id"`
	TeacherID    string            `json:"teacher_id"`
	RecordType   model.RecordType  `json:"record_type"`
	ActualWeek   int               `json:"actual_week"`
	IsSubstitute bool              `json:"is_substitute"`
	SubstituteNote string          `json:"substitute_note"`
	Notes        string            `json:"notes"`
}

func (h *Handler) RecordProgress(c *gin.Context) {
	classPlanID := c.Param("classPlanId")
	var req RecordProgressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.RecordType == "" {
		req.RecordType = model.RecordNormal
	}

	record, err := h.svc.RecordProgress(classPlanID, req.LessonID, req.TeacherID, req.RecordType, req.ActualWeek, req.IsSubstitute)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	record.SubstituteNote = req.SubstituteNote
	record.Notes = req.Notes
	c.JSON(http.StatusCreated, record)
}

func (h *Handler) QuickRecord(c *gin.Context) {
	classPlanID := c.Param("classPlanId")
	teacherID := c.Query("teacherId")
	lessonID := c.Query("lessonId")
	weekStr := c.DefaultQuery("week", "1")
	week, _ := strconv.Atoi(weekStr)
	isSub := c.DefaultQuery("isSubstitute", "false") == "true"

	record, err := h.svc.RecordProgress(classPlanID, lessonID, teacherID, model.RecordNormal, week, isSub)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, record)
}

func (h *Handler) GetCoverageReport(c *gin.Context) {
	classPlanID := c.Param("classPlanId")
	report, err := h.svc.CalculateCoverageReport(classPlanID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, report)
}

func (h *Handler) GetGanttData(c *gin.Context) {
	classPlanID := c.Param("classPlanId")
	data, err := h.svc.GenerateGanttData(classPlanID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, data)
}

func (h *Handler) CompareGantt(c *gin.Context) {
	classPlanIDs := c.QueryArray("classPlanIds")
	if len(classPlanIDs) < 2 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "至少需要两个班级计划进行对比"})
		return
	}
	data, err := h.svc.CompareGantt(classPlanIDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, data)
}

func (h *Handler) ExportExcel(c *gin.Context) {
	grade := c.DefaultQuery("grade", "高一年级")
	semesterID := c.Query("semesterId")

	rows, err := h.svc.GetExcelReport(grade, semesterID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Header("Content-Type", "text/csv; charset=utf-8")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s_期中课程进度汇总表.csv", grade))
	c.Writer.WriteString("\xEF\xBB\xBF")

	writer := csv.NewWriter(c.Writer)
	defer writer.Flush()

	header := []string{"班级", "科目", "计划进度(%)", "实际进度(%)", "偏差(课时)", "是否达标"}
	writer.Write(header)

	for _, row := range rows {
		status := "达标"
		if !row.IsUpToStandard {
			status = "未达标"
		}
		record := []string{
			row.ClassName,
			row.Subject,
			fmt.Sprintf("%.1f", row.PlanProgress),
			fmt.Sprintf("%.1f", row.ActualProgress),
			fmt.Sprintf("%d", row.Deviation),
			status,
		}
		writer.Write(record)
	}
}

func (h *Handler) ListRevisionRequests(c *gin.Context) {
	reqs, err := h.store.ListRevisionRequests()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, reqs)
}

func (h *Handler) CreateRevisionRequest(c *gin.Context) {
	var req model.RevisionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.store.CreateRevisionRequest(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, req)
}

func (h *Handler) ApproveRevision(c *gin.Context) {
	id := c.Param("id")
	var approvelReq struct {
		ApproverID   string `json:"approver_id"`
		ApprovalNote string `json:"approval_note"`
		NewPlanID    string `json:"new_plan_id"`
	}
	if err := c.ShouldBindJSON(&approvelReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req, err := h.store.GetRevisionRequest(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	now := time.Now()
	req.Status = model.RevisionApproved
	req.ApproverID = approvelReq.ApproverID
	req.ApprovalNote = approvelReq.ApprovalNote
	req.NewPlanID = approvelReq.NewPlanID
	req.ApprovedAt = &now

	if err := h.store.UpdateRevisionRequest(req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	origPlan, err := h.store.GetSemester(req.OriginalPlanID)
	if err == nil {
		origPlan.Status = model.StatusArchived
		h.store.UpdateSemester(origPlan)
	}

	newPlan, err := h.store.GetSemester(approvelReq.NewPlanID)
	if err == nil {
		newPlan.Status = model.StatusPublished
		newPlan.ParentPlanID = req.OriginalPlanID
		h.store.UpdateSemester(newPlan)
	}

	c.JSON(http.StatusOK, req)
}

func (h *Handler) RejectRevision(c *gin.Context) {
	id := c.Param("id")
	var rejectReq struct {
		ApproverID   string `json:"approver_id"`
		ApprovalNote string `json:"approval_note"`
	}
	if err := c.ShouldBindJSON(&rejectReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req, err := h.store.GetRevisionRequest(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	req.Status = model.RevisionRejected
	req.ApproverID = rejectReq.ApproverID
	req.ApprovalNote = rejectReq.ApprovalNote

	if err := h.store.UpdateRevisionRequest(req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, req)
}

func (h *Handler) ListKnowledgePoints(c *gin.Context) {
	kps, err := h.store.ListKnowledgePoints()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, kps)
}

func (h *Handler) GetKnowledgePoint(c *gin.Context) {
	id := c.Param("id")
	kp, err := h.store.GetKnowledgePoint(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, kp)
}

func (h *Handler) CreateKnowledgePoint(c *gin.Context) {
	var kp model.KnowledgePoint
	if err := c.ShouldBindJSON(&kp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.store.CreateKnowledgePoint(&kp); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, kp)
}

func (h *Handler) UpdateKnowledgePoint(c *gin.Context) {
	id := c.Param("id")
	var kp model.KnowledgePoint
	if err := c.ShouldBindJSON(&kp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	kp.ID = id
	if err := h.store.UpdateKnowledgePoint(&kp); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, kp)
}

func (h *Handler) DeleteKnowledgePoint(c *gin.Context) {
	id := c.Param("id")
	if err := h.store.DeleteKnowledgePoint(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

func (h *Handler) TopoSort(c *gin.Context) {
	result, err := h.store.TopologicalSort()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *Handler) Dashboard(c *gin.Context) {
	semesters, _ := h.store.ListSemesters()
	classes, _ := h.store.ListClasses()
	teachers, _ := h.store.ListTeachers()
	kps, _ := h.store.ListKnowledgePoints()
	topoResult, _ := h.store.TopologicalSort()

	var criticalDeviations []gin.H
	classPlans, _ := h.store.ListClassPlans()
	for _, cp := range classPlans {
		summary, err := h.svc.CalculateProgressSummary(cp.ID)
		if err != nil {
			continue
		}
		if summary.DeviationLevel == "critical" || summary.DeviationLevel == "warning" {
			classInfo, _ := h.store.GetClass(cp.ClassID)
			className := cp.ClassID
			if classInfo != nil {
				className = classInfo.Name
			}
			criticalDeviations = append(criticalDeviations, gin.H{
				"class_plan_id": cp.ID,
				"class_name":    className,
				"subject":       cp.Subject,
				"deviation":     summary.Deviation,
				"level":         summary.DeviationLevel,
			})
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"semester_count":       len(semesters),
		"class_count":          len(classes),
		"teacher_count":        len(teachers),
		"knowledge_point_count": len(kps),
		"topo_has_cycle":       topoResult != nil && topoResult.Cycle != nil && topoResult.Cycle.HasCycle,
		"cycle_info":           topoResult.Cycle,
		"critical_deviations":  criticalDeviations,
	})
}
