<template>
  <div class="dev-page-wrapper">
    <Navbar />
    
    <div class="dev-container">
      <header class="dev-header">
        <h1 class="font-brand">TramaTex Design System</h1>
        <p>Librería de componentes y estándares visuales para desarrolladores.</p>
      </header>

      <nav class="dev-nav card">
        <a href="#colors">Colores</a>
        <a href="#typography">Tipografía</a>
        <a href="#buttons">Botones</a>
        <a href="#tags">Tags y Trazabilidad</a>
        <a href="#forms">Formularios</a>
        <a href="#license">Licencia</a>
      </nav>

      <div class="dev-content">
        <!-- SECCIÓN: COLORES -->
        <section id="colors" class="card">
          <div class="section-header-dev">
            <span class="material-symbols-outlined">palette</span>
            <h2>Paleta de Colores</h2>
          </div>
          
          <div class="color-section-group">
            <h3>Colores de Marca</h3>
            <div class="color-grid brand-grid">
              <div class="color-swatch-wrapper" v-for="color in brandColors" :key="color.name">
                <div class="color-swatch" :style="{ backgroundColor: color.var }"></div>
                <p>{{ color.name }}</p>
                <code>{{ color.var }}</code>
              </div>
            </div>
          </div>

          <div class="color-section-group mt-8">
            <h3>Colores de Estado y Sistema</h3>
            <div class="color-grid system-grid">
              <div class="color-swatch-wrapper" v-for="color in systemColors" :key="color.name">
                <div class="color-swatch" :style="{ backgroundColor: color.var }"></div>
                <p>{{ color.name }}</p>
                <code>{{ color.var }}</code>
              </div>
            </div>
          </div>
        </section>

        <!-- SECCIÓN: TIPOGRAFÍA -->
        <section id="typography" class="card">
          <div class="section-header-dev">
            <span class="material-symbols-outlined">text_fields</span>
            <h2>Tipografía</h2>
          </div>
          <div class="typography-showcase">
            <h1 class="font-brand" style="color: var(--color-secondary); font-size: 3rem;">TramaTex Brand Font</h1>
            <hr>
            <h1>Heading 1 - Títulos de Página</h1>
            <h2>Heading 2 - Títulos de Sección</h2>
            <h3>Heading 3 - Títulos de Tarjeta</h3>
            <p>Cuerpo de texto estándar para descripciones y párrafos largos.</p>
            <p class="text-muted">Texto secundario o deshabilitado (Muted).</p>
            <code>Código monoespaciado para IDs, SKUs y Referencias</code>
          </div>
        </section>

        <!-- SECCIÓN: BOTONES -->
        <section id="buttons" class="card">
          <div class="section-header-dev">
            <span class="material-symbols-outlined">smart_button</span>
            <h2>Botones Estandarizados</h2>
          </div>
          <table class="dev-table">
            <thead>
              <tr>
                <th>Tipo</th>
                <th>Muestra</th>
                <th>Con Icono</th>
                <th>Deshabilitado</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="btn in buttonTypes" :key="btn.class">
                <td>
                  <strong>{{ btn.label }}</strong><br>
                  <small>{{ btn.desc }}</small>
                </td>
                <td><button :class="['btn', btn.class]">{{ btn.text }}</button></td>
                <td>
                  <button :class="['btn', btn.class]">
                    <span class="material-symbols-outlined">{{ btn.icon }}</span>
                    <span>{{ btn.text }}</span>
                  </button>
                </td>
                <td><button :class="['btn', btn.class]" disabled>{{ btn.text }}</button></td>
              </tr>
            </tbody>
          </table>
        </section>

        <!-- SECCIÓN: TAGS Y TRAZABILIDAD -->
        <section id="tags" class="card">
          <div class="section-header-dev">
            <span class="material-symbols-outlined">label</span>
            <h2>Tags y Trazabilidad</h2>
          </div>
          
          <div class="dev-grid-2">
            <!-- Summary Tags -->
            <div class="dev-example-box">
              <h3>Summary Tags (KPIs)</h3>
              <p class="text-muted mb-4">Uso: Resumen de datos clave. Sin bordes, solo sombra y color de icono.</p>
              <div class="flex-col gap-4">
                <div class="summary-tag">
                  <div class="icon blue"><span class="material-symbols-outlined">person</span></div>
                  <div class="tag-content"><label>Cliente</label><strong>Textil del Norte S.L.</strong></div>
                </div>
                <div class="summary-tag">
                  <div class="icon green"><span class="material-symbols-outlined">payments</span></div>
                  <div class="tag-content"><label>Total</label><strong>1.240,50 €</strong></div>
                </div>
              </div>
            </div>

            <!-- Related Tags -->
            <div class="dev-example-box">
              <h3>Related Tags (Navegación)</h3>
              <p class="text-muted mb-4">Uso: Trazabilidad y enlaces. Con borde izquierdo de color e icono de salto.</p>
              <div class="flex-col gap-4">
                <div class="related-tag-card highlight-info">
                  <div class="tag-icon"><span class="material-symbols-outlined">description</span></div>
                  <div class="tag-content"><label>Origen</label><strong>Presupuesto #45</strong></div>
                  <span class="material-symbols-outlined jump-icon">open_in_new</span>
                </div>
                <div class="related-tag-card success">
                  <div class="tag-icon success"><span class="material-symbols-outlined">receipt</span></div>
                  <div class="tag-content"><label>Destino</label><strong>Factura #A-102</strong></div>
                  <span class="material-symbols-outlined jump-icon">open_in_new</span>
                </div>
              </div>
            </div>
          </div>
        </section>

        <!-- SECCIÓN: FORMULARIOS -->
        <section id="forms" class="card">
          <div class="section-header-dev">
            <span class="material-symbols-outlined">edit_note</span>
            <h2>Estructura de Datos y Formularios</h2>
          </div>
          
          <div class="dev-info-box mb-8">
            <p>TramaTex utiliza una estructura jerárquica para la entrada y visualización de datos: 
               <strong>FormSection</strong> (Contenedor) > <strong>Layout</strong> (Grid/Row) > <strong>DataRow/Input</strong> (Elemento).</p>
          </div>

          <div class="dev-grid-2">
            <!-- Ejemplo de FormSection + DataRow (Lectura) -->
            <div class="dev-example-box">
              <h3>Modo Lectura (DataRow)</h3>
              <div class="standard-border-p">
                <FormSection title="Sección de Ejemplo" icon="info">
                  <DataRow label="Nombre del Objeto" value="Referencia Estándar" icon="tag" />
                  <DataRow label="Estado Actual" icon="verified">
                    <span class="status-badge status-success">Activo</span>
                  </DataRow>
                  <DataRow label="Observaciones" value="Texto largo de ejemplo que demuestra el ajuste de línea en las fichas de detalle." />
                </FormSection>
              </div>
            </div>

            <!-- Ejemplo de Edición -->
            <div class="dev-example-box">
              <h3>Modo Edición (Layouts)</h3>
              <div class="standard-border-p">
                <FormSection title="Formulario de Edición" icon="edit">
                  <div class="form-row">
                    <div class="form-group">
                      <label>Campo A *</label>
                      <input type="text" class="form-input" placeholder="Obligatorio">
                    </div>
                    <div class="form-group">
                      <label>Campo B (Fecha)</label>
                      <input type="date" class="form-input">
                    </div>
                  </div>
                  <div class="form-group mt-4">
                    <label>Área de Texto</label>
                    <textarea class="form-textarea" rows="2" placeholder="Notas..."></textarea>
                  </div>
                </FormSection>
              </div>
            </div>
          </div>

          <!-- Selectores Avanzados -->
          <div class="dev-section-group mt-12">
            <h3>Selectores Especializados</h3>
            <div class="dev-grid-2 mt-4">
              <div class="dev-card-mini">
                <div class="icon-circle"><span class="material-symbols-outlined">person_search</span></div>
                <div>
                  <strong>PartySelector</strong>
                  <p class="text-muted">Búsqueda reactiva de Clientes/Proveedores con filtros de rol.</p>
                </div>
              </div>
              <div class="dev-card-mini">
                <div class="icon-circle"><span class="material-symbols-outlined">inventory_2</span></div>
                <div>
                  <strong>VariantSelector</strong>
                  <p class="text-muted">Explorador de catálogo con soporte para variantes JIT y stock.</p>
                </div>
              </div>
            </div>
          </div>
        </section>

        <!-- SECCIÓN: LICENCIA -->
        <section id="license" class="card">
          <div class="section-header-dev">
            <span class="material-symbols-outlined">copyright</span>
            <h2>Licencia y Autoría</h2>
          </div>
          <div class="license-content">
            <p><strong>TramaTex</strong> es un proyecto MIT de código abierto.</p>
            <div class="license-box">
              <p>Copyright (c) 2026 Jorge Cortés Villalba</p>
              <p>Este sistema de diseño es propiedad intelectual del autor y se distribuye como parte del núcleo de TramaTex.</p>
            </div>
          </div>
        </section>
      </div>
    </div>
  </div>
</template>

<script setup>
import { RouterLink } from 'vue-router';
import Navbar from '@/components/layout/Navbar.vue';
import DataRow from '@/components/shared/DataRow.vue';

const brandColors = [
  { name: 'Primary (Brand)', var: 'var(--color-primary)' },
  { name: 'Primary Light', var: 'var(--color-primary-light)' },
  { name: 'Secondary', var: 'var(--color-secondary)' },
  { name: 'Secondary Light', var: 'var(--color-secondary-light)' },
];

const systemColors = [
  { name: 'Success', var: 'var(--color-success)' },
  { name: 'Warning', var: 'var(--color-warning)' },
  { name: 'Error', var: 'var(--color-error)' },
  { name: 'Background', var: 'var(--color-background)' },
  { name: 'Surface', var: 'var(--color-surface)' },
];

const buttonTypes = [
  { label: 'Primary (Action)', class: 'btn-primary', text: 'Editar / Nuevo', icon: 'edit', desc: 'Habilitar edición o iniciar la creación de un nuevo objeto.' },
  { label: 'Secondary (Commit)', class: 'btn-secondary', text: 'Guardar / Generar', icon: 'save', desc: 'Persistir cambios o realizar procesos técnicos de salida (PDF).' },
  { label: 'Success (Drive)', class: 'btn-success', text: 'Confirmar / Albarán', icon: 'local_shipping', desc: 'Impulsar el documento al siguiente estado del flujo de negocio.' },
  { label: 'Nav (Link)', class: 'btn-nav', text: 'Volver', icon: 'arrow_back', desc: 'Navegación entre entidades sin modificar datos.' },
  { label: 'Outline', class: 'btn-outline', text: 'Imprimir', icon: 'print', desc: 'Herramientas de apoyo o acciones secundarias.' },
  { label: 'Danger', class: 'btn-danger', text: 'Anular / Eliminar', icon: 'cancel', desc: 'Acciones de ruptura, cancelación o borrado.' },
];
</script>

<style scoped>
.dev-page-wrapper { background-color: var(--color-background); min-height: 100vh; }
.dev-container { max-width: 1200px; margin: 0 auto; padding: 3rem 1rem; }
.dev-header { margin-bottom: 3rem; text-align: center; }
.dev-header h1 { font-size: 3rem; color: var(--color-secondary); margin-bottom: 0.5rem; }

.dev-nav { display: flex; justify-content: center; gap: 1.5rem; padding: 1.25rem; position: sticky; top: 80px; z-index: 10; margin-bottom: 3rem; border-radius: 100px; }
.dev-nav a { text-decoration: none; font-weight: 800; color: var(--color-text-secondary); font-size: 0.8rem; text-transform: uppercase; letter-spacing: 0.05em; }
.dev-nav a:hover { color: var(--color-primary); }

.dev-content { display: flex; flex-direction: column; gap: 3rem; }
section { padding: 2.5rem; }

.section-header-dev { display: flex; align-items: center; gap: 1rem; margin-bottom: 2rem; border-bottom: 2px solid var(--color-primary); padding-bottom: 0.75rem; }
.section-header-dev h2 { margin: 0; text-transform: uppercase; letter-spacing: 0.05em; font-size: 1.25rem; }
.section-header-dev .material-symbols-outlined { color: var(--color-primary); font-size: 32px; }

.color-section-group h3 { font-size: 0.9rem; color: var(--color-text-secondary); text-transform: uppercase; margin-bottom: 1.5rem; border-left: 3px solid var(--color-primary); padding-left: 0.75rem; }

.color-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(140px, 1fr)); gap: 1.5rem; }
.color-swatch { height: 80px; border-radius: 12px; margin-bottom: 0.5rem; border: 1px solid var(--color-border); box-shadow: var(--box-shadow-sm); }

.dev-table { width: 100%; border-collapse: collapse; }
.dev-table th { text-align: left; padding: 1rem; background: var(--color-background); font-size: 0.75rem; text-transform: uppercase; }
.dev-table td { padding: 1.25rem 1rem; border-bottom: 1px solid var(--color-border); }

.dev-grid-2 { display: grid; grid-template-columns: 1fr 1fr; gap: 2rem; }
.dev-example-box { background: #f8fafc; padding: 1.5rem; border-radius: 12px; border: 1px solid var(--color-border); }
.dev-example-box h3 { margin-top: 0; margin-bottom: 1.25rem; font-size: 0.85rem; color: var(--color-text-secondary); text-transform: uppercase; }

/* REPLICATED COMPONENT STYLES FOR THE SHOWCASE */
.summary-tag { background: white; border: 1px solid var(--color-border); border-radius: 12px; padding: 0.6rem 1.25rem; display: flex; align-items: center; gap: 1rem; box-shadow: var(--box-shadow-sm); }
.summary-tag .icon { width: 40px; height: 40px; border-radius: 10px; display: flex; align-items: center; justify-content: center; }
.summary-tag .icon.blue { background: rgba(59, 130, 246, 0.1); color: #2563eb; }
.summary-tag .icon.green { background: rgba(34, 197, 94, 0.1); color: #16a34a; }
.tag-content { display: flex; flex-direction: column; gap: 0.15rem; line-height: 1.2; }
.tag-content label { font-size: 0.65rem; font-weight: 700; text-transform: uppercase; color: var(--color-text-secondary); letter-spacing: 0.025em; }
.tag-content strong { font-size: 0.95rem; color: var(--color-text-primary); }

.related-tag-card { padding: 0.6rem 1rem; background: white; border: 1px solid var(--color-border); border-left: 4px solid var(--color-secondary); border-radius: 10px; display: flex; align-items: center; gap: 0.75rem; transition: all 0.2s; }
.related-tag-card.highlight-info { border-left-color: #2563eb; }
.related-tag-card.success { border-left-color: #16a34a; }
.tag-icon { width: 36px; height: 36px; border-radius: 8px; background: #f1f5f9; display: flex; align-items: center; justify-content: center; }
.jump-icon { font-size: 16px; color: var(--color-border); margin-left: auto; }

.license-box { background: var(--color-background); padding: 1.5rem; border-radius: 12px; border-left: 4px solid var(--color-secondary); margin-top: 1rem; }
.flex-col { display: flex; flex-direction: column; }
.gap-4 { gap: 1rem; }
.mb-4 { margin-bottom: 1rem; }
.mt-8 { margin-top: 2rem; }

@media (max-width: 1024px) { .dev-grid-2 { grid-template-columns: 1fr; } }
@media (max-width: 768px) { .dev-nav { display: none; } }
</style>
