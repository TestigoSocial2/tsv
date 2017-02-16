# Transparencia Mexicana

### Extensión al Estandar de Contrataciones Abiertas
---

Esta extensión agrega información necesaria de acuerdo a los procedimientos de contratación
pública en México.

### Motivación
---

Facilitar la adoptción del estandar en México.

### Recomendaciones
---

- Agregar un campo `extensions` al objeto que se está modificando, en este caso `release`,
  en donde poder agregar distintas extensiones
- Cada extensión deberá utilizar la notación de dominio inverso para generar un identificador
  y minimizar posibles colisiones de nombre - [Más información](https://en.wikipedia.org/wiki/Reverse_domain_name_notation)
- Cada extensión deberá ser un obejto y contener un campo `version` que indique la versión
  utilizada bajo el estandar de versiones semanticas - [Más información](http://semver.org)

__Ejemplo:__

```json
{
  ...
  "extensions": {
    "mx.org.tm": {
      "version": "0.1.0",
      ...
    }
  }
}
```

### Campos de la Extensión
---

- `planning.budget.multipleYears`  
  __Tipo:__ `bool`  
  __Propósito:__ Indica si el proyecto cuenta con financiamiento por más un ejercicio fiscal

- `planning.budget.exchangeRate`  
  __Tipo:__ `number`  
  __Propósito:__ Cuando los montos estén indicados en una moneda extranjera y se quiera incluir el tipo de cambio a pesos mexicanos utilizado como referencia

- `tender.scope`  
  __Tipo:__ `codelist`  
  __Propósito:__ Nombre del carácter del procedimiento que se utiliza para la contratación. El carácter define quien puede participar con base en el país en el que se encuentra registrada la empresa interesada.

- `tender.procurementStyle`  
  __Tipo:__ `codelist`  
  __Propósito:__ Nombre de la forma del procedimiento que se utiliza para la contratación.

- `tender.hasSocialWitness`  
  __Tipo:__ `bool`  
  __Propósito:__ Campo binario para indicar si hubo o no un testigo social en el proceso de licitación

- `tender.socialWitness`  
  __Tipo:__ `Organization`  
  __Propósito:__ Datos de la entidad que funge como testigo social durante el proceso de contratación

- `tender.requiringEntity`  
  __Tipo:__ `Organization`  
  __Propósito:__ Área solicitante

- `tender.technicalEntity`  
  __Tipo:__ `Organization`  
  __Propósito:__ Área técnica

- `contract.administratorEntity`  
  __Tipo:__ `Organization`  
  __Propósito:__ Área administrativa

- `contract.hasModifications`  
  __Tipo:__ `bool`  
  __Propósito:__ Indica si se realizarón ajustes y/o modificaciones administrativas durante el proceso de contratación

- `implementation.hasAudits`  
  __Tipo:__ `bool`  
  __Propósito:__ Indica si se realizarón auditorias durante el proceso de contratación

- `implementation.audits`  
  __Tipo:__ `Audit array`  
  __Propósito:__ Proporcionar información sobre las auditorias realizadas

### Objeto `Audit`
---

- `type`  
  __Tipo:__ `codelist`  
  __Propósito:__ Especifica el tipo de auditoria

- `auditor`  
  __Tipo:__ `Organization`  
  __Propósito:__ Entidad encargada de llevar a cabo el proceso de auditoria

- `observations`  
  __Tipo:__ `string array`  
  __Propósito:__ Obervaciones realizadas durante el proceso de auditoria

- `actions`  
  __Tipo:__ `string array`  
  __Propósito:__ Acciones realizadas durante el proceso de auditoria

__Ejemplo:__
```json
{
  "type": "internal",
  "auditor": {},
  "observations": ["..."],
  "actions": ["..."]
}
```

## Listas

### `tender.scope`

Valor | Descripción
----- | -----------
national | Licitación Nacional
international | Licitación Internacional
internationalTreaty | Licitación Internacional bajo tratados

### `tender.procurementStyle`
Valor | Descripción
----- | -----------
presence | Licitación Presencial
electronic | Licitación Electrónica
multiple | Licitación Mixta


### `audit.type`
Valor | Descripción
----- | -----------
internal | Auditoria Interna
external | Auditoria Externa

__Ejemplo:__

```json
{
  "version": "0.1.0",
  "planning": {
    "budget": {
      "multipleYears": true,
      "exchangeRate": 20.252
    }
  },
  "tender": {
    "scope": "national",
    "procurementStyle": "multiple",
    "hasSocialWitness": true,
    "socialWitness": {
      "id": "RFC170219PK0",
      "name": "Empresa S.A. de C.V.",
      "designationDate": ""
    },
    "requiringEntity": {},
    "technicalEntity": {}
  },
  "award": {
    "complains": {
      "accepted": 1,
      "notAccepted": 2
    }
  },
  "contract": {
    "administratorEntity": {},
    "hasModifications": false
  },
  "implementation": {
    "hasAudits": true,
    "audits":[{
      "type": "external",
      "auditorEntity": {},
      "observations": [""],
      "actions": [""]
    }]
  }
}
```
