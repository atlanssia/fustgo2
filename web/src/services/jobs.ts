import apiClient from './api'

export interface Job {
  id: string
  name: string
  description: string
  enabled: boolean
  status: string
  pipeline_config: string
  schedule_config: string
  created_at: string
  updated_at: string
  last_execution_time?: string
  next_execution_time?: string
}

export interface JobCreateRequest {
  name: string
  description?: string
  pipeline_config: string
  schedule_config?: string
}

export const jobService = {
  // 获取任务列表
  list: async (): Promise<Job[]> => {
    const response = await apiClient.get('/jobs')
    return response.data || []
  },

  // 获取任务详情
  get: async (id: string): Promise<Job> => {
    const response = await apiClient.get(`/jobs/${id}`)
    return response.data
  },

  // 创建任务
  create: async (data: JobCreateRequest): Promise<Job> => {
    const response = await apiClient.post('/jobs', data)
    return response.data
  },

  // 更新任务
  update: async (id: string, data: Partial<JobCreateRequest>): Promise<Job> => {
    const response = await apiClient.put(`/jobs/${id}`, data)
    return response.data
  },

  // 删除任务
  delete: async (id: string): Promise<void> => {
    await apiClient.delete(`/jobs/${id}`)
  },

  // 执行任务
  execute: async (id: string): Promise<void> => {
    await apiClient.post(`/jobs/${id}/execute`)
  },

  // 暂停任务
  pause: async (id: string): Promise<void> => {
    await apiClient.post(`/jobs/${id}/pause`)
  },

  // 恢复任务
  resume: async (id: string): Promise<void> => {
    await apiClient.post(`/jobs/${id}/resume`)
  },
}
