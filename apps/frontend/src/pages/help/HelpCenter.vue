<template>
  <div class="help-center-page">
    <BasePageHeader 
      title="Centro de Capacitación" 
      :breadcrumbs="[{ label: 'Inicio', to: '/' }, { label: 'Ayuda' }]"
    >
      <template #icon><BookOpen :size="28" /></template>
    </BasePageHeader>

    <div class="help-layout">
      <!-- Barra lateral de navegación de ayuda -->
      <aside class="help-sidebar">
        <nav class="help-nav">
          <button 
            v-for="(block, index) in blocks" 
            :key="block.id"
            class="nav-block-btn"
            :class="{ active: activeBlock === block.id }"
            @click="activeBlock = block.id"
          >
            <span class="block-num">{{ index + 1 }}</span>
            <component :is="block.icon" :size="20" />
            <span>{{ block.title }}</span>
          </button>
        </nav>

        <div class="keyboard-hint mt-8">
          <p>Usa <kbd>1</kbd>-<kbd>4</kbd> para saltar entre bloques</p>
        </div>
      </aside>

      <!-- Contenido Principal -->
      <main class="help-content-area">
        <Transition name="fade-slide" mode="out-in">
          <div :key="activeBlock" class="block-content">
            <header class="block-header">
              <component :is="currentBlock.icon" :size="32" class="text-primary" />
              <div>
                <h2>{{ currentBlock.title }}</h2>
                <p class="subtitle">{{ currentBlock.description }}</p>
              </div>
            </header>

            <div class="articles-grid">
              <section v-for="article in currentBlock.articles" :key="article.title" class="help-card">
                <h3>{{ article.title }}</h3>
                <p>{{ article.content }}</p>
                
                <div v-if="article.steps" class="steps-list">
                  <div v-for="(step, sIdx) in article.steps" :key="sIdx" class="step-item">
                    <span class="step-badge">{{ sIdx + 1 }}</span>
                    <span>{{ step }}</span>
                  </div>
                </div>

                <div v-if="article.shortcuts" class="shortcuts-mini-grid">
                  <div v-for="sc in article.shortcuts" :key="sc.key" class="mini-sc">
                    <kbd>{{ sc.key }}</kbd>
                    <span>{{ sc.label }}</span>
                  </div>
                </div>
              </section>
            </div>
          </div>
        </Transition>
      </main>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onBeforeUnmount } from 'vue'
import { 
  BookOpen, 
  Info, 
  CreditCard, 
  Factory, 
  Keyboard, 
  Layout, 
  Layers, 
  Zap,
  MousePointer2,
  ArrowRightLeft
} from 'lucide-vue-next'
import BasePageHeader from '@/components/shared/BasePageHeader.vue'

const activeBlock = ref('foundations')

const blocks = [
  {
    id: 'foundations',
    title: 'Fundamentos',
    icon: Layout,
    description: 'Conceptos básicos y arquitectura visual de TramaTex.',
    articles: [
      {
        title: 'La Interfaz en 3 Capas',
        content: 'Toda página de gestión se divide en: 1. Identidad (Fija arriba), 2. Contexto (KPIs y trazabilidad) y 3. Trabajo (Formularios y líneas).',
      },
      {
        title: 'Maestros y Variantes JIT',
        content: 'No creamos miles de variantes. El sistema genera la variante exacta (Talla/Color) en el momento de la venta (Just-In-Time).',
      }
    ]
  },
  {
    id: 'sales',
    title: 'Ciclo de Ventas',
    icon: CreditCard,
    description: 'Desde la oferta inicial hasta el cobro final.',
    articles: [
      {
        title: 'Flujo Documental',
        content: 'El camino recomendado para mantener la trazabilidad es:',
        steps: [
          'Crear Presupuesto para el cliente.',
          'Confirmar Presupuesto -> Se genera Pedido de Venta.',
          'Generar Albarán desde el pedido para logística.',
          'Emitir Factura final.'
        ]
      },
      {
        title: 'Pricing Inteligente',
        content: 'El sistema calcula precios basándose en: Tarifa base + Modificadores de atributo (ej: XL +1€) - Descuento de cliente.',
      }
    ]
  },
  {
    id: 'mes',
    title: 'Operativa MES',
    icon: Factory,
    description: 'Gestión de taller y configuraciones técnicas.',
    articles: [
      {
        title: 'Configuración Técnica',
        content: 'En un pedido, puedes añadir "Trabajos MES". Selecciona un Tipo de Trabajo (ej: Serigrafía) y una Posición (ej: Espalda).',
      },
      {
        title: 'Terminal de Taller',
        content: 'Diseñado para pantallas táctiles. Los operarios ven las tareas pendientes y marcan inicio/fin para control de tiempos real.',
      }
    ]
  },
  {
    id: 'keyboard',
    title: 'Masterclass Teclado',
    icon: Keyboard,
    description: 'Cómo operar al 100% sin usar el ratón.',
    articles: [
      {
        title: 'Navegación de Líneas',
        content: 'Domina las tablas de productos con estas teclas:',
        shortcuts: [
          { key: 'Enter', label: 'Siguiente campo / Nueva línea' },
          { key: '↑ / ↓', label: 'Subir / Bajar valores' },
          { key: 'Insert', label: 'Abrir buscador rápido' },
          { key: 'Ctrl+Supr', label: 'Eliminar línea actual' }
        ]
      },
      {
        title: 'Atajos Globales',
        content: 'Funciona en cualquier parte del sistema.',
        shortcuts: [
          { key: 'Alt+1..5', label: 'Cambiar de módulo' },
          { key: 'Alt+H', label: 'Menú de Ayuda' },
          { key: 'Ctrl+S', label: 'Guardar formulario' },
          { key: 'Ctrl+B', label: 'Colapsar Sidebar' },
          { key: 'F1', label: 'Ayuda Contextual' }
        ]
      }
    ]
  }
]

const currentBlock = computed(() => blocks.find(b => b.id === activeBlock.value)!)

function handleGlobalKeydown(e: KeyboardEvent) {
  const isInput = ['INPUT', 'TEXTAREA'].includes(document.activeElement?.tagName || '')
  if (isInput) return

  if (e.key >= '1' && e.key <= '4') {
    activeBlock.value = blocks[parseInt(e.key) - 1].id
  }
}

onMounted(() => {
  window.addEventListener('keydown', handleGlobalKeydown)
  // Manejar hash si venimos de un enlace directo
  if (window.location.hash === '#glossary') {
    activeBlock.value = 'foundations'
  }
})

onBeforeUnmount(() => {
  window.removeEventListener('keydown', handleGlobalKeydown)
})
</script>

<style scoped>
.help-center-page {
  padding: 1.5rem;
  background: var(--color-background);
  min-height: 100vh;
}

.help-layout {
  display: grid;
  grid-template-columns: 280px 1fr;
  gap: 2rem;
  margin-top: 2rem;
  max-width: 1400px;
  margin-left: auto;
  margin-right: auto;
}

/* Sidebar */
.help-sidebar {
  position: sticky;
  top: 100px;
  height: fit-content;
}

.help-nav {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.nav-block-btn {
  display: flex;
  align-items: center;
  gap: 1rem;
  padding: 1rem;
  background: white;
  border: 1px solid var(--color-border);
  border-radius: 12px;
  cursor: pointer;
  transition: all 0.2s;
  text-align: left;
  color: var(--color-text-secondary);
}

.nav-block-btn:hover {
  border-color: var(--color-primary);
  background: #fffbeb;
}

.nav-block-btn.active {
  background: var(--color-secondary);
  border-color: var(--color-secondary);
  color: white;
  box-shadow: var(--box-shadow-md);
}

.block-num {
  width: 24px; height: 24px;
  display: flex; align-items: center; justify-content: center;
  background: rgba(0,0,0,0.05);
  border-radius: 6px;
  font-size: 0.75rem;
  font-weight: 800;
}
.active .block-num { background: rgba(255,255,255,0.2); }

.keyboard-hint {
  padding: 1rem;
  background: #f8fafc;
  border-radius: 8px;
  font-size: 0.8rem;
  color: var(--color-text-secondary);
  text-align: center;
}

/* Content Area */
.help-content-area {
  min-width: 0;
}

.block-header {
  display: flex;
  align-items: center;
  gap: 1.5rem;
  margin-bottom: 2.5rem;
  padding-bottom: 1.5rem;
  border-bottom: 2px solid var(--color-border-soft);
}

.block-header h2 { margin: 0; font-size: 2rem; color: var(--color-text-primary); }
.subtitle { margin: 0.25rem 0 0; color: var(--color-text-secondary); font-size: 1.1rem; }

.articles-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(400px, 1fr));
  gap: 1.5rem;
}

.help-card {
  background: white;
  padding: 1.5rem;
  border-radius: 16px;
  border: 1px solid var(--color-border);
  box-shadow: var(--box-shadow-sm);
}

.help-card h3 {
  margin: 0 0 1rem;
  font-size: 1.25rem;
  color: var(--color-secondary);
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.help-card p {
  color: var(--color-text-primary);
  line-height: 1.6;
  font-size: 1rem;
  margin-bottom: 1.25rem;
}

/* Steps */
.steps-list { display: flex; flex-direction: column; gap: 0.75rem; }
.step-item { display: flex; align-items: center; gap: 1rem; font-size: 0.95rem; font-weight: 500; }
.step-badge {
  width: 24px; height: 24px; flex-shrink: 0;
  background: var(--color-primary); color: black;
  border-radius: 50%; display: flex; align-items: center; justify-content: center;
  font-size: 0.75rem; font-weight: 900;
}

/* Mini Shortcuts */
.shortcuts-mini-grid {
  display: grid;
  grid-template-columns: 1fr;
  gap: 0.5rem;
  background: #f1f5f9;
  padding: 1rem;
  border-radius: 12px;
}
.mini-sc { display: flex; justify-content: space-between; align-items: center; }
.mini-sc kbd { font-size: 0.7rem; min-width: 80px; text-align: center; }
.mini-sc span { font-size: 0.85rem; color: var(--color-text-secondary); }

/* Transitions */
.fade-slide-enter-active, .fade-slide-leave-active { transition: all 0.3s ease; }
.fade-slide-enter-from { opacity: 0; transform: translateY(10px); }
.fade-slide-leave-to { opacity: 0; transform: translateY(-10px); }

@media (max-width: 1024px) {
  .help-layout { grid-template-columns: 1fr; }
  .help-sidebar { position: static; }
  .help-nav { flex-direction: row; overflow-x: auto; padding-bottom: 1rem; }
  .nav-block-btn { white-space: nowrap; }
}
</style>
