package model

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/go-tour/orm/demo4/internal/errs"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_parseModel(t *testing.T) {
	tests := []struct {
		name    string
		input   any
		want    *Model
		wantErr error
	}{
		{
			name:  "ptr",
			input: &TestModel{},
			want: &Model{
				TableName: "test_model",
				FieldMap: map[string]*Field{
					"Id": {
						ColName: "id",
					},
					"FirstName": {
						ColName: "first_name",
					},
					"Age": {
						ColName: "age",
					},
					"LastName": {
						ColName: "last_name",
					},
				},
			},
		},
		{
			name:    "struct",
			input:   TestModel{},
			wantErr: errs.ErrPointerOnly,
		},
		{
			name:    "map",
			input:   map[string]string{},
			wantErr: errs.ErrPointerOnly,
		},
		{
			name:    "nil",
			input:   nil,
			wantErr: errors.New("orm: 不支持 nil"),
		},
		{
			name:  "nil with type",
			input: (*TestModel)(nil),
			want: &Model{
				TableName: "test_model",
				FieldMap: map[string]*Field{
					"Id": {
						ColName: "id",
					},
					"FirstName": {
						ColName: "first_name",
					},
					"Age": {
						ColName: "age",
					},
					"LastName": {
						ColName: "last_name",
					},
				},
			},
		},

		{
			name: "column tag",
			input: func() any {
				// 我们把测试结构体定义在方法内部，防止被其它用例访问
				type ColumnTag struct {
					ID uint64 `orm:"column=id"`
				}
				return &ColumnTag{}
			}(),
			want: &Model{
				TableName: "column_tag",
				FieldMap: map[string]*Field{
					// 默认是 i_d
					"ID": {
						ColName: "id",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &registry{}
			m, err := r.Register(tt.input)
			assert.Equal(t, tt.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tt.want, m)
		})
	}
}

func TestSwitch(t *testing.T) {
	Switch(nil)
	Switch((*TestModel)(nil))
}

func Switch(val any) {
	switch v := val.(type) {
	case nil:
		fmt.Println("hello, nil")
	case *TestModel:
		fmt.Printf("hello, test Model %v \n", v)
		if v == nil {
			fmt.Printf("hello, test Model nil %v \n", v)
		}
	}
}

type TestModel struct {
	Id        int64
	FirstName string
	Age       int8
	LastName  *sql.NullString
}
