import React, { useState, useEffect } from 'react'
import { api } from '../api.js'

function Classes() {
  const [classes, setClasses] = useState([])
  const [showModal, setShowModal] = useState(false)
  const [editing, setEditing] = useState(null)
  const [form, setForm] = useState({ name: '', grade: '高一年级', student_count: 45 })

  useEffect(() => { loadData() }, [])

  const loadData = async () => {
    try {
      const res = await api.listClasses()
      setClasses(res)
    } catch (err) { alert(err.message) }
  }

  const openCreate = () => {
    setEditing(null)
    setForm({ name: '', grade: '高一年级', student_count: 45 })
    setShowModal(true)
  }

  const openEdit = (c) => {
    setEditing(c)
    setForm({ ...c })
    setShowModal(true)
  }

  const handleSubmit = async () => {
    try {
      if (editing) {
        await api.updateClass(editing.id, form)
      } else {
        await api.createClass({ ...form, id: 'cls_' + Date.now() })
      }
      setShowModal(false)
      loadData()
    } catch (err) { alert(err.message) }
  }

  const handleDelete = async (id) => {
    if (!confirm('确定删除？')) return
    try {
      await api.deleteClass(id)
      loadData()
    } catch (err) { alert(err.message) }
  }

  return (
    <div>
      <div className="page-header">
        <h1>班级管理</h1>
        <button className="btn btn-primary" onClick={openCreate}>+ 新增班级</button>
      </div>

      <div className="card">
        <table>
          <thead>
            <tr>
              <th>班级名称</th>
              <th>年级</th>
              <th>学生数</th>
              <th>教师数</th>
              <th>操作</th>
            </tr>
          </thead>
          <tbody>
            {classes.map(c => (
              <tr key={c.id}>
                <td><strong>{c.name}</strong></td>
                <td>{c.grade}</td>
                <td>{c.student_count}人</td>
                <td>{c.teacher_ids?.length || 0}位</td>
                <td>
                  <button className="btn" style={{ marginRight: 6 }} onClick={() => openEdit(c)}>编辑</button>
                  <button className="btn btn-danger" onClick={() => handleDelete(c.id)}>删除</button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
        {classes.length === 0 && <div className="empty-state">暂无班级</div>}
      </div>

      {showModal && (
        <div className="modal-overlay" onClick={() => setShowModal(false)}>
          <div className="modal" onClick={e => e.stopPropagation()}>
            <div className="modal-header">
              <h3>{editing ? '编辑班级' : '新增班级'}</h3>
              <button className="close-btn" onClick={() => setShowModal(false)}>×</button>
            </div>
            <div className="form-group">
              <label>班级名称</label>
              <input value={form.name} onChange={e => setForm({ ...form, name: e.target.value })} placeholder="如：高一(1)班" />
            </div>
            <div className="form-group">
              <label>年级</label>
              <select value={form.grade} onChange={e => setForm({ ...form, grade: e.target.value })}>
                <option>高一年级</option><option>高二年级</option><option>高三年级</option>
                <option>初一年级</option><option>初二年级</option><option>初三年级</option>
              </select>
            </div>
            <div className="form-group">
              <label>学生人数</label>
              <input type="number" value={form.student_count} onChange={e => setForm({ ...form, student_count: parseInt(e.target.value) })} />
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

export default Classes
