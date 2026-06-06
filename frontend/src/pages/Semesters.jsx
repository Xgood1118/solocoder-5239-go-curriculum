import React, { useState, useEffect } from 'react'
import { api } from '../api.js'

function Semesters() {
  const [semesters, setSemesters] = useState([])
  const [selectedSem, setSelectedSem] = useState(null)
  const [selectedUnit, setSelectedUnit] = useState(null)
  const [showSemModal, setShowSemModal] = useState(false)
  const [showUnitModal, setShowUnitModal] = useState(false)
  const [showLessonModal, setShowLessonModal] = useState(false)
  const [suggestion, setSuggestion] = useState(null)
  const [semForm, setSemForm] = useState({ name: '', grade: '高一年级', subject: '数学', total_weeks: 18 })
  const [unitForm, setUnitForm] = useState({ title: '', description: '', knowledge_point_ids: [], lesson_count: 6 })
  const [lessonForm, setLessonForm] = useState({ title: '', content: '', knowledge_point_ids: [], plan_week: 1 })
  const [kps, setKps] = useState([])
  const [editingUnit, setEditingUnit] = useState(null)
  const [editingLesson, setEditingLesson] = useState(null)

  useEffect(() => {
    loadSemesters()
    loadKps()
  }, [])

  const loadSemesters = async () => {
    try {
      const res = await api.listSemesters()
      setSemesters(res)
      if (res.length > 0 && !selectedSem) {
        setSelectedSem(res[0])
      }
    } catch (err) {
      alert(err.message)
    }
  }

  const loadKps = async () => {
    try {
      const res = await api.listKnowledgePoints()
      setKps(res)
    } catch (err) {
      console.error(err)
    }
  }

  const selectSem = async (sem) => {
    setSelectedSem(sem)
    setSelectedUnit(null)
    setSuggestion(null)
    try {
      const detail = await api.getSemester(sem.id)
      setSelectedSem(detail)
    } catch (err) {
      console.error(err)
    }
  }

  const selectUnit = async (unit) => {
    setSelectedUnit(unit)
    try {
      const sug = await api.suggestUnitLessons(selectedSem.id, unit.id)
      setSuggestion(sug)
    } catch (err) {
      console.error(err)
    }
  }

  const openUnitModal = (unit = null) => {
    setEditingUnit(unit)
    if (unit) {
      setUnitForm({ ...unit, knowledge_point_ids: unit.knowledge_point_ids || [] })
    } else {
      setUnitForm({ title: '', description: '', knowledge_point_ids: [], lesson_count: 6 })
    }
    setShowUnitModal(true)
  }

  const openLessonModal = (lesson = null) => {
    setEditingLesson(lesson)
    if (lesson) {
      setLessonForm({ ...lesson, knowledge_point_ids: lesson.knowledge_point_ids || [] })
    } else {
      setLessonForm({ title: '', content: '', knowledge_point_ids: [], plan_week: 1 })
    }
    setShowLessonModal(true)
  }

  const handleCreateSem = async () => {
    try {
      const id = 'sem_' + Date.now()
      await api.createSemester({ ...semForm, id })
      setShowSemModal(false)
      loadSemesters()
    } catch (err) {
      alert(err.message)
    }
  }

  const handleSaveUnit = async () => {
    try {
      if (editingUnit) {
        await api.updateUnit(selectedSem.id, editingUnit.id, unitForm)
      } else {
        const id = 'unit_' + Date.now()
        await api.createUnit(selectedSem.id, { ...unitForm, id })
      }
      setShowUnitModal(false)
      const detail = await api.getSemester(selectedSem.id)
      setSelectedSem(detail)
    } catch (err) {
      alert(err.message)
    }
  }

  const handleSaveLesson = async () => {
    try {
      if (editingLesson) {
        await api.updateLesson(selectedSem.id, selectedUnit.id, editingLesson.id, lessonForm)
      } else {
        const id = 'lesson_' + Date.now()
        await api.createLesson(selectedSem.id, selectedUnit.id, { ...lessonForm, id })
      }
      setShowLessonModal(false)
      const detail = await api.getSemester(selectedSem.id)
      setSelectedSem(detail)
      const unit = detail.units.find(u => u.id === selectedUnit.id)
      if (unit) setSelectedUnit(unit)
    } catch (err) {
      alert(err.message)
    }
  }

  const handleDeleteUnit = async (unitId) => {
    if (!confirm('确定删除此单元？')) return
    try {
      await api.deleteUnit(selectedSem.id, unitId)
      const detail = await api.getSemester(selectedSem.id)
      setSelectedSem(detail)
      setSelectedUnit(null)
    } catch (err) {
      alert(err.message)
    }
  }

  const handleDeleteLesson = async (lessonId) => {
    if (!confirm('确定删除此课时？')) return
    try {
      await api.deleteLesson(selectedSem.id, selectedUnit.id, lessonId)
      const detail = await api.getSemester(selectedSem.id)
      setSelectedSem(detail)
      const unit = detail.units.find(u => u.id === selectedUnit.id)
      if (unit) setSelectedUnit(unit)
    } catch (err) {
      alert(err.message)
    }
  }

  const toggleKP = (kpId, formType) => {
    const form = formType === 'unit' ? unitForm : lessonForm
    const setForm = formType === 'unit' ?
      (v) => setUnitForm(v) : (v) => setLessonForm(v)
    const list = form.knowledge_point_ids || []
    if (list.includes(kpId)) {
      setForm({ ...form, knowledge_point_ids: list.filter(id => id !== kpId) })
    } else {
      setForm({ ...form, knowledge_point_ids: [...list, kpId] })
    }
  }

  const statusTag = (status) => {
    const map = {
      draft: { text: '草稿', class: 'tag' },
      published: { text: '已发布', class: 'tag tag-blue' },
      executing: { text: '执行中', class: 'tag tag-green' },
      completed: { text: '已完成', class: 'tag tag-purple' },
      archived: { text: '已归档', class: 'tag' },
    }
    const info = map[status] || { text: status, class: 'tag' }
    return <span className={info.class}>{info.text}</span>
  }

  return (
    <div>
      <div className="page-header">
        <h1>教学大纲</h1>
        <button className="btn btn-primary" onClick={() => setShowSemModal(true)}>+ 新建学期</button>
      </div>

      <div style={{ display: 'flex', gap: 16, minHeight: 'calc(100vh - 120px)' }}>
        <div className="card" style={{ width: 260, flexShrink: 0 }}>
          <div className="card-title">学期列表</div>
          {semesters.map(sem => (
            <div
              key={sem.id}
              className={`unit-item ${selectedSem?.id === sem.id ? 'active' : ''}`}
              onClick={() => selectSem(sem)}
            >
              <div style={{ fontWeight: 500, marginBottom: 4 }}>{sem.name}</div>
              <div style={{ fontSize: 12, color: '#999' }}>
                {sem.grade} · {sem.subject} · {sem.units?.length || 0}单元
              </div>
              <div style={{ marginTop: 6 }}>{statusTag(sem.status)}</div>
            </div>
          ))}
          {semesters.length === 0 && <div className="empty-state">暂无学期</div>}
        </div>

        <div className="card" style={{ width: 320, flexShrink: 0 }}>
          {selectedSem ? (
            <>
              <div className="card-title flex-between">
                <span>单元列表</span>
                <button className="btn btn-primary" style={{ padding: '4px 12px', fontSize: 12 }} onClick={() => openUnitModal()}>+ 单元</button>
              </div>
              {selectedSem.units?.map((unit, idx) => (
                <div
                  key={unit.id}
                  className={`unit-item ${selectedUnit?.id === unit.id ? 'active' : ''}`}
                  onClick={() => selectUnit(unit)}
                >
                  <div style={{ fontWeight: 500, marginBottom: 4 }}>
                    第{idx + 1}单元 {unit.title}
                  </div>
                  <div style={{ fontSize: 12, color: '#999' }}>
                    {unit.lesson_count} 课时 · {unit.knowledge_point_ids?.length || 0}知识点
                  </div>
                </div>
              ))}
              {(!selectedSem.units || selectedSem.units.length === 0) && (
                <div className="empty-state">暂无单元</div>
              )}
            </>
          ) : (
            <div className="empty-state">请选择学期</div>
          )}
        </div>

        <div className="card" style={{ flex: 1 }}>
          {selectedUnit ? (
            <>
              <div className="card-title flex-between">
                <span>{selectedUnit.title}</span>
                <div style={{ display: 'flex', gap: 6 }}>
                  <button className="btn" style={{ padding: '4px 12px', fontSize: 12 }} onClick={() => openUnitModal(selectedUnit)}>编辑单元</button>
                  <button className="btn btn-danger" style={{ padding: '4px 12px', fontSize: 12 }} onClick={() => handleDeleteUnit(selectedUnit.id)}>删除</button>
                </div>
              </div>

              <p style={{ color: '#666', marginBottom: 16 }}>{selectedUnit.description || '暂无描述'}</p>

              {suggestion && (
                <div className="card" style={{ background: '#f0f9ff', border: '1px solid #bae7ff' }}>
                  <div style={{ fontWeight: 500, marginBottom: 6 }}>💡 智能课时建议</div>
                  <div style={{ fontSize: 13, color: '#555' }}>
                    建议 <strong>{suggestion.minLessons} - {suggestion.maxLessons}</strong> 课时
                    <div style={{ marginTop: 4, fontSize: 12, color: '#888' }}>{suggestion.reason}</div>
                  </div>
                </div>
              )}

              <div className="flex-between" style={{ marginBottom: 12 }}>
                <strong>课时列表（{selectedUnit.lessons?.length || 0}）</strong>
                <button className="btn btn-primary" style={{ padding: '4px 12px', fontSize: 12 }} onClick={() => openLessonModal()}>+ 课时</button>
              </div>

              {selectedUnit.lessons?.map((lesson, idx) => (
                <div key={lesson.id} className="lesson-item flex-between">
                  <div>
                    <span style={{ fontWeight: 500 }}>第{idx + 1}课 {lesson.title}</span>
                    <span style={{ fontSize: 12, color: '#999', marginLeft: 10 }}>
                      计划第{lesson.plan_week}周
                    </span>
                  </div>
                  <div>
                    <button className="btn" style={{ padding: '2px 10px', fontSize: 12, marginRight: 6 }} onClick={() => openLessonModal(lesson)}>编辑</button>
                    <button className="btn btn-danger" style={{ padding: '2px 10px', fontSize: 12 }} onClick={() => handleDeleteLesson(lesson.id)}>删除</button>
                  </div>
                </div>
              ))}
              {(!selectedUnit.lessons || selectedUnit.lessons.length === 0) && (
                <div className="empty-state">暂无课时</div>
              )}

              <div style={{ marginTop: 20 }}>
                <strong>单元知识点（{selectedUnit.knowledge_point_ids?.length || 0}）</strong>
                <div style={{ marginTop: 10, display: 'flex', flexWrap: 'wrap', gap: 6 }}>
                  {selectedUnit.knowledge_point_ids?.map(kpId => {
                    const kp = kps.find(k => k.id === kpId)
                    return kp ? <span key={kpId} className="tag tag-blue">{kp.name}</span> : null
                  })}
                </div>
              </div>
            </>
          ) : (
            <div className="empty-state">请选择单元查看详情</div>
          )}
        </div>
      </div>

      {showSemModal && (
        <div className="modal-overlay" onClick={() => setShowSemModal(false)}>
          <div className="modal" onClick={e => e.stopPropagation()}>
            <div className="modal-header"><h3>新建学期</h3><button className="close-btn" onClick={() => setShowSemModal(false)}>×</button></div>
            <div className="form-group">
              <label>学期名称</label>
              <input value={semForm.name} onChange={e => setSemForm({ ...semForm, name: e.target.value })} placeholder="如：2024-2025学年第一学期" />
            </div>
            <div className="form-group">
              <label>年级</label>
              <select value={semForm.grade} onChange={e => setSemForm({ ...semForm, grade: e.target.value })}>
                <option>高一年级</option><option>高二年级</option><option>高三年级</option>
                <option>初一年级</option><option>初二年级</option><option>初三年级</option>
              </select>
            </div>
            <div className="form-group">
              <label>科目</label>
              <select value={semForm.subject} onChange={e => setSemForm({ ...semForm, subject: e.target.value })}>
                <option>数学</option><option>语文</option><option>英语</option>
                <option>物理</option><option>化学</option><option>生物</option>
                <option>历史</option><option>地理</option><option>政治</option>
              </select>
            </div>
            <div className="form-group">
              <label>总周数</label>
              <input type="number" value={semForm.total_weeks} onChange={e => setSemForm({ ...semForm, total_weeks: parseInt(e.target.value) })} />
            </div>
            <div className="modal-footer">
              <button className="btn" onClick={() => setShowSemModal(false)}>取消</button>
              <button className="btn btn-primary" onClick={handleCreateSem}>创建</button>
            </div>
          </div>
        </div>
      )}

      {showUnitModal && (
        <div className="modal-overlay" onClick={() => setShowUnitModal(false)}>
          <div className="modal" onClick={e => e.stopPropagation()}>
            <div className="modal-header"><h3>{editingUnit ? '编辑单元' : '新增单元'}</h3><button className="close-btn" onClick={() => setShowUnitModal(false)}>×</button></div>
            <div className="form-group">
              <label>单元标题</label>
              <input value={unitForm.title} onChange={e => setUnitForm({ ...unitForm, title: e.target.value })} />
            </div>
            <div className="form-group">
              <label>描述</label>
              <textarea rows={2} value={unitForm.description || ''} onChange={e => setUnitForm({ ...unitForm, description: e.target.value })} />
            </div>
            <div className="form-group">
              <label>课时数</label>
              <input type="number" value={unitForm.lesson_count} onChange={e => setUnitForm({ ...unitForm, lesson_count: parseInt(e.target.value) })} />
            </div>
            <div className="form-group">
              <label>知识点</label>
              <div style={{ maxHeight: 150, overflowY: 'auto', border: '1px solid #e8e8e8', padding: 10, borderRadius: 4 }}>
                {kps.map(kp => (
                  <label key={kp.id} style={{ display: 'block', marginBottom: 4, cursor: 'pointer' }}>
                    <input type="checkbox"
                      checked={(unitForm.knowledge_point_ids || []).includes(kp.id)}
                      onChange={() => toggleKP(kp.id, 'unit')}
                    />
                    <span style={{ marginLeft: 6 }}>{kp.name}</span>
                  </label>
                ))}
              </div>
            </div>
            <div className="modal-footer">
              <button className="btn" onClick={() => setShowUnitModal(false)}>取消</button>
              <button className="btn btn-primary" onClick={handleSaveUnit}>保存</button>
            </div>
          </div>
        </div>
      )}

      {showLessonModal && (
        <div className="modal-overlay" onClick={() => setShowLessonModal(false)}>
          <div className="modal" onClick={e => e.stopPropagation()}>
            <div className="modal-header"><h3>{editingLesson ? '编辑课时' : '新增课时'}</h3><button className="close-btn" onClick={() => setShowLessonModal(false)}>×</button></div>
            <div className="form-group">
              <label>课时标题</label>
              <input value={lessonForm.title} onChange={e => setLessonForm({ ...lessonForm, title: e.target.value })} />
            </div>
            <div className="form-group">
              <label>内容</label>
              <textarea rows={3} value={lessonForm.content || ''} onChange={e => setLessonForm({ ...lessonForm, content: e.target.value })} />
            </div>
            <div className="form-group">
              <label>计划周次</label>
              <input type="number" min={1} value={lessonForm.plan_week} onChange={e => setLessonForm({ ...lessonForm, plan_week: parseInt(e.target.value) })} />
            </div>
            <div className="form-group">
              <label>知识点</label>
              <div style={{ maxHeight: 120, overflowY: 'auto', border: '1px solid #e8e8e8', padding: 10, borderRadius: 4 }}>
                {kps.map(kp => (
                  <label key={kp.id} style={{ display: 'block', marginBottom: 4, cursor: 'pointer' }}>
                    <input type="checkbox"
                      checked={(lessonForm.knowledge_point_ids || []).includes(kp.id)}
                      onChange={() => toggleKP(kp.id, 'lesson')}
                    />
                    <span style={{ marginLeft: 6 }}>{kp.name}</span>
                  </label>
                ))}
              </div>
            </div>
            <div className="modal-footer">
              <button className="btn" onClick={() => setShowLessonModal(false)}>取消</button>
              <button className="btn btn-primary" onClick={handleSaveLesson}>保存</button>
            </div>
          </div>
        </div>
      )}
    </div>
  )
}

export default Semesters
