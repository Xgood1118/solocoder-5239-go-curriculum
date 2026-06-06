import React, { useState, useEffect } from 'react'
import { api } from '../api.js'

function KnowledgePoints() {
  const [kps, setKps] = useState([])
  const [topoResult, setTopoResult] = useState(null)
  const [showModal, setShowModal] = useState(false)
  const [editing, setEditing] = useState(null)
  const [form, setForm] = useState({
    name: '',
    description: '',
    difficulty: 3,
    prerequisites: [],
    is_key_point: false,
    is_difficult: false,
    is_exam_point: false,
    estimated_lessons: 2,
  })

  useEffect(() => {
    loadData()
    loadTopo()
  }, [])

  const loadData = async () => {
    try {
      const res = await api.listKnowledgePoints()
      setKps(res)
    } catch (err) {
      alert(err.message)
    }
  }

  const loadTopo = async () => {
    try {
      const res = await api.topoSort()
      setTopoResult(res)
    } catch (err) {
      console.error(err)
    }
  }

  const openCreate = () => {
    setEditing(null)
    setForm({
      name: '',
      description: '',
      difficulty: 3,
      prerequisites: [],
      is_key_point: false,
      is_difficult: false,
      is_exam_point: false,
      estimated_lessons: 2,
    })
    setShowModal(true)
  }

  const openEdit = (kp) => {
    setEditing(kp)
    setForm({ ...kp })
    setShowModal(true)
  }

  const handleSubmit = async () => {
    try {
      if (editing) {
        await api.updateKnowledgePoint(editing.id, form)
      } else {
        const id = 'kp_' + Date.now()
        await api.createKnowledgePoint({ ...form, id })
      }
      setShowModal(false)
      loadData()
      loadTopo()
    } catch (err) {
      alert(err.message)
    }
  }

  const handleDelete = async (id) => {
    if (!confirm('确定删除此知识点？')) return
    try {
      await api.deleteKnowledgePoint(id)
      loadData()
      loadTopo()
    } catch (err) {
      alert(err.message)
    }
  }

  const togglePrereq = (kpId) => {
    const pres = form.prerequisites || []
    if (pres.includes(kpId)) {
      setForm({ ...form, prerequisites: pres.filter(p => p !== kpId) })
    } else {
      setForm({ ...form, prerequisites: [...pres, kpId] })
    }
  }

  const renderStars = (level) => {
    return '⭐'.repeat(level) + '☆'.repeat(5 - level)
  }

  return (
    <div>
      <div className="page-header">
        <h1>知识点管理</h1>
        <div style={{ display: 'flex', gap: 8 }}>
          <button className="btn" onClick={loadTopo}>🔄 拓扑排序</button>
          <button className="btn btn-primary" onClick={openCreate}>+ 新增知识点</button>
        </div>
      </div>

      {topoResult?.cycle?.has_cycle && (
        <div className="card" style={{ borderLeft: '4px solid #ff4d4f' }}>
          <strong style={{ color: '#ff4d4f' }}>⚠️ 检测到依赖环：</strong>
          {topoResult.cycle.cycle.map((node, i) => (
            <span key={i} className="tag tag-red" style={{ marginLeft: 8 }}>{node}</span>
          ))}
        </div>
      )}

      {topoResult?.sorted && !topoResult.cycle?.has_cycle && (
        <div className="card">
          <div className="card-title">拓扑排序结果</div>
          <div style={{ display: 'flex', flexWrap: 'wrap', gap: 8 }}>
            {topoResult.sorted.map((id, i) => {
              const kp = kps.find(k => k.id === id)
              return (
                <span key={id} className="tag tag-blue">
                  第{topoResult.levels[id]}层 - {kp?.name || id}
                </span>
              )
            })}
          </div>
        </div>
      )}

      <div className="card">
        <div className="card-title">知识点列表</div>
        <table>
          <thead>
            <tr>
              <th>名称</th>
              <th>难度</th>
              <th>预估课时</th>
              <th>标签</th>
              <th>前置依赖</th>
              <th>操作</th>
            </tr>
          </thead>
          <tbody>
            {kps.map(kp => (
              <tr key={kp.id}>
                <td><strong>{kp.name}</strong></td>
                <td>{renderStars(kp.difficulty)}</td>
                <td>{kp.estimated_lessons} 课时</td>
                <td>
                  {kp.is_key_point && <span className="tag tag-red">重点</span>}
                  {kp.is_difficult && <span className="tag tag-orange">难点</span>}
                  {kp.is_exam_point && <span className="tag tag-purple">考点</span>}
                </td>
                <td style={{ fontSize: 12, color: '#999' }}>
                  {kp.prerequisites?.length > 0
                    ? kp.prerequisites.map(pid => kps.find(k => k.id === pid)?.name || pid).join(', ')
                    : '无'}
                </td>
                <td>
                  <button className="btn" style={{ marginRight: 6 }} onClick={() => openEdit(kp)}>编辑</button>
                  <button className="btn btn-danger" onClick={() => handleDelete(kp.id)}>删除</button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
        {kps.length === 0 && <div className="empty-state">暂无知识点</div>}
      </div>

      {showModal && (
        <div className="modal-overlay" onClick={() => setShowModal(false)}>
          <div className="modal" onClick={e => e.stopPropagation()}>
            <div className="modal-header">
              <h3>{editing ? '编辑知识点' : '新增知识点'}</h3>
              <button className="close-btn" onClick={() => setShowModal(false)}>×</button>
            </div>

            <div className="form-group">
              <label>名称</label>
              <input type="text" value={form.name} onChange={e => setForm({ ...form, name: e.target.value })} />
            </div>

            <div className="form-group">
              <label>描述</label>
              <textarea rows={3} value={form.description || ''} onChange={e => setForm({ ...form, description: e.target.value })} />
            </div>

            <div className="form-group">
              <label>难度系数（1-5）</label>
              <select value={form.difficulty} onChange={e => setForm({ ...form, difficulty: parseInt(e.target.value) })}>
                {[1, 2, 3, 4, 5].map(n => (
                  <option key={n} value={n}>{n} 星 - {['入门', '简单', '中等', '较难', '困难'][n - 1]}</option>
                ))}
              </select>
            </div>

            <div className="form-group">
              <label>预估课时</label>
              <input type="number" min={1} value={form.estimated_lessons} onChange={e => setForm({ ...form, estimated_lessons: parseInt(e.target.value) || 1 })} />
            </div>

            <div className="form-group">
              <label>前置依赖知识点</label>
              <div style={{ display: 'flex', flexWrap: 'wrap', gap: 8, maxHeight: 150, overflowY: 'auto', border: '1px solid #e8e8e8', padding: 10, borderRadius: 4 }}>
                {kps.filter(k => k.id !== editing?.id).map(kp => (
                  <label key={kp.id} style={{ display: 'flex', alignItems: 'center', gap: 4, cursor: 'pointer' }}>
                    <input
                      type="checkbox"
                      checked={(form.prerequisites || []).includes(kp.id)}
                      onChange={() => togglePrereq(kp.id)}
                    />
                    <span style={{ fontSize: 13 }}>{kp.name}</span>
                  </label>
                ))}
              </div>
            </div>

            <div style={{ display: 'flex', gap: 20 }}>
              <label style={{ display: 'flex', alignItems: 'center', gap: 6, cursor: 'pointer' }}>
                <input type="checkbox" checked={form.is_key_point} onChange={e => setForm({ ...form, is_key_point: e.target.checked })} />
                重点
              </label>
              <label style={{ display: 'flex', alignItems: 'center', gap: 6, cursor: 'pointer' }}>
                <input type="checkbox" checked={form.is_difficult} onChange={e => setForm({ ...form, is_difficult: e.target.checked })} />
                难点
              </label>
              <label style={{ display: 'flex', alignItems: 'center', gap: 6, cursor: 'pointer' }}>
                <input type="checkbox" checked={form.is_exam_point} onChange={e => setForm({ ...form, is_exam_point: e.target.checked })} />
                考点
              </label>
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

export default KnowledgePoints
