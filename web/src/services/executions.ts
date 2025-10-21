import apiClient from './api'

export interface Execution {
  id: string
  job_id: string
  status: string
  start_time: string
  end_time?: string
  duration: number
  records_read: number
  records_written: number
  records_filtered: number
  records_error: number
  bytes_read: number
  bytes_written: number
  error_message?: string
  created_at: string
  updated_at: string
}

export const executionService = {
  // 获取执行历史列表
  list: async (params?: { job_id?: string; limit?: number }): Promise<Execution[]> => {
    const response = await apiClient.get('/executions', { params })
    return response.data || []
  },

  // 获取执行详情
  get: async (id: string): Promise<Execution> => {
    const response = await apiClient.get(`/executions/${id}`)
    return response.data
  },

  // 取消执行
  cancel: async (id: string): Promise<void> => {
    await apiClient.post(`/executions/${id}/cancel`)
  },
}
