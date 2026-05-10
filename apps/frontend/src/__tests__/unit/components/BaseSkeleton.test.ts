import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import BaseSkeleton from '@/components/shared/BaseSkeleton.vue'

describe('BaseSkeleton.vue', () => {
  it('renders with default props', () => {
    const wrapper = mount(BaseSkeleton)
    expect(wrapper.classes()).toContain('skeleton')
    expect(wrapper.classes()).toContain('skeleton-text')
  })

  it('renders correct type class', () => {
    const wrapper = mount(BaseSkeleton, {
      props: { type: 'circle' }
    })
    expect(wrapper.classes()).toContain('skeleton-circle')
  })

  it('applies custom styles', () => {
    const wrapper = mount(BaseSkeleton, {
      props: {
        width: '100px',
        height: '50px',
        borderRadius: '10px'
      }
    })
    const style = wrapper.attributes('style')
    expect(style).toContain('width: 100px;')
    expect(style).toContain('height: 50px;')
    expect(style).toContain('border-radius: 10px;')
  })

  it('disables animation when requested', () => {
    const wrapper = mount(BaseSkeleton, {
      props: { animated: false }
    })
    expect(wrapper.classes()).toContain('no-animation')
  })
})
