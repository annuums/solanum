# Example for Dependency Injection

이 디렉토리에는 Solanum 웹 프레임워크에 PostgreSQL 의존성 주입을 적용한 두 가지 버전의 예제가 있습니다.

해당 예제는 PostgreSQL 17 버전과 함께 동작합니다.

- **MVC 버전**: `docs/examples/database/mvc`

```bash
docs/examples/database
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
