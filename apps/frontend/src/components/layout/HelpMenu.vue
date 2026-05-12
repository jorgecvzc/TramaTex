<template>
  <div class="help-menu-container" ref="menuContainer">
    <!-- Trigger Button -->
    <button 
      class="help-trigger" 
      :class="{ 'is-active': isOpen, 'is-expanded': sidebarExpanded }"
      @click="toggleMenu"
      title="Ayuda y Soporte (Alt+H)"
    >
      <HelpCircle :size="24" />
      <span v-if="sidebarExpanded" class="text">Ayuda y Soporte</span>
    </button>

    <!-- Dropdown Popover -->
    <Transition name="slide-up">
      <div v-if="isOpen" class="help-dropdown" @keydown="handleKeyDown">
        <header class="dropdown-header">
          <HelpCircle :size="16" />
          <span>Centro de Soporte</span>
        </header>
        
        <div class="dropdown-body">
          <button 
            v-for="(item, index) in menuItems" 
            :key="item.id"
            class="dropdown-item"
            :class="{ 'focused': activeIndex === index }"
            @click="executeAction(item)"
            @mouseenter="activeIndex = index"
          >
            <div class="item-icon">
              <component :is="item.icon" :size="18" />
            </div>
            <div class="item-content">
              <span class="item-label">{{ item.label }}</span>
              <span class="item-shortcut">{{ item.shortcutLabel }}</span>
            </div>
          </button>
        </div>
      </div>
    </Transition>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onBeforeUnmount } from 'vue'
import { useRouter } from 'vue-router'
import { 
  HelpCircle, 
  BookOpen, 
  Keyboard, 
  Info,
  LifeBuoy
} from 'lucide-vue-next'

const props = defineProps<{
  sidebarExpanded: boolean
}>()

const emit = defineEmits(['open-shortcuts', 'open-contextual-help'])

const router = useRouter()
const isOpen = ref(false)
const activeIndex = ref(0)
const menuContainer = ref<HTMLElement | null>(null)

const menuItems = [
  { id: 'manual', label: 'Manual de Usuario', icon: BookOpen, shortcutLabel: 'Alt+M', action: () => router.push('/help') },
  { id: 'shortcuts', label: 'Mapa de Atajos', icon: Keyboard, shortcutLabel: '?', action: () => emit('open-shortcuts') },
  { id: 'contextual', label: 'Guía de esta página', icon: Info, shortcutLabel: 'F1', action: () => {
    window.dispatchEvent(new CustomEvent('tramatex-contextual-help'))
  }},
  { id: 'support', label: 'Soporte Técnico', icon: LifeBuoy, shortcutLabel: '', action: () => window.open('mailto:soporte@tramatex.local', '_blank') },
]

function toggleMenu() {
  isOpen.value = !isOpen.value
  if (isOpen.value) activeIndex.value = 0
}

function closeMenu() {
  isOpen.value = false
}

function handleKeyDown(e: KeyboardEvent) {
  if (e.key === 'ArrowDown') {
    e.preventDefault()
    activeIndex.value = (activeIndex.value + 1) % menuItems.length
  } else if (e.key === 'ArrowUp') {
    e.preventDefault()
    activeIndex.value = (activeIndex.value - 1 + menuItems.length) % menuItems.length
  } else if (e.key === 'Enter') {
    e.preventDefault()
    executeAction(menuItems[activeIndex.value])
  } else if (e.key === 'Escape') {
    e.preventDefault()
    closeMenu()
  }
}

function executeAction(item: any) {
  item.action()
  closeMenu()
}

function handleGlobalClick(e: MouseEvent) {
  if (menuContainer.value && !menuContainer.value.contains(e.target as Node)) {
    closeMenu()
  }
}

function handleGlobalKeydown(e: KeyboardEvent) {
  if (e.altKey && e.key.toLowerCase() === 'h') {
    e.preventDefault()
    toggleMenu()
  }
}

onMounted(() => {
  window.addEventListener('click', handleGlobalClick)
  window.addEventListener('keydown', handleGlobalKeydown)
  window.addEventListener('tramatex-contextual-help', () => emit('open-contextual-help'))
})

onBeforeUnmount(() => {
  window.removeEventListener('click', handleGlobalClick)
  window.removeEventListener('keydown', handleGlobalKeydown)
  window.removeEventListener('tramatex-contextual-help', () => emit('open-contextual-help'))
})
</script>

<style scoped>
.help-menu-container {
  position: relative;
  width: 100%;
}

.help-trigger {
  width: 100%;
  display: flex;
  align-items: center;
  gap: 1.25rem;
  padding: 0.75rem 1rem;
  background: transparent;
  border: none;
  color: rgba(255, 255, 255, 0.7);
  cursor: pointer;
  transition: all 0.2s;
  border-radius: 8px;
  margin-bottom: 0.5rem;
}

.help-trigger:hover, .help-trigger.is-active {
  background: rgba(255, 255, 255, 0.1);
  color: white;
}

.help-trigger.is-active {
  color: var(--color-primary);
}

.help-trigger .text {
  font-size: 0.9rem;
  font-weight: 600;
  white-space: nowrap;
}

.help-dropdown {
  position: absolute;
  bottom: calc(100% + 10px);
  left: 0;
  width: 260px;
  background: white;
  border-radius: 12px;
  box-shadow: 0 10px 40px rgba(0,0,0,0.2);
  border: 1px solid var(--color-border);
  overflow: hidden;
  z-index: 3000;
}

.dropdown-header {
  padding: 0.75rem 1rem;
  background: #f8fafc;
  border-bottom: 1px solid var(--color-border);
  display: flex;
  align-items: center;
  gap: 0.5rem;
  font-size: 0.7rem;
  font-weight: 800;
  color: var(--color-text-secondary);
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.dropdown-body {
  padding: 0.5rem;
}

.dropdown-item {
  width: 100%;
  display: flex;
  align-items: center;
  gap: 0.75rem;
  padding: 0.6rem 0.75rem;
  border: none;
  background: transparent;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s;
  text-align: left;
}

.dropdown-item.focused {
  background: #f1f5f9;
}

.dropdown-item.focused .item-icon {
  color: var(--color-secondary);
  transform: scale(1.1);
}

.item-icon {
  color: var(--color-text-secondary);
  transition: transform 0.2s;
}

.item-content {
  flex: 1;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.item-label {
  font-size: 0.9rem;
  font-weight: 500;
  color: var(--color-text-primary);
}

.item-shortcut {
  font-size: 0.7rem;
  font-weight: 700;
  color: var(--color-text-secondary);
  background: #f1f5f9;
  padding: 0.1rem 0.4rem;
  border-radius: 4px;
  min-width: 35px;
  text-align: center;
}

/* Animations */
.slide-up-enter-active, .slide-up-leave-active {
  transition: all 0.25s cubic-bezier(0.4, 0, 0.2, 1);
}
.slide-up-enter-from {
  opacity: 0;
  transform: translateY(10px);
}
.slide-up-leave-to {
  opacity: 0;
  transform: translateY(5px);
}
</style>
