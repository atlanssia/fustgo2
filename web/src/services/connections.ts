import apiClient from './api'

export interface Connection {
  id: string
  name: string
  type: string
  description: string
  config: string
  enabled: boolean
  created_at: string
  updated_at: string
  last_tested_at?: string
  last_test_status?: string
}

export interface ConnectionCreateRequest {
  name: string
  type: string
  description?: string
  config: Record<string, any>
}

export const connectionService = {
  // 获取连接列表
  list: async (): Promise<Connection[]> => {
    const response = await apiClient.get('/connections')
    return response.data || []
  },

  // 获取连接详情
  get: async (id: string): Promise<Connection> => {
    const response = await apiClient.get(`/connections/${id}`)
    return response.data
  },

  // 创建连接
  create: async (data: ConnectionCreateRequest): Promise<Connection> => {
    const response = await apiClient.post('/connections', data)
    return response.data
  },

  // 更新连接
  update: async (id: string, data: Partial<ConnectionCreateRequest>): Promise<Connection> => {
    const response = await apiClient.put(`/connections/${id}`, data)
    return response.data
  },

  // 删除连接
  delete: async (id: string): Promise<void> => {
    await apiClient.delete(`/connections/${id}`)
  },

  // 测试连接
  test: async (id: string): Promise<{ success: boolean; message: string }> => {
    const response = await apiClient.post(`/connections/${id}/test`)
    return response.data
  },
}
