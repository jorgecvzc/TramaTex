# Estrategia Técnica: Asignación de Tareas y Planificación MES (Post-MVP)

Este documento define la inteligencia de planificación del taller, permitiendo optimizar el uso de los recursos humanos y técnicos (máquinas).

---

## 1. Planificador de Carga de Trabajo (Scheduler)

### 1.1 Asignación por Capacidad
- El responsable de planta podrá visualizar la carga de trabajo de cada máquina/operario mediante un **Diagrama de Gantt** interactivo.
- **Balanceo de Carga**: El sistema sugerirá la asignación de tareas a los operarios con menor carga acumulada o mayor especialización en el tipo de trabajo requerido.

---

## 2. Gestión de Operarios y Turnos

### 2.1 Control de Presencia y Disponibilidad
- Gestión de turnos de trabajo, pausas y ausencias.
- El sistema impedirá asignar una tarea a un operario que no esté en su turno de trabajo o que esté de baja/vacaciones.

### 2.2 Skill Mapping y Seguridad por Atributos (ABAC)
- **Matriz de Competencias**: Definición de habilidades por operario (ej: Experto en Máquina X, Junior en Máquina Y).
- **Control de Acceso (ABAC)**: El sistema restringirá la visualización de archivos de diseño o la operación de máquinas según los atributos del operario (Certificaciones vigentes, Nivel de experiencia).
- **Seguridad Industrial**: Impedir el inicio de una fase de producción si el operario no tiene el atributo de "Formación de Seguridad" activo en su ficha.

---

## 3. Registro de Tiempos y Rendimiento

### 3.1 Captura de Tiempos en Tiempo Real
- Los operarios iniciarán y pausarán sus tareas directamente desde el terminal táctil (fichaje industrial).
- **Control de Inactividad**: El sistema requerirá indicar un motivo (ej: Mantenimiento, Falta de Material) si una máquina se detiene durante el turno.

### 3.2 Gamificación y KPIs de Operario
- Visualización para el operario de su rendimiento diario frente a la media o el objetivo de producción, fomentando la mejora continua.

---

*Última actualización: 2026-04-27*
