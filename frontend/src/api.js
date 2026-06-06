const API_BASE = '/api'

async function request(url, options = {}) {
  const res = await fetch(`${API_BASE}${url}`, {
    headers: {
      'Content-Type': 'application/json',
      ...options.headers,
    },
    ...options,
  })
  if (!res.ok) {
    const err = await res.json().catch(() => ({ error: '请求失败' }))
    throw new Error(err.error || '请求失败')
  }
  if (res.headers.get('content-type')?.includes('text/csv')) {
    return res
  }
  return res.json()
}

export const api = {
  getDashboard: () => request('/dashboard'),

  listKnowledgePoints: () => request('/knowledge-points'),
  getKnowledgePoint: (id) => request(`/knowledge-points/${id}`),
  createKnowledgePoint: (data) => request('/knowledge-points', { method: 'POST', body: JSON.stringify(data) }),
  updateKnowledgePoint: (id, data) => request(`/knowledge-points/${id}`, { method: 'PUT', body: JSON.stringify(data) }),
  deleteKnowledgePoint: (id) => request(`/knowledge-points/${id}`, { method: 'DELETE' }),
  topoSort: () => request('/knowledge-points/topo-sort'),

  listSemesters: () => request('/semesters'),
  getSemester: (id) => request(`/semesters/${id}`),
  createSemester: (data) => request('/semesters', { method: 'POST', body: JSON.stringify(data) }),
  updateSemester: (id, data) => request(`/semesters/${id}`, { method: 'PUT', body: JSON.stringify(data) }),
  deleteSemester: (id) => request(`/semesters/${id}`, { method: 'DELETE' }),

  listUnits: (semesterId) => request(`/semesters/${semesterId}/units`),
  createUnit: (semesterId, data) => request(`/semesters/${semesterId}/units`, { method: 'POST', body: JSON.stringify(data) }),
  updateUnit: (semesterId, unitId, data) => request(`/semesters/${semesterId}/units/${unitId}`, { method: 'PUT', body: JSON.stringify(data) }),
  deleteUnit: (semesterId, unitId) => request(`/semesters/${semesterId}/units/${unitId}`, { method: 'DELETE' }),
  suggestUnitLessons: (semesterId, unitId) => request(`/semesters/${semesterId}/units/${unitId}/suggest-lessons`),

  listLessons: (semesterId, unitId) => request(`/semesters/${semesterId}/units/${unitId}/lessons`),
  createLesson: (semesterId, unitId, data) => request(`/semesters/${semesterId}/units/${unitId}/lessons`, { method: 'POST', body: JSON.stringify(data) }),
  updateLesson: (semesterId, unitId, lessonId, data) => request(`/semesters/${semesterId}/units/${unitId}/lessons/${lessonId}`, { method: 'PUT', body: JSON.stringify(data) }),
  deleteLesson: (semesterId, unitId, lessonId) => request(`/semesters/${semesterId}/units/${unitId}/lessons/${lessonId}`, { method: 'DELETE' }),

  listClasses: () => request('/classes'),
  createClass: (data) => request('/classes', { method: 'POST', body: JSON.stringify(data) }),
  updateClass: (id, data) => request(`/classes/${id}`, { method: 'PUT', body: JSON.stringify(data) }),
  deleteClass: (id) => request(`/classes/${id}`, { method: 'DELETE' }),

  listTeachers: () => request('/teachers'),
  createTeacher: (data) => request('/teachers', { method: 'POST', body: JSON.stringify(data) }),
  updateTeacher: (id, data) => request(`/teachers/${id}`, { method: 'PUT', body: JSON.stringify(data) }),
  deleteTeacher: (id) => request(`/teachers/${id}`, { method: 'DELETE' }),

  listClassPlans: (classId) => request(`/class-plans${classId ? `?classId=${classId}` : ''}`),
  getClassPlan: (id) => request(`/class-plans/${id}`),
  createClassPlan: (data) => request('/class-plans', { method: 'POST', body: JSON.stringify(data) }),
  updateClassPlan: (id, data) => request(`/class-plans/${id}`, { method: 'PUT', body: JSON.stringify(data) }),

  getProgressSummary: (classPlanId) => request(`/class-plans/${classPlanId}/progress`),
  listProgressRecords: (classPlanId) => request(`/class-plans/${classPlanId}/progress/records`),
  recordProgress: (classPlanId, data) => request(`/class-plans/${classPlanId}/progress/record`, { method: 'POST', body: JSON.stringify(data) }),
  quickRecord: (classPlanId, params) => {
    const qs = new URLSearchParams(params).toString()
    return request(`/class-plans/${classPlanId}/progress/quick?${qs}`, { method: 'POST' })
  },

  getCoverageReport: (classPlanId) => request(`/class-plans/${classPlanId}/coverage`),
  getGanttData: (classPlanId) => request(`/class-plans/${classPlanId}/gantt`),
  compareGantt: (classPlanIds) => {
    const qs = classPlanIds.map(id => `classPlanIds=${id}`).join('&')
    return request(`/gantt/compare?${qs}`)
  },

  exportExcel: (grade, semesterId) => {
    const qs = new URLSearchParams({ grade, ...(semesterId && { semesterId }) }).toString()
    return request(`/reports/excel?${qs}`)
  },

  listRevisionRequests: () => request('/revision-requests'),
  createRevisionRequest: (data) => request('/revision-requests', { method: 'POST', body: JSON.stringify(data) }),
  approveRevision: (id, data) => request(`/revision-requests/${id}/approve`, { method: 'POST', body: JSON.stringify(data) }),
  rejectRevision: (id, data) => request(`/revision-requests/${id}/reject`, { method: 'POST', body: JSON.stringify(data) }),
}
