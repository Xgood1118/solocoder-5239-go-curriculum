import React, { useState, useEffect } from 'react'
import { api } from '../api.js'

function Revisions() {
  const [requests, setRequests] = useState([])
  const [semesters, setSemesters] = useState([])
  const [showModal, setShowModal] = useState(false)
  const [form, setForm] = useState({ original_plan_id: '', title: '', reason: '', applicant_id: '' })

  useEffect(() => {
    loadData()
    loadSemesters()
  }, [])

  const loadData = async () => {
    try {
      const res = await api.listRevisionRequests()
      setRequests(res)
    } catch (err) { alert(err.message) }
  }

  const loadSemesters = async () => {
    try {
      const res = await api.listSemesters()
      setSemesters(res)
    } catch (err) { console.error(err) }
  }

  const handleSubmit = async () => {
    try {
      await api.createRevisionRequest({
        ...form,
        id: 'rev_' + Date.now(),
        originalPlanId: form.original_plan_id,
        applicantId: form.applicant_id,
      })
      setShowModal(false)
      loadData()
    } catch (err) { alert(err.message) }
  }

  const handleApprove = async (id) => {
    const newPlanId = prompt('请输入新计划ID：')
    if (!newPlanId) return
    try {
      await api.approveRevision(id, {
        approver_id: 'admin',
        approval_note: '同意修订',
        new_plan_id: newPlanId,
      })
      loadData()
    } catch (err) { alert(err.message) }
  }

  const handleReject = async (id) => {
    const note = prompt('请输入驳回理由：')
    if (!note) return
    try {
      await api.rejectRevision(id, {
        approver_id: 'admin',
        approval_note: note,
      })
      loadData()
    } catch (err) { alert(err.message) }
  }

  const statusTag = (status) => {
    const map = {
      pending: { text: '待审批', class: 'tag tag-orange' },
      approved: { text: '已通过', class: 'tag tag-green' },
      rejected: { text: '已驳回', class: 'tag tag-red' },
    }
    const info = map[status] || { text: status, class: 'tag' }
    return <span className={info.class}>{info.text}</span>
  }

  return (
    <div>
      <div className="page-header">
        <h1>修订申请</h1>
        <button className="btn btn-primary" onClick={() => setShowModal(true)}>+ 提交申请</button>
      </div>

      <div className="card">
        <table>
          <thead>
            <tr>
              <th>申请标题</th>
              <th>原计划</th>
              <th>状态</th>
              <th>申请人</th>
              <th>申请时间</th>
              <th>操作</th>
            </tr>
          </thead>
          <tbody>
            {requests.map(r => {
              const sem = semesters.find(s => s.id === r.original_plan_id || s.id === r.originalPlanId)
              return (
                <tr key={r.id}>
                  <td><strong>{r.title}</strong></td>
                  <td>{sem?.name || r.original_plan_id || r.originalPlanId}</td>
                  <td>{statusTag(r.status)}</td>
                  <td>{r.applicant_id || r.applicantId || '-'}</td>
                  <td style={{ fontSize: 12, color: '#999' }}>
                    {new Date(r.created_at || r.createdAt).toLocaleDateString()}
                  </td>
                  <td>
                    {r.status === 'pending' && (
                      <>
                        <button className="btn" style={{ marginRight: 6 }} onClick={() => handleApprove(r.id)}>通过</button>
                        <button className="btn btn-danger" onClick={() => handleReject(r.id)}>驳回</button>
                      </>
                    )}
                  </td>
                </tr>
              )
            })}
          </tbody>
        </table>
        {requests.length === 0 && <div className="empty-state">暂无修订申请</div>}
      </div>

      {showModal && (
        <div className="modal-overlay" onClick={() => setShowModal(false)}>
          <div className="modal" onClick={e => e.stopPropagation()}>
            <div className="modal-header">
              <h3>提交修订申请</h3>
              <button className="close-btn" onClick={() => setShowModal(false)}>×</button>
            </div>
            <div className="form-group">
              <label>原计划</label>
              <select value={form.original_plan_id} onChange={e => setForm({ ...form, original_plan_id: e.target.value })}>
                <option value="">请选择</option>
                {semesters.map(s => (
                  <option key={s.id} value={s.id}>{s.name}（{s.subject}）</option>
                ))}
              </select>
            </div>
            <div className="form-group">
              <label>申请标题</label>
              <input value={form.title} onChange={e => setForm({ ...form, title: e.target.value })} placeholder="简要描述修订内容" />
            </div>
            <div className="form-group">
              <label>修订原因</label>
              <textarea rows={3} value={form.reason} onChange={e => setForm({ ...form, reason: e.target.value })} />
            </div>
            <div className="form-group">
              <label>申请人</label>
              <input value={form.applicant_id} onChange={e => setForm({ ...form, applicant_id: e.target.value })} placeholder="教师姓名或ID" />
            </div>
            <div className="modal-footer">
              <button className="btn" onClick={() => setShowModal(false)}>取消</button>
              <button className="btn btn-primary" onClick={handleSubmit}>提交申请</button>
            </div>
          </div>
        </div>
      )}
    </div>
  )
}

export default Revisions
