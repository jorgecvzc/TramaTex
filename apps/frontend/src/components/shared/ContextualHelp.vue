<template>
  <Transition name="slide-right">
    <aside v-if="show" class="contextual-help-panel" @keydown.esc="$emit('close')">
      <header class="panel-header">
        <div class="flex items-center gap-3">
          <Info :size="24" class="text-secondary" />
          <div class="title-group">
            <h3>Guía Rápida</h3>
            <span class="route-label">{{ routeName }}</span>
          </div>
        </div>
        <button class="btn-close" @click="$emit('close')"><X :size="20" /></button>
      </header>

      <div class="panel-body">
        <!-- Content reacciona a la ruta actual -->
        <section class="help-section">
          <h4>¿Qué puedo hacer aquí?</h4>
          <p>{{ currentHelp.description }}</p>
        </section>

        <section class="help-section" v-if="currentHelp.actions?.length">
          <h4>Acciones Principales</h4>
          <ul class="actions-list">
            <li v-for="action in currentHelp.actions" :key="action.label">
              <strong>{{ action.label }}:</strong> {{ action.text }}
            </li>
          </ul>
        </section>

        <section class="help-section" v-if="currentHelp.tips?.length">
          <h4>Pro-Tips de Teclado</h4>
          <div v-for="tip in currentHelp.tips" :key="tip.key" class="tip-card">
            <kbd>{{ tip.key }}</kbd>
            <span>{{ tip.text }}</span>
          </div>
        </section>

        <div class="panel-footer">
          <button class="btn btn-primary w-full justify-center" @click="goToFullHelp">
            <BookOpen :size="18" /> Ver Manual Completo
          </button>
        </div>
      </div>
    </aside>
  </Transition>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { Info, X, BookOpen } from 'lucide-vue-next'

const props = defineProps<{ show: boolean }>()
const emit = defineEmits(['close'])

const route = useRoute()
const router = useRouter()

const routeName = computed(() => {
  const path = route.path
  if (path.includes('/sales/quotes')) return 'Presupuestos'
  if (path.includes('/sales/orders')) return 'Pedidos'
  if (path.includes('/parties')) return 'Entidades'
  if (path.includes('/mes')) return 'Taller'
  return 'General'
})

const helpContent: any = {
  'Presupuestos': {
    description: 'Gestión de ofertas comerciales. Los presupuestos no afectan al stock ni a la producción hasta que son confirmados.',
    actions: [
      { label: 'Crear', text: 'Define cliente, líneas y validez.' },
      { label: 'Confirmar', text: 'Convierte la oferta en un Pedido en firme.' }
    ],
    tips: [
      { key: 'Insert', text: 'Añadir línea de producto' },
      { key: 'Ctrl+S', text: 'Guardar borrador' }
    ]
  },
  'Pedidos': {
    description: 'Compromiso de venta. Desde aquí se lanzan los trabajos al taller (MES) y se preparan las entregas.',
    actions: [
      { label: 'Lanzar MES', text: 'Envía las instrucciones técnicas al taller.' },
      { label: 'Albaranar', text: 'Registra la salida física de mercancía.' }
    ],
    tips: [
      { key: 'F12', text: 'Finalizar y procesar' },
      { key: 'Alt+5', text: 'Ir al terminal de taller' }
    ]
  },
  'General': {
    description: 'Bienvenido al panel central de TramaTex.',
    actions: [],
    tips: [
      { key: 'Ctrl+K', text: 'Búsqueda global' },
      { key: 'Alt+1..5', text: 'Cambiar de módulo' }
    ]
  }
}

const currentHelp = computed(() => helpContent[routeName.value] || helpContent['General'])

function goToFullHelp() {
  router.push('/help')
  emit('close')
}
</script>

<style scoped>
.contextual-help-panel {
  position: fixed; top: 0; right: 0; bottom: 0; width: 360px;
  background: white; border-left: 1px solid var(--color-border);
  z-index: 9999; box-shadow: var(--box-shadow-lg);
  display: flex; flex-direction: column;
}

.panel-header {
  padding: 1.5rem; border-bottom: 1px solid var(--color-border);
  display: flex; justify-content: space-between; align-items: flex-start;
  background: #f8fafc;
}

.title-group h3 { margin: 0; font-size: 1.1rem; font-weight: 800; color: var(--color-text-primary); }
.route-label { font-size: 0.75rem; font-weight: 700; color: var(--color-secondary); text-transform: uppercase; }

.btn-close { background: none; border: none; color: var(--color-text-secondary); cursor: pointer; padding: 0.25rem; }

.panel-body { padding: 1.5rem; overflow-y: auto; flex: 1; display: flex; flex-direction: column; gap: 1.5rem; }

.help-section h4 { font-size: 0.8rem; font-weight: 800; text-transform: uppercase; color: var(--color-text-secondary); margin-bottom: 0.75rem; border-bottom: 1px solid #f1f5f9; padding-bottom: 0.25rem; }
.help-section p { font-size: 0.95rem; line-height: 1.6; color: var(--color-text-primary); }

.actions-list { padding: 0; list-style: none; margin: 0; }
.actions-list li { font-size: 0.9rem; margin-bottom: 0.5rem; }

.tip-card { display: flex; align-items: center; gap: 1rem; padding: 0.75rem; background: #f0f9ff; border: 1px solid #bae6fd; border-radius: 8px; margin-bottom: 0.5rem; }
.tip-card kbd { min-width: 60px; text-align: center; background: white; border: 1px solid #94a3b8; font-size: 0.75rem; }
.tip-card span { font-size: 0.85rem; font-weight: 500; color: #0369a1; }

.panel-footer { margin-top: auto; padding-top: 1rem; }

/* Transitions */
.slide-right-enter-active, .slide-right-leave-active { transition: transform 0.3s cubic-bezier(0.4, 0, 0.2, 1); }
.slide-right-enter-from, .slide-right-leave-to { transform: translateX(100%); }
</style>
