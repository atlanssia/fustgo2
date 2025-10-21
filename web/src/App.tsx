import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom'
import MainLayout from './layouts/MainLayout'
import Dashboard from './pages/Dashboard'
import Jobs from './pages/Jobs'
import Executions from './pages/Executions'
import Connections from './pages/Connections'
import Plugins from './pages/Plugins'
import Pipelines from './pages/Pipelines'
import PipelineDesigner from './components/PipelineDesigner'

function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<MainLayout />}>
          <Route index element={<Navigate to="/dashboard" replace />} />
          <Route path="dashboard" element={<Dashboard />} />
          <Route path="jobs" element={<Jobs />} />
          <Route path="executions" element={<Executions />} />
          <Route path="connections" element={<Connections />} />
          <Route path="plugins" element={<Plugins />} />
          <Route path="pipelines" element={<Pipelines />} />
          <Route path="config" element={<PipelineDesigner />} />
        </Route>
      </Routes>
    </BrowserRouter>
  )
}

export default App
