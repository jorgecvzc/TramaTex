import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import Navbar from '@/components/layout/Navbar.vue'
import { createPinia, setActivePinia } from 'pinia'
import { createRouter, createWebHistory } from 'vue-router'

// Mock de componentes hijos para evitar errores de carga
vi.mock('@/components/layout/UserMenu.vue', () => ({
  default: { template: '<div>UserMenu</div>' }
}))
vi.mock('@/components/shared/GlobalSearch.vue', () => ({
  default: { template: '<div>GlobalSearch</div>' }
}))

describe('Global Keyboard Shortcuts', () => {
  let router: any

  beforeEach(() => {
    setActivePinia(createPinia())
    router = createRouter({
      history: createWebHistory(),
      routes: [
        { path: '/', component: { template: '<div>Home</div>' } },
        { path: '/parties', component: { template: '<div>Parties</div>' } },
        { path: '/parties/new', component: { template: '<div>New Party</div>' } }
      ]
    })
  })

  it('Ctrl+K opens global search', async () => {
    const wrapper = mount(Navbar, {
      global: { plugins: [router] }
    })
    
    const event = new KeyboardEvent('keydown', { ctrlKey: true, key: 'k' })
    document.dispatchEvent(event)
    
    // El componente GlobalSearch debería recibir prop show: true
    // En el Navbar, esto se gestiona con la ref showSearch
    expect((wrapper.vm as any).showSearch).toBe(true)
  })

  it('Alt+N redirects to context-aware new page', async () => {
    await router.push('/parties')
    mount(Navbar, {
      global: { plugins: [router] }
    })
    
    const pushSpy = vi.spyOn(router, 'push')
    const event = new KeyboardEvent('keydown', { altKey: true, key: 'n' })
    document.dispatchEvent(event)
    
    expect(pushSpy).toHaveBeenCalledWith('/parties/new')
  })

  it('Alt+R dispatches custom refresh event', async () => {
    mount(Navbar, {
      global: { plugins: [router] }
    })
    
    const dispatchSpy = vi.spyOn(window, 'dispatchEvent')
    const event = new KeyboardEvent('keydown', { altKey: true, key: 'r' })
    document.dispatchEvent(event)
    
    expect(dispatchSpy).toHaveBeenCalledWith(expect.objectContaining({ type: 'tramatex-refresh' }))
  })
})
