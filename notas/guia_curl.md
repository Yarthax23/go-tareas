# 📄 Referencia Técnica: Uso de curl para el servidor REST

Este archivo documenta cómo interactuar con los endpoints de la API REST desarrollada en Go, utilizando `curl`. Sirve como guía rápida de referencia para pruebas manuales y depuración.

---

## 🧪 Cómo probar los endpoints con `curl`

### 📌 1. Crear un nuevo post (`POST`)

#### JSON:

```bash
curl -X POST localhost:8080/post \
  -H "Content-Type: application/json" \
  -d '{"contenido": "Un posteo"}'
```

#### Texto plano:

```bash
curl -X POST localhost:8080/post \
  -H "Content-Type: text/plain" \
  --data "Este es un post en texto plano"
```

---

### 📌 2. Ver todos los posteos (`GET`)

```bash
curl localhost:8080/ver
curl localhost:8080/ver | jq
```

---

### 📌 3. Ver un post por ID (`GET`)

```bash
curl localhost:8080/ver/3
```

---

### 📌 4. Editar un post existente (`PATCH` o `PUT`)

```bash
curl -X PATCH localhost:8080/editar/8 \
  -H "Content-Type: application/json" \
  -d '{"contenido": "Un posteo editado"}'
```

> PATCH modifica parcialmente.
> PUT puede usarse si reemplazás todos los campos.

---

### 📌 5. Borrar un post por ID (`DELETE`)

```bash
curl -X DELETE localhost:8080/borrar/3
```

---

### 📌 6. Alternar campo `listo` (toggle)

```bash
curl -X POST localhost:8080/toggle/4
```

Este endpoint cambia el campo `listo` de 'SI' a 'NO' o viceversa.


## 📝 Notas técnicas

* Todos los datos JSON deben usar comillas dobles (`"`) y estar correctamente escapados si se escriben en consola.
* El `Content-Type` debe coincidir con el tipo de dato que se está enviando.
* Si se olvida el `Content-Type`, el servidor puede rechazar la solicitud.
* Los IDs deben ser enteros válidos y corresponder a registros existentes en la base de datos.

---

## 🔧 Endpoints esperados

| Método | Ruta             | Descripción                     |
| ------ | ---------------- | ------------------------------- |
| POST   | /post            | Crear post (JSON o texto plano) |
| POST   | /toggle/{id}     | Cambiar estado Listo            |
| GET    | /ver             | Listar todos los posts          |
| GET    | /ver/{id}        | Ver post específico por ID      |
| DELETE | /borrar/{id}     | Eliminar post por ID            |
| PUT    | /actualizar/{id} | Reemplazar post completamente   |
| PATCH  | /editar/{id}     | Editar parcialmente un post     |


## Base de datos

Tabla `posts(id SERIAL PRIMARY KEY, contenido TEXT, listo TEXT)`

* El campo `listo` se agregó con:

```sql
ALTER TABLE posts ADD COLUMN listo TEXT DEFAULT 'NO';
```

* Las filas viejas se actualizaron con:

```sql
UPDATE posts SET listo = 'NO' WHERE listo IS NULL;
```

---

## Futuras ideas

* Endpoint `/ver/contenido` si se desea sólo ver texto sin otros campos.
* Agregar campos como fecha de creación, categoría, etiquetas, etc.
* Endpoint para borrar todos los posts (`DELETE /borrar-todo`)
* Autenticación de usuario para proteger endpoints sensibles.

---

Actualizado: 2025-05
