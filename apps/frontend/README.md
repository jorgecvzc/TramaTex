# Frontend TramaTex (Vue + Vite)

## Configuración de emisor para impresión (Sales)

Los documentos de Sales (presupuestos, pedidos, albaranes y facturas) usan un perfil de emisor para cabecera/pie de impresión.

### Variables de entorno (`VITE_`)

Puedes definir estos valores en `.env` (o por entorno):

- `VITE_PRINT_ISSUER_NAME`
- `VITE_PRINT_ISSUER_TAX_LABEL`
- `VITE_PRINT_ISSUER_TAX_ID`
- `VITE_PRINT_ISSUER_ADDRESS`
- `VITE_PRINT_ISSUER_CITY`
- `VITE_PRINT_ISSUER_CONTACT`

Ejemplo:

```env
VITE_PRINT_ISSUER_NAME=TramaTex S.L.
VITE_PRINT_ISSUER_TAX_LABEL=CIF
VITE_PRINT_ISSUER_TAX_ID=B12345678
VITE_PRINT_ISSUER_ADDRESS=C/ Ejemplo 123
VITE_PRINT_ISSUER_CITY=28001 Madrid
VITE_PRINT_ISSUER_CONTACT=info@tramatex.com · +34 900 000 000
```

### Override en runtime (`localStorage`)

La pantalla Admin `Administración / Impresión` guarda un override en:

- clave: `tramatex_print_issuer_profile`

Si existe esa clave, sus valores tienen prioridad sobre `.env` para ese navegador/usuario.

### Precedencia de valores

1. Valores por defecto internos
2. Variables `VITE_PRINT_ISSUER_*`
3. Override de `localStorage` (`tramatex_print_issuer_profile`)
