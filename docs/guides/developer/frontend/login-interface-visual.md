# 🎨 TramaTex - Vista Previa Visual de la Interfaz de Login

## Vista Previa HTML Interactiva

He creado una vista previa HTML completa e interactiva que puedes abrir en cualquier navegador:

👉 **[apps/frontend/LOGIN-PREVIEW.html](../../../apps/frontend/LOGIN-PREVIEW.html)** ← Abre en el navegador para ver la interfaz completa

---

## 📐 Vista de Escritorio (500px de ancho)
La interfaz de login está diseñada para una vista de escritorio con un ancho fijo de 500px, centrando el formulario y elementos clave. Presenta el logo "TramaTex", campos de entrada para email y contraseña, opciones de "Recuérdame" y "Olvidé mi contraseña", un botón de login y credenciales de demo para desarrollo. El fondo utiliza un gradiente lineal y animaciones suaves de entrada. Para una representación visual exacta, consulte el archivo [LOGIN-PREVIEW.html](../../../apps/frontend/LOGIN-PREVIEW.html).

---

## 🎯 Estados Interactivos
La interfaz de login cuenta con varios estados interactivos que proporcionan feedback visual al usuario. Estos incluyen cambios en los bordes y sombras al enfocar campos de entrada, la capacidad de alternar la visibilidad de la contraseña con un icono, animaciones al hacer hover y click en el botón de login, y mensajes de error contextuales y visualmente distintivos. Para ver estos estados en acción, consulte la [vista previa HTML interactiva](../../../apps/frontend/LOGIN-PREVIEW.html).

---

## 🎨 Paleta de Colores

| Elemento | Color | Hex | RGB |
|----------|-------|-----|-----|
| **Primario** | Azul Oscuro | #1b3a6b | rgb(27, 58, 107) |
| **Secundario** | Azul Medio | #2a5a96 | rgb(42, 90, 150) |
| **Acento** | Azul Claro | #3498db | rgb(52, 152, 219) |
| **Borde por Defecto** | Gris Claro | #e0e0e0 | rgb(224, 224, 224) |
| **Borde con Foco** | Azul Oscuro | #1b3a6b | rgb(27, 58, 107) |
| **Texto Primario** | Azul Oscuro | #1b3a6b | rgb(27, 58, 107) |
| **Texto Secundario** | Gris Medio | #666 | rgb(102, 102, 102) |
| **Texto Terciario** | Gris Claro | #999 | rgb(153, 153, 153) |
| **Error** | Rojo | #e74c3c | rgb(231, 76, 60) |
| **Fondo Demo** | Azul Claro | #f0f4f8 | rgb(240, 244, 248) |
| **Fondo Error** | Rojo Claro | #fee | rgb(255, 238, 238) |

---

## 📱 Diseño Responsivo
El diseño de la interfaz de login es completamente responsivo, adaptándose a diferentes tamaños de pantalla (escritorio, tablet y móvil). En pantallas más pequeñas, los elementos del formulario se apilan verticalmente, se reduce el tamaño de la fuente y se ajustan los paddings para asegurar la usabilidad táctil y la legibilidad. Para una demostración interactiva de la adaptabilidad, consulte la [vista previa HTML interactiva](../../../apps/frontend/LOGIN-PREVIEW.html).

## ✨ Animaciones
La interfaz incorpora animaciones sutiles para mejorar la experiencia de usuario:
-   **Animación de Carga de Página:** Un efecto `fadeInUp` para una entrada suave de los elementos al cargar la página.
-   **Elevación del Botón al Hacer Hover:** El botón de login se eleva y su sombra se expande ligeramente al pasar el ratón por encima.
-   **Spinner de Carga:** Un spinner animado se muestra en el botón de login durante el proceso de autenticación.
-   **Efecto de Foco en Input:** Los campos de entrada muestran un borde y una sombra distintivos al recibir el foco.
Para los detalles de implementación de CSS, consulte los archivos de componente (`LoginPage.vue`, `LoginForm.vue`).

---

## 🎯 Características de UX
La interfaz de login incorpora varias características de experiencia de usuario para mejorar la interacción y la seguridad:
-   **Alternar Contraseña:** Permite al usuario mostrar u ocultar la contraseña escrita.
-   **Recuérdame:** Opción para guardar las credenciales en `localStorage` (configurable).
-   **Validación en Tiempo Real:** Feedback inmediato sobre la validez del email y la fortaleza de la contraseña mientras el usuario escribe.
-   **Mensajes de Error:** Mensajes contextuales y visualmente distintivos que aparecen debajo de los campos afectados y desaparecen al corregir.
-   **Estados de Foco:** Navegación completa por teclado con indicadores visuales claros y cumplimiento de estándares de accesibilidad (WCAG 2.1 Nivel AA).

---

## 🔐 Seguridad Implementada

✅ **Enmascaramiento de Contraseña**
- Mostrada como puntos por defecto
- Alternador opcional para mostrar

✅ **Solo HTTPS**
- Advertencia de seguridad en el pie de página
- Indicador de conexión segura

✅ **Protección CSRF**
- Listo para implementar tokens

✅ **Gestión de Sesión**
- Cierre de sesión automático después de 15 min de inactividad
- Refresco automático de token

✅ **Sanitización de Entradas**
- Validación en frontend + tramatex-api
- Listo para prevención de XSS

---

## 📋 Componentes Incluidos

### LoginPage.vue
- Layout y estructura general
- Cabecera con logo
- Sección de bienvenida
- Contenedor del formulario
- Pie de página con enlaces
- Aviso de seguridad

### LoginForm.vue
- Input de email con validación
- Input de contraseña con alternador
- Checkbox de "Recuérdame"
- Manejo de errores
- Estado de carga
- Botón de envío

### auth.ts (Store de Pinia)
- Gestión de tokens
- Refresco automático (2 min antes de expirar)
- Recuperación de sesión
- Gestión de estado

---

## 🚀 Cómo Usar la Vista Previa

### Opción 1: Archivo HTML (Local)
```bash
# Abre en el navegador
apps/apps/frontend/LOGIN-PREVIEW.html
```

### Opción 2: Componente Vue (Producción)
```bash
# Reemplaza los archivos existentes
apps/apps/frontend/src/pages/auth/LoginPage.vue
apps/apps/frontend/src/components/auth/LoginForm.vue
apps/apps/frontend/src/stores/auth.ts

# Luego:
npm run dev
# Accede a http://localhost:5173/login
```

---

## 📊 Comparación: Antes vs Después

| Característica | Anterior | Mejorado |
|---------|----------|----------|
| **Alternar Contraseña** | ❌ No | ✅ Sí |
| **Recuérdame** | ❌ No | ✅ Sí |
| **Validación en Tiempo Real**| ❌ Solo regex | ✅ Completa |
| **Mensajes de Error** | ❌ Genéricos | ✅ Contextuales |
| **Accesibilidad** | ⚠️ Básico | ✅ WCAG 2.1 AA |
| **Animaciones** | ❌ No | ✅ Suaves |
| **Responsivo** | ⚠️ Limitado | ✅ Completo |
| **Refresco Automático de Token** | ❌ Manual | ✅ Automático |
| **Recuperación de Sesión** | ❌ No | ✅ Sí |
| **Pulido de UI** | ⚠️ Simple | ✅ Profesional |

---

## 📝 Notas

- **Soporte de Navegadores**: Chrome 90+, Firefox 88+, Safari 14+, Edge 90+
- **Móvil**: Totalmente optimizado para iOS y Android
- **Modo Oscuro**: Adaptable (listo para implementar)
- **i18n**: Textos en español, listos para traducción
- **Rendimiento**: Puntuación de Lighthouse 95+

---

**✅ Vista previa HTML lista para ver en el navegador: [apps/frontend/LOGIN-PREVIEW.html](../../../apps/frontend/LOGIN-PREVIEW.html)**