import React from 'react'
import { Routes, Route, NavLink } from 'react-router-dom'
import Dashboard from './pages/Dashboard.jsx'
import KnowledgePoints from './pages/KnowledgePoints.jsx'
import Semesters from './pages/Semesters.jsx'
import Progress from './pages/Progress.jsx'
import Compare from './pages/Compare.jsx'
import Classes from './pages/Classes.jsx'
import Teachers from './pages/Teachers.jsx'
import Revisions from './pages/Revisions.jsx'

function App() {
  return (
    <div className="layout">
      <aside className="sidebar">
        <h2>📚 课程规划系统</h2>
        <nav>
          <NavLink to="/" end className={({ isActive }) => isActive ? 'active' : ''}>
            📊 仪表盘
          </NavLink>
          <NavLink to="/semesters" className={({ isActive }) => isActive ? 'active' : ''}>
            📅 教学大纲
          </NavLink>
          <NavLink to="/knowledge" className={({ isActive }) => isActive ? 'active' : ''}>
            🧠 知识点
          </NavLink>
          <NavLink to="/progress" className={({ isActive }) => isActive ? 'active' : ''}>
            📈 进度跟踪
          </NavLink>
          <NavLink to="/compare" className={({ isActive }) => isActive ? 'active' : ''}>
            📋 班级对比
          </NavLink>
          <NavLink to="/classes" className={({ isActive }) => isActive ? 'active' : ''}>
            🏫 班级管理
          </NavLink>
          <NavLink to="/teachers" className={({ isActive }) => isActive ? 'active' : ''}>
            👨‍🏫 教师管理
          </NavLink>
          <NavLink to="/revisions" className={({ isActive }) => isActive ? 'active' : ''}>
            📝 修订申请
          </NavLink>
        </nav>
      </aside>
      <main className="main-content">
        <Routes>
          <Route path="/" element={<Dashboard />} />
          <Route path="/semesters" element={<Semesters />} />
          <Route path="/knowledge" element={<KnowledgePoints />} />
          <Route path="/progress" element={<Progress />} />
          <Route path="/compare" element={<Compare />} />
          <Route path="/classes" element={<Classes />} />
          <Route path="/teachers" element={<Teachers />} />
          <Route path="/revisions" element={<Revisions />} />
        </Routes>
      </main>
    </div>
  )
}

export default App
