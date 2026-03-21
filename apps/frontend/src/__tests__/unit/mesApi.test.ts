import { describe, it, expect, beforeEach, vi, afterEach } from 'vitest'
import { mesApi } from '../../services/mesApi'

globalThis.fetch = vi.fn()

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
      ;(globalThis.fetch as any).mockResolvedValueOnce({
        ok: true,
        json: async () => [
          {
            id: 'task-1',
            name: 'Diseñar',
            description: 'Diseño',
            is_active: true,
          },
        ],
      })

      const result = await mesApi.listTasks({ search: 'dis', is_active: true })

      expect(result).toHaveLength(1)
      expect(result[0].name).toBe('Diseñar')
      const url = (globalThis.fetch as any).mock.calls[0][0] as string
      expect(url).toContain('/mes/tasks')
      expect(url).toContain('search=dis')
      expect(url).toContain('is_active=true')
      expect(globalThis.fetch).toHaveBeenCalledWith(
        expect.any(String),
        expect.objectContaining({
          method: 'GET',
          headers: expect.objectContaining({
            Authorization: 'Bearer test-token',
            'X-User-ID': 'user-123',
          }),
        })
      )
    })

    it('should propagate listTasks backend error', async () => {
      ;(globalThis.fetch as any).mockResolvedValueOnce({
        ok: false,
        status: 500,
        json: async () => ({ message: 'error interno' }),
      })

      await expect(mesApi.listTasks()).rejects.toThrow('No se pudieron cargar las tareas MES')
    })
  })

  describe('tasks and positions', () => {
    it('should create task successfully', async () => {
      ;(globalThis.fetch as any).mockResolvedValueOnce({
        ok: true,
        json: async () => ({
          id: 'task-2',
          name: 'Imprimir',
          is_active: true,
        }),
      })

      const result = await mesApi.createTask({ name: 'Imprimir' })
      expect(result.id).toBe('task-2')
    })

    it('should list positions', async () => {
      ;(globalThis.fetch as any).mockResolvedValueOnce({
        ok: true,
        json: async () => [
          { id: 'pos-1', name: 'Espalda', code: 'BACK', is_active: true },
        ],
      })

      const result = await mesApi.listPositions({ search: 'back' })
      expect(result).toHaveLength(1)
      const url = (globalThis.fetch as any).mock.calls[0][0] as string
      expect(url).toContain('/mes/positions')
      expect(url).toContain('search=back')
    })
  })

  describe('work types', () => {
    it('should list work types', async () => {
      ;(globalThis.fetch as any).mockResolvedValueOnce({
        ok: true,
        json: async () => [{ id: 'wt-1', name: 'Bordado', is_active: true, tasks: [] }],
      })

      const result = await mesApi.listWorkTypes()
      expect(result).toHaveLength(1)
      expect(result[0].id).toBe('wt-1')
      const url = (globalThis.fetch as any).mock.calls[0][0] as string
      expect(url).toContain('/mes/work-types')
    })
  })

  describe('createWorkType', () => {
    it('should create work type successfully', async () => {
      ;(globalThis.fetch as any).mockResolvedValueOnce({
        ok: true,
        json: async () => ({
          id: 'wt-1',
          name: 'Serigrafía',
          description: '1 color',
          is_active: true,
          tasks: [{ task_id: 'task-1', sequence: 1 }],
        }),
      })

      const result = await mesApi.createWorkType({
        name: 'Serigrafía',
        description: '1 color',
        is_active: true,
        task_assignments: [{ task_id: 'task-1', sequence: 1 }],
      })

      expect(result.id).toBe('wt-1')
      expect(globalThis.fetch).toHaveBeenCalledWith(
        expect.stringContaining('/mes/work-types'),
        expect.objectContaining({ method: 'POST' })
      )
    })

    it('should surface backend error message', async () => {
      ;(globalThis.fetch as any).mockResolvedValueOnce({
        ok: false,
        status: 400,
        json: async () => ({ error: 'task sequence must be greater than zero' }),
      })

      await expect(
        mesApi.createWorkType({
          name: 'Invalid',
          task_assignments: [{ task_id: 'task-1', sequence: 0 }],
        })
      ).rejects.toThrow('No se pudo crear el tipo de trabajo MES')
    })
  })

  describe('work orders', () => {
    it('should list work orders with filters', async () => {
      ;(globalThis.fetch as any).mockResolvedValueOnce({
        ok: true,
        json: async () => [
          {
            id: 'work-1',
            work_number: 'MES-2026-001',
            work_name: 'Trabajo A',
            party_id: 'party-1',
            work_setup_id: 'ws-1',
            status: 'PENDING',
            priority: 'NORMAL',
            lines: [],
          },
        ],
      })

      const result = await mesApi.listWorkOrders({ search: 'Trabajo', status: 'DRAFT' })

      expect(result).toHaveLength(1)
      expect(result[0].work_number).toBe('MES-2026-001')
      const url = (globalThis.fetch as any).mock.calls[0][0] as string
      expect(url).toContain('/mes/work-orders')
      expect(url).toContain('search=Trabajo')
      expect(url).toContain('status=DRAFT')
    })

    it('should create work order', async () => {
      ;(globalThis.fetch as any).mockResolvedValueOnce({
        ok: true,
        json: async () => ({
          id: 'work-2',
          work_number: 'MES-2026-002',
          work_name: 'Trabajo B',
          party_id: 'party-1',
          work_setup_id: 'ws-1',
          status: 'PENDING',
          priority: 'NORMAL',
          lines: [],
        }),
      })

      const result = await mesApi.createWorkOrder({
        work_name: 'Trabajo B',
        party_id: 'party-1',
        work_setup_id: 'ws-1',
      })

      expect(result.id).toBe('work-2')
      expect(globalThis.fetch).toHaveBeenCalledWith(
        expect.stringContaining('/mes/work-orders'),
        expect.objectContaining({ method: 'POST' })
      )
    })

    it('should get work order by id', async () => {
      ;(globalThis.fetch as any).mockResolvedValueOnce({
        ok: true,
        json: async () => ({
          id: 'work-1',
          work_number: 'MES-2026-001',
          work_name: 'Trabajo A',
          party_id: 'party-1',
          work_setup_id: 'ws-1',
          status: 'PENDING',
          priority: 'NORMAL',
          lines: [],
        }),
      })

      const result = await mesApi.getWorkOrder('work-1')
      expect(result.id).toBe('work-1')
      expect(globalThis.fetch).toHaveBeenCalledWith(
        expect.stringContaining('/mes/work-orders/work-1'),
        expect.objectContaining({ method: 'GET' })
      )
    })

    it('should propagate create work order validation error', async () => {
      ;(globalThis.fetch as any).mockResolvedValueOnce({
        ok: false,
        status: 400,
        json: async () => ({ error: 'work name is required' }),
      })

      await expect(
        mesApi.createWorkOrder({
          work_name: '',
          party_id: 'party-1',
          work_setup_id: 'ws-1',
        })
      ).rejects.toThrow('No se pudo crear la orden de trabajo MES')
    })

    it('should get dashboard stats', async () => {
      ;(globalThis.fetch as any).mockResolvedValueOnce({
        ok: true,
        json: async () => ({
          total: 10,
          by_status: { DRAFT: 4, IN_PROGRESS: 3, COMPLETED: 3 },
          overdue: 2,
          due_today: 1,
        }),
      })

      const result = await mesApi.getWorkOrderDashboardStats()

      expect(result.total).toBe(10)
      expect(result.by_status.IN_PROGRESS).toBe(3)
      expect(globalThis.fetch).toHaveBeenCalledWith(
        expect.stringContaining('/mes/work-orders/dashboard/stats'),
        expect.objectContaining({ method: 'GET' })
      )
    })

    it('should list overdue work orders with limit', async () => {
      ;(globalThis.fetch as any).mockResolvedValueOnce({
        ok: true,
        json: async () => [
          {
            id: 'work-9',
            work_number: 'MES-2026-009',
            work_name: 'Vencido 1',
            party_id: 'party-1',
            work_setup_id: 'ws-1',
            status: 'IN_PROGRESS',
            priority: 'HIGH',
            due_date: '2026-02-18T00:00:00Z',
            lines: [],
          },
        ],
      })

      const result = await mesApi.listOverdueWorkOrders(20)

      expect(result).toHaveLength(1)
      const url = (globalThis.fetch as any).mock.calls[0][0] as string
      expect(url).toContain('/mes/work-orders/overdue')
      expect(url).toContain('limit=20')
    })
  })
})
