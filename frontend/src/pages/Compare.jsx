import React, { useState, useEffect } from 'react'
import { api } from '../api.js'

function Compare() {
  const [classes, setClasses] = useState([])
  const [classPlans, setClassPlans] = useState([])
  const [selectedPlanIds, setSelectedPlanIds] = useState([])
  const [ganttData, setGanttData] = useState([])
  const [semester, setSemester] = useState(null)

  useEffect(() => {
    loadClasses()
  }, [])

  const loadClasses = async () => {
    try {
      const res = await api.listClasses()
      setClasses(res)
    } catch (err) { console.error(err) }
  }

  const togglePlan = async (planId) => {
    let newSelected
    if (selectedPlanIds.includes(planId)) {
      newSelected = selectedPlanIds.filter(id => id !== planId)
    } else {
      if (selectedPlanIds.length >= 2) {
        alert('最多选择2个班级进行对比')
        return
      }
      newSelected = [...selectedPlanIds, planId]
    }
    setSelectedPlanIds(newSelected)

    if (newSelected.length >= 2) {
      try {
        const data = await api.compareGantt(newSelected)
        setGanttData(data)
      } catch (err) {
        console.error(err)
      }
    } else {
      setGanttData([])
    }
  }

  const onSelectClass = async (classId) => {
    try {
      const res = await api.listClassPlans(classId)
      setClassPlans(res)
    } catch (err) { console.error(err) }
  }

  const totalWeeks = 18

  const colorMap = ['#1890ff', '#52c41a', '#faad14', '#722ed1']

  return (
    <div>
      <div className="page-header">
        <h1>班级进度对比</h1>
      </div>

      <div className="card" style={{ marginBottom: 20 }}>
        <div className="card-title">选择班级计划（最多2个）</div>
        <div style={{ display: 'flex', gap: 20, flexWrap: 'wrap' }}>
          {classes.map(c => (
            <div key={c.id} style={{ minWidth: 200 }}>
              <div style={{ fontWeight: 500, marginBottom: 8 }} onClick={() => onSelectClass(c.id)} style={{ cursor: 'pointer' }}>
                📚 {c.name}
              </div>
              {classPlans.filter(p => p.class_id === c.id || p.classId === c.id).map(p => {
                const pid = p.id
                const selected = selectedPlanIds.includes(pid)
                return (
                  <div key={pid}
                    style={{
                      padding: '6px 10px',
                      marginBottom: 4,
                      background: selected ? '#e6f7ff' : '#fafafa',
                      border: `1px solid ${selected ? '#1890ff' : '#e8e8e8'}`,
                      borderRadius: 4,
                      cursor: 'pointer',
                      fontSize: 13,
                    }}
                    onClick={() => togglePlan(pid)}
                  >
                    {selected ? '✅ ' : '⬜ '}{p.subject}
                  </div>
                )
              })}
            </div>
          ))}
        </div>
      </div>

      {ganttData.length >= 2 && (
        <div className="card">
          <div className="card-title">甘特图对比</div>
          <div className="gantt-container">
            <div className="gantt-chart">
              <div className="gantt-header">
                <div className="label">单元</div>
                <div className="weeks">
                  {Array.from({ length: totalWeeks }, (_, i) => (
                    <div key={i} className="week">第{i + 1}周</div>
                  ))}
                </div>
              </div>

              {ganttData[0]?.bars?.map((bar, idx) => (
                <div key={bar.unit_id} className="gantt-swimlanes">
                  {ganttData.map((data, dataIdx) => (
                    <div key={data.class_id} className="swimlane">
                      <div className="gantt-row">
                        <div className="label" style={{ fontSize: 12, padding: '6px 10px' }}>
                          <span style={{
                            display: 'inline-block',
                            width: 10, height: 10, borderRadius: 2,
                            background: colorMap[dataIdx % colorMap.length],
                            marginRight: 6, verticalAlign: 'middle'
                          }} />
                          {data.class_name}
                        </div>
                        <div className="bars" style={{ position: 'relative', minHeight: 24 }}>
                          {(() => {
                            const b = data.bars[idx]
                            if (!b) return null
                            const left = ((b.start_week - 1) / totalWeeks) * 100
                            const width = ((b.end_week - b.start_week + 1) / totalWeeks) * 100
                            return (
                              <div
                                className="gantt-bar"
                                style={{
                                  left: `${left}%`,
                                  width: `${width}%`,
                                  background: colorMap[dataIdx % colorMap.length],
                                  opacity: 0.7 + dataIdx * 0.1,
                                }}
                                title={`${b.unit_title} (第${b.start_week}-${b.end_week}周)`}
                              >
                                {b.unit_title.length > 6 ? b.unit_title.slice(0, 6) + '...' : b.unit_title}
                              </div>
                            )
                          })()}
                        </div>
                      </div>
                    </div>
                  ))}
                </div>
              ))}
            </div>
          </div>

          <div style={{ marginTop: 16, display: 'flex', gap: 20 }}>
            {ganttData.map((data, idx) => (
              <div key={data.class_id} style={{ display: 'flex', alignItems: 'center', gap: 6 }}>
                <span style={{
                  display: 'inline-block', width: 16, height: 12,
                  background: colorMap[idx % colorMap.length], borderRadius: 2
                }} />
                <span style={{ fontSize: 13 }}>{data.class_name}</span>
              </div>
            ))}
          </div>
        </div>
      )}

      {selectedPlanIds.length < 2 && (
        <div className="empty-state">请选择至少2个班级计划进行对比</div>
      )}
    </div>
  )
}

export default Compare
