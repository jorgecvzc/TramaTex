import { describe, it, expect, beforeEach, vi } from 'vitest'
import { render, screen, fireEvent } from '@testing-library/vue'
import { h } from 'vue'
import BaseCatalog from '../../../components/shared/BaseCatalog.vue'
import { createRouter, createWebHistory } from 'vue-router'

// Mocking icons and sub-components
vi.mock('lucide-vue-next', () => ({
  FilterX: { template: '<div>FilterX</div>' },
  RefreshCw: { template: '<div>RefreshCw</div>' },
  Plus: { template: '<div>Plus</div>' },
  AlertCircle: { template: '<div>AlertCircle</div>' },
}))

vi.mock('@/utils/icons', () => ({
  getIcon: vi.fn(() => ({ template: '<div>Icon</div>' }))
}))

// Mock scrollIntoView as JSDOM doesn't implement it
window.HTMLElement.prototype.scrollIntoView = vi.fn()

const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/', component: { template: '<div>Home</div>' } },
    { path: '/create', name: 'create', component: { template: '<div>Create</div>' } },
  ],
})

describe('BaseCatalog Component', () => {
  const defaultProps = {
    title: 'Test Catalog',
    breadcrumbs: [{ label: 'Home', to: '/' }],
    items: [
      { id: '1', name: 'Item 1' },
      { id: '2', name: 'Item 2' },
      { id: '3', name: 'Item 3' },
    ],
  }

  it('should render items correctly', () => {
    render(BaseCatalog, {
      props: defaultProps,
      global: { plugins: [router] },
      slots: {
        item: (props: any) => h('td', { 'data-testid': 'item-name' }, props.item.name),
      },
    })

    expect(screen.getByText('Item 1')).toBeInTheDocument()
    expect(screen.getByText('Item 2')).toBeInTheDocument()
    expect(screen.getByText('Item 3')).toBeInTheDocument()
  })

  it('should navigate with ArrowDown and ArrowUp', async () => {
    const { container } = render(BaseCatalog, {
      props: defaultProps,
      global: { plugins: [router] },
      slots: {
        item: (props: any) => h('td', props.item.name),
      },
    })

    const rows = container.querySelectorAll('tbody tr')
    
    // Initial state: no selection
    expect(container.querySelector('.is-selected')).toBeNull()

    // ArrowDown -> select first item
    await fireEvent.keyDown(window, { key: 'ArrowDown' })
    expect(rows[0]).toHaveClass('is-selected')

    // ArrowDown -> select second item
    await fireEvent.keyDown(window, { key: 'ArrowDown' })
    expect(rows[0]).not.toHaveClass('is-selected')
    expect(rows[1]).toHaveClass('is-selected')

    // ArrowUp -> back to first item
    await fireEvent.keyDown(window, { key: 'ArrowUp' })
    expect(rows[0]).toHaveClass('is-selected')
    expect(rows[1]).not.toHaveClass('is-selected')
  })

  it('should emit click-item when Enter is pressed on a selected item', async () => {
    const { emitted } = render(BaseCatalog, {
      props: defaultProps,
      global: { plugins: [router] },
      slots: {
        item: (props: any) => h('td', props.item.name),
      },
    })

    // Select second item
    await fireEvent.keyDown(window, { key: 'ArrowDown' })
    await fireEvent.keyDown(window, { key: 'ArrowDown' })

    // Press Enter
    await fireEvent.keyDown(window, { key: 'Enter' })

    expect(emitted()['click-item']).toBeTruthy()
    expect(emitted()['click-item'][0]).toEqual([defaultProps.items[1]])
  })

  it('should update selectedIndex when a row is clicked', async () => {
    const { container } = render(BaseCatalog, {
      props: defaultProps,
      global: { plugins: [router] },
      slots: {
        item: (props: any) => h('td', props.item.name),
      },
    })

    const rows = container.querySelectorAll('tbody tr')
    
    await fireEvent.click(rows[2])
    
    expect(rows[2]).toHaveClass('is-selected')
  })

  it('should reset selection with Escape', async () => {
    const { container } = render(BaseCatalog, {
      props: defaultProps,
      global: { plugins: [router] },
      slots: {
        item: (props: any) => h('td', props.item.name),
      },
    })

    const rows = container.querySelectorAll('tbody tr')
    
    await fireEvent.keyDown(window, { key: 'ArrowDown' })
    expect(rows[0]).toHaveClass('is-selected')

    await fireEvent.keyDown(window, { key: 'Escape' })
    expect(container.querySelector('.is-selected')).toBeNull()
  })

  it('should not navigate if items is empty', async () => {
    const { container } = render(BaseCatalog, {
      props: { ...defaultProps, items: [] },
      global: { plugins: [router] },
    })

    await fireEvent.keyDown(window, { key: 'ArrowDown' })
    expect(container.querySelector('.is-selected')).toBeNull()
  })

  it('should not navigate if typing in an input', async () => {
    const { container } = render(BaseCatalog, {
      props: defaultProps,
      global: { plugins: [router] },
      slots: {
        filters: () => h('input', { type: 'text', 'data-testid': 'search-input' }),
        item: (props: any) => h('td', props.item.name),
      },
    })

    const input = screen.getByTestId('search-input')
    await fireEvent.keyDown(input, { key: 'ArrowDown' })
    
    expect(container.querySelector('.is-selected')).toBeNull()
  })
})
