package demo

import (
	"context"
	"database/sql"
)

// Querier SELECT 语句
type Querier[T any] interface {
	// user := xxx.Get(ctx)

	// 不再需要写成
	// var user User
	// Get(ctx, &user)
	Get(ctx context.Context) (*T, error)
	GetMulti(ctx context.Context) ([]*T, error)
}

// Executor UPDATE, DELETE, INSERT
type Executor interface {
	Exec(ctx context.Context) (sql.Result, error)
}

// QueryBuilder 作为构建SQL这一个单独步骤的顶级抽象
type QueryBuilder interface {
	Build() (*Query, error)
}

// db.Exec
// db.QueryRow
// db.Query
type Query struct {
	SQL  string
	Args []any
}
