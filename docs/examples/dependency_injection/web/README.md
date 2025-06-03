# Example for Dependency Injection

- This directory contains two versions of an example applying PostgreSQL dependency injection in the Solanum web framework.
- The examples are designed to work with PostgreSQL version 17.

```bash
docs/examples/dependency_injection/web
├── mvc
└── hexagonal
```

## User 테이블 DDL
```sql
CREATE TABLE users (
  id SERIAL PRIMARY KEY,
  name TEXT NOT NULL,
  email TEXT UNIQUE NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
```
