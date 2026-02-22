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

  describe('service groups', () => {
    it('should list service groups', async () => {
      ;(globalThis.fetch as any).mockResolvedValueOnce({
        ok: true,
        json: async () => [{ id: 'sg-1', name: 'Bordado', is_active: true, tasks: [] }],
      })

      const result = await mesApi.listServiceGroups()
      expect(result).toHaveLength(1)
      expect(result[0].id).toBe('sg-1')
    })
  })

  describe('createServiceGroup', () => {
    it('should create service group successfully', async () => {
      ;(globalThis.fetch as any).mockResolvedValueOnce({
        ok: true,
        json: async () => ({
          id: 'sg-1',
          name: 'Serigrafía',
          description: '1 color',
          is_active: true,
          tasks: [{ task_id: 'task-1', sequence: 1 }],
        }),
      })

      const result = await mesApi.createServiceGroup({
        name: 'Serigrafía',
        description: '1 color',
        is_active: true,
        task_assignments: [{ task_id: 'task-1', sequence: 1 }],
      })

      expect(result.id).toBe('sg-1')
      expect(globalThis.fetch).toHaveBeenCalledWith(
        expect.stringContaining('/mes/service-groups'),
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
        mesApi.createServiceGroup({
          name: 'Invalid',
          task_assignments: [{ task_id: 'task-1', sequence: 0 }],
        })
      ).rejects.toThrow('No se pudo crear el grupo de servicio MES')
    })
  })

  describe('works', () => {
    it('should list mes works with filters', async () => {
      ;(globalThis.fetch as any).mockResolvedValueOnce({
        ok: true,
        json: async () => [
          {
            id: 'work-1',
            work_number: 'MES-2026-001',
            work_name: 'Trabajo A',
            party_id: 'party-1',
            tangible_group_id: 'group-1',
            status: 'DRAFT',
            priority: 'NORMAL',
            service_groups: [],
          },
        ],
      })

      const result = await mesApi.listWorks({ search: 'Trabajo', status: 'DRAFT' })

      expect(result).toHaveLength(1)
      expect(result[0].work_number).toBe('MES-2026-001')
      const url = (globalThis.fetch as any).mock.calls[0][0] as string
      expect(url).toContain('/mes/works')
      expect(url).toContain('search=Trabajo')
      expect(url).toContain('status=DRAFT')
    })

    it('should create mes work', async () => {
      ;(globalThis.fetch as any).mockResolvedValueOnce({
        ok: true,
        json: async () => ({
          id: 'work-2',
          work_number: 'MES-2026-002',
          work_name: 'Trabajo B',
          party_id: 'party-1',
          tangible_group_id: 'group-1',
          status: 'DRAFT',
          priority: 'NORMAL',
          service_groups: [],
        }),
      })

      const result = await mesApi.createWork({
        work_name: 'Trabajo B',
        party_id: 'party-1',
        tangible_group_id: 'group-1',
        service_group_assignments: [
          {
            service_group_id: 'sg-1',
            position_id: 'pos-1',
            sequence: 1,
          },
        ],
      })

      expect(result.id).toBe('work-2')
      expect(globalThis.fetch).toHaveBeenCalledWith(
        expect.stringContaining('/mes/works'),
        expect.objectContaining({ method: 'POST' })
      )
    })

    it('should get mes work by id', async () => {
      ;(globalThis.fetch as any).mockResolvedValueOnce({
        ok: true,
        json: async () => ({
          id: 'work-1',
          work_number: 'MES-2026-001',
          work_name: 'Trabajo A',
          party_id: 'party-1',
          tangible_group_id: 'group-1',
          status: 'DRAFT',
          priority: 'NORMAL',
          service_groups: [],
        }),
      })

      const result = await mesApi.getWork('work-1')
      expect(result.id).toBe('work-1')
      expect(globalThis.fetch).toHaveBeenCalledWith(
        expect.stringContaining('/mes/works/work-1'),
        expect.objectContaining({ method: 'GET' })
      )
    })

    it('should propagate create work validation error', async () => {
      ;(globalThis.fetch as any).mockResolvedValueOnce({
        ok: false,
        status: 400,
        json: async () => ({ error: 'work name is required' }),
      })

      await expect(
        mesApi.createWork({
          work_name: '',
          party_id: 'party-1',
          tangible_group_id: 'group-1',
          service_group_assignments: [{ service_group_id: 'sg-1', position_id: 'pos-1', sequence: 1 }],
        })
      ).rejects.toThrow('No se pudo crear el trabajo MES')
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

      const result = await mesApi.getWorkDashboardStats()

      expect(result.total).toBe(10)
      expect(result.by_status.IN_PROGRESS).toBe(3)
      expect(globalThis.fetch).toHaveBeenCalledWith(
        expect.stringContaining('/mes/works/dashboard/stats'),
        expect.objectContaining({ method: 'GET' })
      )
    })

    it('should list overdue works with limit', async () => {
      ;(globalThis.fetch as any).mockResolvedValueOnce({
        ok: true,
        json: async () => [
          {
            id: 'work-9',
            work_number: 'MES-2026-009',
            work_name: 'Vencido 1',
            party_id: 'party-1',
            tangible_group_id: 'group-1',
            status: 'IN_PROGRESS',
            priority: 'HIGH',
            due_date: '2026-02-18T00:00:00Z',
            service_groups: [],
          },
        ],
      })

      const result = await mesApi.listOverdueWorks(20)

      expect(result).toHaveLength(1)
      const url = (globalThis.fetch as any).mock.calls[0][0] as string
      expect(url).toContain('/mes/works/overdue')
      expect(url).toContain('limit=20')
    })
  })
})
