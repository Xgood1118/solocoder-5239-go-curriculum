import React, { useState, useEffect } from 'react'
import { api } from '../api.js'

function Teachers() {
  const [teachers, setTeachers] = useState([])
  const [showModal, setShowModal] = useState(false)
  const [editing, setEditing] = useState(null)
  const [form, setForm] = useState({ name: '', subject: '数学', class_ids: [] })
  const [classes, setClasses] = useState([])

  useEffect(() => {
    loadData()
    loadClasses()
  }, [])

  const loadData = async () => {
    try {
      const res = await api.listTeachers()
      setTeachers(res)
    } catch (err) { alert(err.message) }
  }

  const loadClasses = async () => {
    try {
      const res = await api.listClasses()
      setClasses(res)
    } catch (err) { console.error(err) }
  }

  const openCreate = () => {
    setEditing(null)
    setForm({ name: '', subject: '数学', class_ids: [] })
    setShowModal(true)
  }

  const openEdit = (t) => {
    setEditing(t)
    setForm({ ...t, class_ids: t.class_ids || [] })
    setShowModal(true)
  }

  const toggleClass = (classId) => {
    const list = form.class_ids || []
    if (list.includes(classId)) {
      setForm({ ...form, class_ids: list.filter(id => id !== classId) })
    } else {
      setForm({ ...form, class_ids: [...list, classId] })
    }
  }

  const handleSubmit = async () => {
    try {
      if (editing) {
        await api.updateTeacher(editing.id, form)
      } else {
        await api.createTeacher({ ...form, id: 'tch_' + Date.now() })
      }
      setShowModal(false)
      loadData()
    } catch (err) { alert(err.message) }
  }

  const handleDelete = async (id) => {
    if (!confirm('确定删除？')) return
    try {
      await api.deleteTeacher(id)
      loadData()
    } catch (err) { alert(err.message) }
  }

  return (
    <div>
      <div className="page-header">
        <h1>教师管理</h1>
        <button className="btn btn-primary" onClick={openCreate}>+ 新增教师</button>
      </div>

      <div className="card">
        <table>
          <thead>
            <tr>
              <th>姓名</th>
              <th>科目</th>
              <th>任教班级</th>
              <th>操作</th>
            </tr>
          </thead>
          <tbody>
            {teachers.map(t => (
              <tr key={t.id}>
                <td><strong>{t.name}</strong></td>
                <td>{t.subject}</td>
                <td style={{ fontSize: 13, color: '#666' }}>
                  {t.class_ids?.length || 0} 个班
                </td>
                <td>
                  <button className="btn" style={{ marginRight: 6 }} onClick={() => openEdit(t)}>编辑</button>
                  <button className="btn btn-danger" onClick={() => handleDelete(t.id)}>删除</button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
        {teachers.length === 0 && <div className="empty-state">暂无教师</div>}
      </div>

      {showModal && (
        <div className="modal-overlay" onClick={() => setShowModal(false)}>
          <div className="modal" onClick={e => e.stopPropagation()}>
            <div className="modal-header">
              <h3>{editing ? '编辑教师' : '新增教师'}</h3>
              <button className="close-btn" onClick={() => setShowModal(false)}>×</button>
            </div>
            <div className="form-group">
              <label>姓名</label>
              <input value={form.name} onChange={e => setForm({ ...form, name: e.target.value })} />
            </div>
            <div className="form-group">
              <label>科目</label>
              <select value={form.subject} onChange={e => setForm({ ...form, subject: e.target.value })}>
                <option>数学</option><option>语文</option><option>英语</option>
                <option>物理</option><option>化学</option><option>生物</option>
                <option>历史</option><option>地理</option><option>政治</option>
                <option>体育</option>
              </select>
            </div>
            <div className="form-group">
              <label>任教班级</label>
              <div style={{ maxHeight: 150, overflowY: 'auto', border: '1px solid #e8e8e8', padding: 10, borderRadius: 4 }}>
                {classes.map(c => (
                  <label key={c.id} style={{ display: 'block', marginBottom: 4, cursor: 'pointer' }}>
                    <input type="checkbox"
                      checked={(form.class_ids || []).includes(c.id)}
                      onChange={() => toggleClass(c.id)}
                    />
                    <span style={{ marginLeft: 6 }}>{c.name}（{c.grade}）</span>
                  </label>
                ))}
              </div>
            </div>
            <div className="modal-footer">
              <button className="btn" onClick={() => setShowModal(false)}>取消</button>
              <button className="btn btn-primary" onClick={handleSubmit}>保存</button>
            </div>
          </div>
        </div>
      )}
    </div>
  )
}

export default Teachers
