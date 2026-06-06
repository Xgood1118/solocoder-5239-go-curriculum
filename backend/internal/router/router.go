package router

import (
	"curriculum/internal/handler"

	"github.com/gin-gonic/gin"
)

func SetupRouter(h *handler.Handler) *gin.Engine {
	r := gin.Default()

	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	api := r.Group("/api")
	{
		api.GET("/dashboard", h.Dashboard)

		kp := api.Group("/knowledge-points")
		{
			kp.GET("", h.ListKnowledgePoints)
			kp.POST("", h.CreateKnowledgePoint)
			kp.GET("/topo-sort", h.TopoSort)
			kp.GET("/:id", h.GetKnowledgePoint)
			kp.PUT("/:id", h.UpdateKnowledgePoint)
			kp.DELETE("/:id", h.DeleteKnowledgePoint)
		}

		semesters := api.Group("/semesters")
		{
			semesters.GET("", h.ListSemesters)
			semesters.POST("", h.CreateSemester)
			semesters.GET("/:id", h.GetSemester)
			semesters.PUT("/:id", h.UpdateSemester)
			semesters.DELETE("/:id", h.DeleteSemester)

			units := semesters.Group("/:semesterId/units")
			{
				units.GET("", h.ListUnits)
				units.POST("", h.CreateUnit)
				units.GET("/:unitId/suggest-lessons", h.SuggestUnitLessons)
				units.PUT("/:unitId", h.UpdateUnit)
				units.DELETE("/:unitId", h.DeleteUnit)

				lessons := units.Group("/:unitId/lessons")
				{
					lessons.GET("", h.ListLessons)
					lessons.POST("", h.CreateLesson)
					lessons.PUT("/:lessonId", h.UpdateLesson)
					lessons.DELETE("/:lessonId", h.DeleteLesson)
				}
			}
		}

		classes := api.Group("/classes")
		{
			classes.GET("", h.ListClasses)
			classes.POST("", h.CreateClass)
			classes.GET("/:id", h.GetClass)
			classes.PUT("/:id", h.UpdateClass)
			classes.DELETE("/:id", h.DeleteClass)
		}

		teachers := api.Group("/teachers")
		{
			teachers.GET("", h.ListTeachers)
			teachers.POST("", h.CreateTeacher)
			teachers.GET("/:id", h.GetTeacher)
			teachers.PUT("/:id", h.UpdateTeacher)
			teachers.DELETE("/:id", h.DeleteTeacher)
		}

		classPlans := api.Group("/class-plans")
		{
			classPlans.GET("", h.ListClassPlans)
			classPlans.POST("", h.CreateClassPlan)
			classPlans.GET("/:id", h.GetClassPlan)
			classPlans.PUT("/:id", h.UpdateClassPlan)
			classPlans.DELETE("/:id", h.DeleteClassPlan)

			progress := classPlans.Group("/:classPlanId/progress")
			{
				progress.GET("", h.GetProgressSummary)
				progress.GET("/records", h.ListProgressRecords)
				progress.POST("/record", h.RecordProgress)
				progress.POST("/quick", h.QuickRecord)
			}

			classPlans.GET("/:classPlanId/coverage", h.GetCoverageReport)
			classPlans.GET("/:classPlanId/gantt", h.GetGanttData)
		}

		api.GET("/gantt/compare", h.CompareGantt)
		api.GET("/reports/excel", h.ExportExcel)

		revisions := api.Group("/revision-requests")
		{
			revisions.GET("", h.ListRevisionRequests)
			revisions.POST("", h.CreateRevisionRequest)
			revisions.POST("/:id/approve", h.ApproveRevision)
			revisions.POST("/:id/reject", h.RejectRevision)
		}
	}

	return r
}
