# Plantilla de Documentación de Módulo

**Nombre del Módulo:** [Nombre]  
**Bounded Context:** [Contexto Delimitado]  
**Responsabilidad Principal:** [Descripción breve]  
**Entidades Raíz:** [Agregados principales]  
**Dependencias:** [Módulos de los que depende]  

---

## 1. Especificación del Módulo (module-spec.md)

### Objetivo
[Describir qué problema resuelve y qué valor proporciona]

### Alcance
[Qué está incluido, qué no]

### Restricciones
[Limitaciones técnicas o de negocio]

---

## 2. Modelo de Dominio (domain-model.md)

### Entidades Principales

```
[DiagramER simplificado]
```

#### [Entidad 1]
- Responsabilidad:
- Value Objects:
- Reglas de Negocio:

#### [Entidad 2]
- ...

### Value Objects
- [VO1]: [Descripción]
- [VO2]: [Descripción]

### Servicios de Dominio
[Si aplica]

---

## 3. Casos de Uso (use-cases.md)

### Caso de Uso 1: [Nombre]
- **Actor:** [Quién lo usa]
- **Precondiciones:** [Qué debe ser verdadero]
- **Flujo Normal:**
  1. ...
  2. ...
- **Flujos Alternativos:** [Si aplica]
- **Postcondiciones:** [Qué debe ser verdadero después]

### Caso de Uso 2: [Nombre]
- ...

---

## 4. Contratos de API (api-contracts.md)

### Endpoint 1
```json
POST /api/[modulo]/[recurso]
Request: { ... }
Response: { ... }
Errores: [...]
```

### Endpoint 2
```
...
```

---

## 5. Decisiones Técnicas

### [Decisión 1]
**Alternativas Consideradas:**
**Decisión Tomada:**
**Justificación:**

---

## 6. Tests

### Cobertura Objetivo
- Dominio: ≥80%
- Aplicación: ≥70%
- Infraestructura: ≥60%

### Casos de Prueba Críticos
- [Test 1]
- [Test 2]

---

## 7. Notas y Pendientes

- [ ] Tarea 1
- [ ] Tarea 2

---

**Última Actualización:** [Fecha]  
**Responsable:** [Persona/Equipo]
