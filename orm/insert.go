package orm

import (
	"context"
	"github.com/go-tour/orm/internal/errs"
)

type UpsertBuilder[T any] struct {
	i               *Inserter[T]
	conflictColumns []string
}

type Upsert struct {
	conflictColumns []string
	assigns         []Assignable
}

func (o *UpsertBuilder[T]) ConflictColumns(cols ...string) *UpsertBuilder[T] {
	o.conflictColumns = cols
	return o
}

// Update 也可以看做是一个终结方法，重新回到 Inserter 里面
func (o *UpsertBuilder[T]) Update(assigns ...Assignable) *Inserter[T] {
	o.i.upsert = &Upsert{
		conflictColumns: o.conflictColumns,
		assigns:         assigns,
	}
	return o.i
}

type Inserter[T any] struct {
	builder
	values  []*T
	columns []string
	upsert  *Upsert
	sess    Session
}

func NewInserter[T any](sess Session) *Inserter[T] {
	c := sess.getCore()
	return &Inserter[T]{
		sess: sess,
		builder: builder{
			core:    c,
			dialect: c.dialect,
			quoter:  c.dialect.quoter(),
		},
	}
}

func (i *Inserter[T]) Values(vals ...*T) *Inserter[T] {
	i.values = vals
	return i
}

func (i *Inserter[T]) OnDuplicateKey() *UpsertBuilder[T] {
	return &UpsertBuilder[T]{
		i: i,
	}
}

// Columns 指定要插入的列
// TODO 目前我们只支持指定具体的列，但是不支持复杂的表达式
// 例如不支持 VALUES(..., now(), now()) 这种在 MySQL 里面常用的
func (i *Inserter[T]) Columns(cols ...string) *Inserter[T] {
	i.columns = cols
	return i
}

func (i *Inserter[T]) Build() (*Query, error) {
	if len(i.values) == 0 {
		return nil, errs.ErrInsertZeroRow
	}
	if i.model == nil {
		m, err := i.r.Get(i.values[0])
		if err != nil {
			return nil, err
		}
		i.model = m
	}

	i.sb.WriteString("INSERT INTO ")
	i.quote(i.model.TableName)
	i.sb.WriteString("(")

	fields := i.model.Fields
	if len(i.columns) != 0 {
		fields = make([]*model.Field, 0, len(i.columns))
		for _, c := range i.columns {
			field, ok := i.model.FieldMap[c]
			if !ok {
				return nil, errs.NewErrUnknownField(c)
			}
			fields = append(fields, field)
		}
	}

	// (len(i.values) + 1) 中 +1 是考虑到 UPSERT 语句会传递额外的参数
	i.args = make([]any, 0, len(fields)*(len(i.values)+1))
	for idx, fd := range fields {
		if idx > 0 {
			i.sb.WriteByte(',')
		}
		i.quote(fd.ColName)
	}

	i.sb.WriteString(") VALUES")
	for vIdx, val := range i.values {
		if vIdx > 0 {
			i.sb.WriteByte(',')
		}
		refVal := i.valCreator(val, i.model)
		i.sb.WriteByte('(')
		for fIdx, field := range fields {
			if fIdx > 0 {
				i.sb.WriteByte(',')
			}
			i.sb.WriteByte('?')
			fdVal, err := refVal.Field(field.GoName)
			if err != nil {
				return nil, err
			}
			i.addArgs(fdVal)
		}
		i.sb.WriteByte(')')
	}

	if i.upsert != nil {
		err := i.core.dialect.buildUpsert(&i.builder, i.upsert)
		if err != nil {
			return nil, err
		}
	}

	i.sb.WriteString(";")
	return &Query{
		SQL:  i.sb.String(),
		Args: i.args,
	}, nil
}

func (i *Inserter[T]) Exec(ctx context.Context) Result {
	if i.model == nil {
		m, err := i.r.Get(new(T))
		if err != nil {
			return Result{
				err: err,
			}
		}
		i.model = m
	}
	return exec(ctx, i.sess, i.core, &QueryContext{
		builder: i,
		Type:    "INSERT",
		Model:   i.model,
	})
}
