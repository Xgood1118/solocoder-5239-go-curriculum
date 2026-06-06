import React, { useState, useEffect } from 'react'
import { api } from '../api.js'

function Dashboard() {
  const [data, setData] = useState(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    loadData()
  }, [])

  const loadData = async () => {
    try {
      const res = await api.getDashboard()
      setData(res)
    } catch (err) {
      console.error(err)
    } finally {
      setLoading(false)
    }
  }

  if (loading) return <div>加载中...</div>

  return (
    <div>
      <div className="page-header">
        <h1>仪表盘</h1>
        <button className="btn btn-primary" onClick={loadData}>刷新</button>
      </div>

      <div className="stat-grid">
        <div className="stat-card">
          <div className="label">学期大纲</div>
          <div className="value primary">{data.semester_count}</div>
        </div>
        <div className="stat-card">
          <div className="label">班级数量</div>
          <div className="value success">{data.class_count}</div>
        </div>
        <div className="stat-card">
          <div className="label">教师数量</div>
          <div className="value">{data.teacher_count}</div>
        </div>
        <div className="stat-card">
          <div className="label">知识点数量</div>
          <div className="value warning">{data.knowledge_point_count}</div>
        </div>
      </div>

      {data.topo_has_cycle && (
        <div className="card" style={{ borderLeft: '4px solid #ff4d4f' }}>
          <div className="card-title" style={{ color: '#ff4d4f' }}>
            ⚠️ 知识点依赖环检测
          </div>
          <p style={{ color: '#ff4d4f', marginBottom: 10 }}>
            检测到知识点依赖存在环，请尽快修改！
          </p>
          <p>
            <strong>环上节点：</strong>
            {data.cycle_info?.cycle?.map((node, i) => (
              <span key={i} className="tag tag-red">{node}</span>
            ))}
          </p>
        </div>
      )}

      <div className="card">
        <div className="card-title">进度异常班级</div>
        {data.critical_deviations?.length === 0 ? (
          <div className="empty-state">暂无异常，所有班级进度正常 🎉</div>
        ) : (
          <table>
            <thead>
              <tr>
                <th>班级</th>
                <th>科目</th>
                <th>偏差（课时）</th>
                <th>级别</th>
                <th>操作</th>
              </tr>
            </thead>
            <tbody>
              {data.critical_deviations?.map((item, i) => (
                <tr key={i}>
                  <td>{item.class_name}</td>
                  <td>{item.subject}</td>
                  <td className={item.deviation < 0 ? 'danger-text' : 'success-text'}>
                    {item.deviation > 0 ? '+' : ''}{item.deviation}
                  </td>
                  <td>
                    <span className={`tag ${item.level === 'critical' ? 'tag-red' : 'tag-orange'}`}>
                      {item.level === 'critical' ? '严重滞后' : '滞后预警'}
                    </span>
                  </td>
                  <td>
                    <button className="btn">查看详情</button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
      </div>

      <div className="card">
        <div className="card-title">快捷操作</div>
        <div style={{ display: 'flex', gap: 10, flexWrap: 'wrap' }}>
          <button className="btn btn-primary">📝 登记进度</button>
          <button className="btn">📊 导出报表</button>
          <button className="btn">🔄 拓扑排序</button>
          <button className="btn">📋 修订申请</button>
        </div>
      </div>
    </div>
  )
}

export default Dashboard
