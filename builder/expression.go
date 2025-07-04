// builder/expression.go
package builder

import "fmt"

// Expression 代表一个可被别名的表达式
type Expression struct {
	Text  string
	Alias string
}

// As 创建一个带别名的表达式
func As(expression, alias string) Expression {
	return Expression{Text: expression, Alias: alias}
}

// String 实现 Stringer 接口
func (e Expression) String() string {
	if e.Alias != "" {
		return fmt.Sprintf("%s AS %s", e.Text, e.Alias)
	}
	return e.Text
}
