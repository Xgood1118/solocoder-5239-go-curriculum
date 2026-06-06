import React, { useState, useEffect } from 'react'
import { api } from '../api.js'

function Progress() {
  const [classes, setClasses] = useState([])
  const [classPlans, setClassPlans] = useState([])
  const [selectedPlan, setSelectedPlan] = useState(null)
  const [summary, setSummary] = useState(null)
  const [records, setRecords] = useState([])
  const [coverage, setCoverage] = useState(null)
  const [lessons, setLessons] = useState([])
  const [semester, setSemester] = useState(null)
  const [showRecordModal, setShowRecordModal] = useState(false)
  const [recordForm, setRecordForm] = useState({
    lesson_id: '',
    teacher_id: '',
    record_type: 'normal',
    actual_week: 1,
    is_substitute: false,
    notes: '',
  })
  const [teachers, setTeachers] = useState([])
  const [activeTab, setActiveTab] = useState('summary')

  useEffect(() => {
    loadClasses()
    loadTeachers()
  }, [])

  const loadClasses = async () => {
    try {
      const res = await api.listClasses()
      setClasses(res)
    } catch (err) { console.error(err) }
  }

  const loadTeachers = async () => {
    try {
      const res = await api.listTeachers()
      setTeachers(res)
    } catch (err) { console.error(err) }
  }

  const onSelectClass = async (classId) => {
    try {
      const res = await api.listClassPlans(classId)
      setClassPlans(res)
      setSelectedPlan(null)
      setSummary(null)
    } catch (err) { console.error(err) }
  }

  const onSelectPlan = async (plan) => {
    setSelectedPlan(plan)
    try {
      const [sum, recs, cov, sem] = await Promise.all([
        api.getProgressSummary(plan.id),
        api.listProgressRecords(plan.id),
        api.getCoverageReport(plan.id),
        api.getSemester(plan.semester_id),
      ])
      setSummary(sum)
      setRecords(recs)
      setCoverage(cov)
      setSemester(sem)
      const allLessons = []
      sem.units?.forEach(u => u.lessons?.forEach(l => allLessons.push(l)))
      setLessons(allLessons)
    } catch (err) {
      console.error(err)
    }
  }

  const handleRecord = async () => {
    try {
      await api.recordProgress(selectedPlan.id, recordForm)
      setShowRecordModal(false)
      onSelectPlan(selectedPlan)
    } catch (err) {
      alert(err.message)
    }
  }

  const deviationColor = (level) => {
    if (level === 'critical') return 'danger-text'
    if (level === 'warning') return 'warning-text'
    if (level === 'ahead') return 'success-text'
    return ''
  }

  const deviationText = (level) => {
    const map = { normal: '正常', warning: '滞后预警', critical: '严重滞后', ahead: '超前' }
    return map[level] || level
  }

  return (
    <div>
      <div className="page-header">
        <h1>进度跟踪</h1>
        <button className="btn btn-primary" onClick={() => setShowRecordModal(true)} disabled={!selectedPlan}>
          📝 登记进度
        </button>
      </div>

      <div style={{ display: 'flex', gap: 16, minHeight: 'calc(100vh - 120px)' }}>
        <div className="card" style={{ width: 220, flexShrink: 0 }}>
          <div className="card-title">班级</div>
          {classes.map(c => (
            <div key={c.id} className="unit-item" onClick={() => onSelectClass(c.id)}>
              <div style={{ fontWeight: 500 }}>{c.name}</div>
              <div style={{ fontSize: 12, color: '#999' }}>{c.grade} · {c.student_count}人</div>
            </div>
          ))}
        </div>

        <div className="card" style={{ width: 260, flexShrink: 0 }}>
          <div className="card-title">教学计划</div>
          {classPlans.map(p => (
            <div key={p.id}
              className={`unit-item ${selectedPlan?.id === p.id ? 'active' : ''}`}
              onClick={() => onSelectPlan(p)}
            >
              <div style={{ fontWeight: 500 }}>{p.subject}</div>
              <div style={{ fontSize: 12, color: '#999' }}>
                状态：{p.status === 'executing' ? '执行中' : p.status}
                {p.is_locked ? ' · 已锁定' : ''}
              </div>
            </div>
          ))}
          {classPlans.length === 0 && <div className="empty-state">暂无计划</div>}
        </div>

        <div className="card" style={{ flex: 1 }}>
          {selectedPlan ? (
            <>
              <div className="tabs">
                <div className={`tab ${activeTab === 'summary' ? 'active' : ''}`} onClick={() => setActiveTab('summary')}>进度概览</div>
                <div className={`tab ${activeTab === 'records' ? 'active' : ''}`} onClick={() => setActiveTab('records')}>登记记录</div>
                <div className={`tab ${activeTab === 'coverage' ? 'active' : ''}`} onClick={() => setActiveTab('coverage')}>覆盖率</div>
              </div>

              {activeTab === 'summary' && summary && (
                <div>
                  <div className="stat-grid">
                    <div className="stat-card">
                      <div className="label">总课时</div>
                      <div className="value primary">{summary.total_lessons}</div>
                    </div>
                    <div className="stat-card">
                      <div className="label">已完成</div>
                      <div className="value success">{summary.completed_lessons}</div>
                    </div>
                    <div className="stat-card">
                      <div className="label">计划进度</div>
                      <div className="value">{summary.plan_lessons_by_week}</div>
                    </div>
                    <div className="stat-card">
                      <div className="label">偏差</div>
                      <div className={`value ${deviationColor(summary.deviation_level)}`}>
                        {summary.deviation > 0 ? '+' : ''}{summary.deviation}
                      </div>
                    </div>
                  </div>

                  <div className="card">
                    <div className="card-title">进度状态</div>
                    <p>
                      当前第 <strong>{summary.current_week}</strong> 周 ·
                      <span className={deviationColor(summary.deviation_level)} style={{ marginLeft: 8 }}>
                        {deviationText(summary.deviation_level)}
                      </span>
                    </p>
                    <div className="progress-bar mt-10">
                      <div
                        className={`fill ${summary.deviation_level === 'critical' ? 'danger' : summary.deviation_level === 'warning' ? 'warning' : ''}`}
                        style={{ width: `${(summary.completed_lessons / summary.total_lessons) * 100}%` }}
                      />
                    </div>
                    <p style={{ marginTop: 8, fontSize: 12, color: '#999' }}>
                      完成度：{((summary.completed_lessons / summary.total_lessons) * 100).toFixed(1)}%
                    </p>
                  </div>

                  <div className="card">
                    <div className="card-title">覆盖率</div>
                    <div style={{ marginBottom: 12 }}>
                      <div className="flex-between" style={{ marginBottom: 4 }}>
                        <span>重点覆盖率</span>
                        <span className={summary.key_coverage_rate >= 0.8 ? 'success-text' : 'danger-text'}>
                          {(summary.key_coverage_rate * 100).toFixed(1)}%
                        </span>
                      </div>
                      <div className="progress-bar">
                        <div className="fill" style={{ width: `${summary.key_coverage_rate * 100}%` }} />
                      </div>
                    </div>
                    <div style={{ marginBottom: 12 }}>
                      <div className="flex-between" style={{ marginBottom: 4 }}>
                        <span>难点覆盖率</span>
                        <span className={summary.difficult_coverage_rate >= 0.8 ? 'success-text' : 'danger-text'}>
                          {(summary.difficult_coverage_rate * 100).toFixed(1)}%
                        </span>
                      </div>
                      <div className="progress-bar">
                        <div className="fill" style={{ width: `${summary.difficult_coverage_rate * 100}%` }} />
                      </div>
                    </div>
                    <div>
                      <div className="flex-between" style={{ marginBottom: 4 }}>
                        <span>考点覆盖率</span>
                        <span className={summary.exam_coverage_rate >= 0.8 ? 'success-text' : 'danger-text'}>
                          {(summary.exam_coverage_rate * 100).toFixed(1)}%
                        </span>
                      </div>
                      <div className="progress-bar">
                        <div className="fill" style={{ width: `${summary.exam_coverage_rate * 100}%` }} />
                      </div>
                    </div>
                  </div>
                </div>
              )}

              {activeTab === 'records' && (
                <div>
                  <table>
                    <thead>
                      <tr>
                        <th>时间</th>
                        <th>课时</th>
                        <th>类型</th>
                        <th>教师</th>
                        <th>周次</th>
                        <th>备注</th>
                      </tr>
                    </thead>
                    <tbody>
                      {records.map(r => {
                        const lesson = lessons.find(l => l.id === r.lesson_id)
                        const teacher = teachers.find(t => t.id === r.teacher_id)
                        const typeMap = {
                          normal: { text: '正常', class: 'tag tag-green' },
                          substitute: { text: '代课', class: 'tag tag-orange' },
                          study_self: { text: '自习', class: 'tag' },
                          pe: { text: '体育', class: 'tag tag-blue' },
                          meeting: { text: '会议', class: 'tag' },
                          occupied: { text: '占用', class: 'tag tag-red' },
                        }
                        const typeInfo = typeMap[r.record_type] || { text: r.record_type, class: 'tag' }
                        return (
                          <tr key={r.id}>
                            <td>{new Date(r.actual_date).toLocaleDateString()}</td>
                            <td>{lesson?.title || r.lesson_id}</td>
                            <td>
                              <span className={typeInfo.class}>{typeInfo.text}</span>
                              {r.is_substitute && <span className="tag tag-orange" style={{ marginLeft: 4 }}>代课</span>}
                            </td>
                            <td>{teacher?.name || r.teacher_id}</td>
                            <td>第{r.actual_week}周</td>
                            <td style={{ fontSize: 12, color: '#999' }}>{r.notes || '-'}</td>
                          </tr>
                        )
                      })}
                    </tbody>
                  </table>
                  {records.length === 0 && <div className="empty-state">暂无登记记录</div>}
                </div>
              )}

              {activeTab === 'coverage' && coverage && (
                <div>
                  {coverage.unpassed_units?.length > 0 && (
                    <div className="card" style={{ borderLeft: '4px solid #ff4d4f' }}>
                      <strong style={{ color: '#ff4d4f' }}>
                        ⚠️ 未达标章节（{coverage.unpassed_units.length}）
                      </strong>
                      <div style={{ marginTop: 8 }}>
                        {coverage.unpassed_units.map(uid => {
                          const unit = semester?.units?.find(u => u.id === uid)
                          return <span key={uid} className="tag tag-red" style={{ marginRight: 6 }}>{unit?.title || uid}</span>
                        })}
                      </div>
                    </div>
                  )}

                  <div className="card">
                    <div className="card-title">重点覆盖率</div>
                    {coverage.key_points?.map(item => (
                      <div key={item.unit_id} style={{ marginBottom: 10 }}>
                        <div className="flex-between">
                          <span>{item.unit_name}</span>
                          <span className={item.is_passed ? 'success-text' : 'danger-text'}>
                            {(item.rate * 100).toFixed(1)}%
                          </span>
                        </div>
                        <div className="progress-bar" style={{ marginTop: 4 }}>
                          <div className={`fill ${!item.is_passed ? 'danger' : ''}`} style={{ width: `${item.rate * 100}%` }} />
                        </div>
                      </div>
                    ))}
                  </div>
                </div>
              )}
            </>
          ) : (
            <div className="empty-state">请选择班级和教学计划</div>
          )}
        </div>
      </div>

      {showRecordModal && (
        <div className="modal-overlay" onClick={() => setShowRecordModal(false)}>
          <div className="modal" onClick={e => e.stopPropagation()}>
            <div className="modal-header">
              <h3>登记进度</h3>
              <button className="close-btn" onClick={() => setShowRecordModal(false)}>×</button>
            </div>

            <div className="form-group">
              <label>课时</label>
              <select value={recordForm.lesson_id} onChange={e => setRecordForm({ ...recordForm, lesson_id: e.target.value })}>
                <option value="">请选择课时</option>
                {lessons.map(l => (
                  <option key={l.id} value={l.id}>{l.title}</option>
                ))}
              </select>
            </div>

            <div className="form-group">
              <label>登记类型</label>
              <select value={recordForm.record_type} onChange={e => setRecordForm({ ...recordForm, record_type: e.target.value })}>
                <option value="normal">正常上课</option>
                <option value="substitute">代课</option>
                <option value="study_self">自习课</option>
                <option value="pe">体育课</option>
                <option value="meeting">会议占用</option>
                <option value="occupied">其他占用</option>
              </select>
            </div>

            <div className="form-group">
              <label>教师</label>
              <select value={recordForm.teacher_id} onChange={e => setRecordForm({ ...recordForm, teacher_id: e.target.value })}>
                <option value="">请选择教师</option>
                {teachers.map(t => (
                  <option key={t.id} value={t.id}>{t.name}（{t.subject}）</option>
                ))}
              </select>
            </div>

            <div className="form-group">
              <label>实际周次</label>
              <input type="number" min={1} value={recordForm.actual_week}
                onChange={e => setRecordForm({ ...recordForm, actual_week: parseInt(e.target.value) })} />
            </div>

            <div className="form-group">
              <label style={{ display: 'flex', alignItems: 'center', gap: 6 }}>
                <input type="checkbox"
                  checked={recordForm.is_substitute}
                  onChange={e => setRecordForm({ ...recordForm, is_substitute: e.target.checked })}
                />
                标记为代课
              </label>
            </div>

            <div className="form-group">
              <label>备注</label>
              <textarea rows={2} value={recordForm.notes || ''}
                onChange={e => setRecordForm({ ...recordForm, notes: e.target.value })} />
            </div>

            <div className="modal-footer">
              <button className="btn" onClick={() => setShowRecordModal(false)}>取消</button>
              <button className="btn btn-primary" onClick={handleRecord}>确认登记</button>
            </div>
          </div>
        </div>
      )}
    </div>
  )
}

export default Progress
