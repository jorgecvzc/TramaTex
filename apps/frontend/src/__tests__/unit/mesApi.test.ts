import { describe, it, expect, beforeEach, vi, afterEach } from 'vitest'
import { mesApi } from '../../services/mesApi'
import { api } from '../../services/api'

describe('MESApi Service', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    localStorage.clear()
    localStorage.setItem('tramatex_auth_token', 'test-token')
    localStorage.setItem('tramatex_user', JSON.stringify({ id: 'user-123' }))
  })

  afterEach(() => {
    vi.resetAllMocks()
  })

  describe('listTasks', () => {
    it('should list tasks with filters', async () => {
      vi.mocked(api.get).mockResolvedValueOnce({
        data: [
          {
            id: 'task-1',
            name: 'Diseñar',
            is_active: true,
          },
        ],
        status: 200,
      })

      const result = await mesApi.listTasks({ search: 'dis', is_active: true })

      expect(result).toHaveLength(1)
      expect(result[0].name).toBe('Diseñar')
      expect(api.get).toHaveBeenCalledWith(
        expect.stringContaining('/mes/tasks'),
        expect.objectContaining({
          params: expect.objectContaining({
            search: 'dis',
            is_active: true
          })
        })
      )
    })
  })

  describe('work orders', () => {
    it('should create work order', async () => {
      vi.mocked(api.post).mockResolvedValueOnce({
        data: { id: 'wo-1' },
        status: 201,
      })

      const result = await mesApi.createWorkOrder({
        order_id: 'order-1',
        description: 'Test WO',
        items: []
      })

      expect(result.id).toBe('wo-1')
      expect(api.post).toHaveBeenCalled()
    })
  })
})
