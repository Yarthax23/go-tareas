# 📘 Guía Rápida de PostgreSQL + Go

## 🧭 Conexión desde terminal

```bash
psql -U <usuario> -d <base_de_datos> -h localhost
# Ejemplo:
psql -U golang -d todos -h localhost
```

---

## 📦 Bases de Datos

| Acción               | Comando                                     |
|----------------------|---------------------------------------------|
| Ver todas las bases  | `\l`                                       |
| Crear nueva base     | `CREATE DATABASE nombre;`                   |
| Cambiar de base      | `\c nombre_base`                            |
| Borrar base          | `DROP DATABASE nombre;` *(superuser)*       |

---

## 📋 Tablas

| Acción                   | Comando SQL                                |
|--------------------------|--------------------------------------------|
| Ver tablas               | `\dt`                                      |
| Ver estructura de tabla  | `\d nombre_tabla`                          |
| Crear tabla              | `CREATE TABLE nombre (...);`               |
| Borrar tabla             | `DROP TABLE nombre;`                       |
| Reiniciar tabla          | `TRUNCATE tareas RESTART IDENTITY;`        |

| Agregar columna          | `ALTER TABLE tareas ADD COLUMN autor TEXT;` | 
| Alterar columna          | `ALTER TABLE tareas ALTER COLUMN listo SET DEFAULT 'NO';`|
| Quitar default           | `ALTER TABLE tareas ALTER COLUMN listo DROP DEFAULT;`|
| Renombar columna         | `ALTER TABLE tareas RENAME COLUMN listo TO nueva_columna;`|

| Popular columna B        | `UPDATE tareas`
| desde una columna A      | `SET resuelto_bool = CASE`
                           | `WHEN resuelto = 'SI' THEN TRUE`
                           | `ELSE FALSE`
                           | `END;`




## Columnas

| Acción                              | Afecta datos viejos | Afecta datos nuevos    |
| ----------------------------------- | ------------------- | ---------------------- |
| `ADD COLUMN autor TEXT`             | `NULL` en viejos    | Debés especificar      |
| `ADD COLUMN autor TEXT DEFAULT 'X'` | `NULL` en viejos    | `'X'` si no ponés nada |
| `UPDATE tareas SET autor = 'X' ...`  | Actualiza viejos    | No afecta nuevos       |
|            `...WHERE autor IS NULL;`|

---

## 📄 Datos

| Acción                    | Comando SQL                                         |
|---------------------------|-----------------------------------------------------|
| Ver todos los datos       | `SELECT * FROM tareas;`                              |
| Insertar un dato          | `INSERT INTO tareas (contenido) VALUES ('Hola');`   |
| Editar un dato            | `UPDATE tareas SET contenido = 'Nuevo' WHERE id = 1;` |
| Borrar un dato específico | `DELETE FROM tareas WHERE id = 1;`                  |
| Borrar todos los datos    | `DELETE FROM tareas;`                                |

---

## 👤 Usuarios y Permisos

| Acción                      | Comando SQL                                             |
|-----------------------------|----------------------------------------------------------|
| Crear usuario               | `CREATE USER golang WITH PASSWORD 'clave123';`           |
| Dar permisos sobre base     | `GRANT ALL PRIVILEGES ON DATABASE todos TO golang;`      |
| Dar permisos sobre tablas   | `GRANT ALL ON TABLE tareas TO golang;`                   |

---

## 🧑‍💻 Desde Go

### Crear tabla:
```go
db.Exec(`CREATE TABLE IF NOT EXISTS tareas (
    id SERIAL PRIMARY KEY,
    contenido TEXT NOT NULL
)`)
```

### Insertar dato:
```go
db.Exec("INSERT INTO tareas (contenido) VALUES ($1)", contenido)
```

### Leer datos:
```go
rows, _ := db.Query("SELECT id, contenido FROM tareas")
```

### Borrar por ID:
```go
db.Exec("DELETE FROM tareas WHERE id = $1", id)
```

---

## 🧽 Para limpiar todo

### Borrar todos los datos:
```sql
DELETE FROM tareas;
```

### Borrar la tabla:
```sql
DROP TABLE tareas;
```

### Borrar base de datos y usuario (como superusuario):
```bash
dropdb todos
dropuser golang
```

---

## 📎 Comandos útiles de psql

| Acción                    | Comando |
|---------------------------|---------|
| Ver comandos disponibles  | `\?`    |
| Ver ayuda de comandos SQL | `\h`    |
| Salir de psql             | `\q`    |

---

## 📚 Ejemplos en Go

### 🔌 Conectar a PostgreSQL

```go
import (
    "database/sql"
    "fmt"
    _ "github.com/lib/pq"
)

func main() {
    connStr := "host=localhost user=golang dbname=todos sslmode=disable"
    db, err := sql.Open("postgres", connStr)
    if err != nil {
        panic(err)
    }
    defer db.Close()

    fmt.Println("Conectado a PostgreSQL")
}
```

---

### 🛠 Crear tabla (si no existe)

```go
_, err = db.Exec(`
    CREATE TABLE IF NOT EXISTS tareas (
        id SERIAL PRIMARY KEY,
        contenido TEXT NOT NULL
    );
`)
if err != nil {
    panic(err)
}
```

---

### 💾 Insertar datos

```go
contenido := "Hola desde Go"
_, err = db.Exec("INSERT INTO tareas (contenido) VALUES ($1)", contenido)
if err != nil {
    panic(err)
}
```

---

### 📖 Ver todos los datos

```go
rows, err := db.Query("SELECT id, contenido FROM tareas")
if err != nil {
    panic(err)
}
defer rows.Close()

for rows.Next() {
    var id int
    var contenido string
    rows.Scan(&id, &contenido)
    fmt.Printf("ID: %d - Contenido: %s\n", id, contenido)
}
```

---

---

## ❓ ¿Qué es `err` y para qué sirve en Go?

### 📌 ¿Qué es `err`?

En Go, muchas funciones devuelven un **valor de error** además del resultado. Este valor se guarda en una variable llamada `err`. Si todo sale bien, `err` será `nil`. Si algo falla, `err` contendrá información sobre el error.

### 📋 Ejemplo simple

```go
resultado, err := hacerAlgo()
if err != nil {
    panic(err)
}
```

---

### ⚠️ ¿Cuándo puede fallar algo?

| Función de Go     | Qué intenta hacer             | Cuándo puede fallar                                 |
|-------------------|-------------------------------|-----------------------------------------------------|
| `sql.Open(...)`   | Conectar a base de datos      | Credenciales mal escritas, base no existe           |
| `db.Exec(...)`    | Crear tabla, insertar datos   | Error en SQL, tabla ya existe, sin permisos         |
| `rows.Next()`     | Leer siguiente resultado      | Error de conexión o datos corruptos                 |

---

### 🧨 ¿Qué hace `panic(err)`?

Detiene el programa y muestra el error. Es útil para detectar fallos críticos mientras estás desarrollando.

---

### ✅ Buenas prácticas (más adelante)

Más adelante, en lugar de `panic`, es mejor:

```go
if err != nil {
    log.Fatalf("Error: %v", err)
}
```

Esto da un mensaje útil y termina el programa de forma más controlada.

---
