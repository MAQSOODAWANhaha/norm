// builder/expression.go
package builder

import (
	"fmt"
	"strings"

	"norm/types"
)

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

// BuildAs 为现有表达式添加别名
func (e Expression) BuildAs(alias string) Expression {
	return Expression{Text: e.Text, Alias: alias}
}

// ExpressionBuilder 表达式构建器 (旧版，逐步废弃，保留用于向后兼容)
type ExpressionBuilder struct {
	parts []string
}

// NewExpression 创建新的表达式构建器
func NewExpression() *ExpressionBuilder {
	return &ExpressionBuilder{
		parts: make([]string, 0),
	}
}

// Property 添加属性表达式
func (eb *ExpressionBuilder) Property(property string) *ExpressionBuilder {
	eb.parts = append(eb.parts, property)
	return eb
}

// Equal 等于比较 (=)
func (eb *ExpressionBuilder) Equal(value interface{}) *ExpressionBuilder {
	eb.parts = append(eb.parts, "=", formatValue(value))
	return eb
}

// NotEqual 不等于比较 (<>)
func (eb *ExpressionBuilder) NotEqual(value interface{}) *ExpressionBuilder {
	eb.parts = append(eb.parts, "<>", formatValue(value))
	return eb
}

// LessThan 小于比较 (<)
func (eb *ExpressionBuilder) LessThan(value interface{}) *ExpressionBuilder {
	eb.parts = append(eb.parts, "<", formatValue(value))
	return eb
}

// LessThanOrEqual 小于等于比较 (<=)
func (eb *ExpressionBuilder) LessThanOrEqual(value interface{}) *ExpressionBuilder {
	eb.parts = append(eb.parts, "<=", formatValue(value))
	return eb
}

// GreaterThan 大于比较 (>)
func (eb *ExpressionBuilder) GreaterThan(value interface{}) *ExpressionBuilder {
	eb.parts = append(eb.parts, ">", formatValue(value))
	return eb
}

// GreaterThanOrEqual 大于等于比较 (>=)
func (eb *ExpressionBuilder) GreaterThanOrEqual(value interface{}) *ExpressionBuilder {
	eb.parts = append(eb.parts, ">=", formatValue(value))
	return eb
}

// Contains 包含操作 (CONTAINS)
func (eb *ExpressionBuilder) Contains(value string) *ExpressionBuilder {
	eb.parts = append(eb.parts, "CONTAINS", formatValue(value))
	return eb
}

// StartsWith 开始于操作 (STARTS WITH)
func (eb *ExpressionBuilder) StartsWith(value string) *ExpressionBuilder {
	eb.parts = append(eb.parts, "STARTS", "WITH", formatValue(value))
	return eb
}

// EndsWith 结束于操作 (ENDS WITH)
func (eb *ExpressionBuilder) EndsWith(value string) *ExpressionBuilder {
	eb.parts = append(eb.parts, "ENDS", "WITH", formatValue(value))
	return eb
}

// Regex 正则表达式匹配 (=~)
func (eb *ExpressionBuilder) Regex(pattern string) *ExpressionBuilder {
	eb.parts = append(eb.parts, "=~", formatValue(pattern))
	return eb
}

// In 在列表中 (IN)
func (eb *ExpressionBuilder) In(values ...interface{}) *ExpressionBuilder {
	var valueStrs []string
	for _, v := range values {
		valueStrs = append(valueStrs, formatValue(v))
	}
	eb.parts = append(eb.parts, "IN", "["+strings.Join(valueStrs, ", ")+"]")
	return eb
}

// IsNull 为空值 (IS NULL)
func (eb *ExpressionBuilder) IsNull() *ExpressionBuilder {
	eb.parts = append(eb.parts, "IS", "NULL")
	return eb
}

// IsNotNull 不为空值 (IS NOT NULL)
func (eb *ExpressionBuilder) IsNotNull() *ExpressionBuilder {
	eb.parts = append(eb.parts, "IS", "NOT", "NULL")
	return eb
}

// And 逻辑与 (AND)
func (eb *ExpressionBuilder) And(condition string) *ExpressionBuilder {
	eb.parts = append(eb.parts, "AND", condition)
	return eb
}

// Or 逻辑或 (OR)
func (eb *ExpressionBuilder) Or(condition string) *ExpressionBuilder {
	eb.parts = append(eb.parts, "OR", condition)
	return eb
}

// Not 逻辑非 (NOT)
func (eb *ExpressionBuilder) Not() *ExpressionBuilder {
	eb.parts = append([]string{"NOT"}, eb.parts...)
	return eb
}

// Build 构建最终的表达式字符串
func (eb *ExpressionBuilder) Build() string {
	return strings.Join(eb.parts, " ")
}

// BuildAs 构建带别名的表达式
func (eb *ExpressionBuilder) BuildAs(alias string) Expression {
	return Expression{
		Text:  eb.Build(),
		Alias: alias,
	}
}

// formatValue 格式化值 (旧版，逐步废弃)
func formatValue(value interface{}) string {
	switch v := value.(type) {
	case string:
		if strings.HasPrefix(v, "$") {
			return v // 参数引用
		}
		return fmt.Sprintf("'%s'", strings.ReplaceAll(v, "'", "'"))
	case int, int64, int32, int16, int8:
		return fmt.Sprintf("%v", v)
	case float64, float32:
		return fmt.Sprintf("%v", v)
	case bool:
		return fmt.Sprintf("%t", v)
	default:
		return fmt.Sprintf("'%v'", v)
	}
}

// =================================================================
// 新版谓词函数 (Predicate Functions) - 返回 types.Condition
// =================================================================

// Eq 等于表达式
func Eq(property string, value interface{}) types.Condition {
	return types.Predicate{Property: property, Operator: types.OpEqual, Value: value}
}

// Ne 不等于表达式
func Ne(property string, value interface{}) types.Condition {
	return types.Predicate{Property: property, Operator: types.OpNotEqual, Value: value}
}

// Lt 小于表达式
func Lt(property string, value interface{}) types.Condition {
	return types.Predicate{Property: property, Operator: types.OpLessThan, Value: value}
}

// Le 小于等于表达式
func Le(property string, value interface{}) types.Condition {
	return types.Predicate{Property: property, Operator: types.OpLessThanOrEqual, Value: value}
}

// Gt 大于表达式
func Gt(property string, value interface{}) types.Condition {
	return types.Predicate{Property: property, Operator: types.OpGreaterThan, Value: value}
}

// Ge 大于等于表达式
func Ge(property string, value interface{}) types.Condition {
	return types.Predicate{Property: property, Operator: types.OpGreaterThanOrEqual, Value: value}
}

// Contains 包含表达式
func Contains(property string, value interface{}) types.Condition {
	return types.Predicate{Property: property, Operator: types.OpContains, Value: value}
}

// StartsWith 开始于表达式
func StartsWith(property string, value interface{}) types.Condition {
	return types.Predicate{Property: property, Operator: types.OpStartsWith, Value: value}
}

// EndsWith 结束于表达式
func EndsWith(property string, value interface{}) types.Condition {
	return types.Predicate{Property: property, Operator: types.OpEndsWith, Value: value}
}

// Regex 正则表达式匹配
func Regex(property string, pattern string) types.Condition {
	return types.Predicate{Property: property, Operator: types.OpRegex, Value: pattern}
}

// In 在列表中表达式
func In(property string, values ...interface{}) types.Condition {
	return types.Predicate{Property: property, Operator: types.OpIn, Value: values}
}

// IsNull 为空表达式
func IsNull(property string) types.Condition {
	return types.Predicate{Property: property, Operator: types.OpIsNull}
}

// IsNotNull 不为空表达式
func IsNotNull(property string) types.Condition {
	return types.Predicate{Property: property, Operator: types.OpIsNotNull}
}

// Not 逻辑非
func Not(condition types.Condition) types.Condition {
	switch c := condition.(type) {
	case types.Predicate:
		c.Not = !c.Not // Toggle the Not flag
		return c
	case types.LogicalGroup:
		// For a group, it's more complex. A simple flag doesn't work well with Cypher syntax.
		// A better approach is to wrap it, but for now, we'll stick to negating predicates.
		// A full implementation might require a "NotGroup" type or similar.
		// For now, we return the group unmodified and log a warning or handle it in the builder.
		return c
	default:
		return condition
	}
}

// And 连接多个条件用 AND
func And(conditions ...types.Condition) types.Condition {
	return types.LogicalGroup{Operator: types.OpAnd, Conditions: conditions}
}

// Or 连接多个条件用 OR
func Or(conditions ...types.Condition) types.Condition {
	return types.LogicalGroup{Operator: types.OpOr, Conditions: conditions}
}

// ================================
// 聚合函数 (Aggregating Functions)
// ================================

// Count 计数函数
func Count(expression string) Expression {
	return Expression{Text: fmt.Sprintf("count(%s)", expression)}
}

// CountDistinct 去重计数函数
func CountDistinct(expression string) Expression {
	return Expression{Text: fmt.Sprintf("count(DISTINCT %s)", expression)}
}

// Sum 求和函数
func Sum(expression string) Expression {
	return Expression{Text: fmt.Sprintf("sum(%s)", expression)}
}

// Avg 平均值函数
func Avg(expression string) Expression {
	return Expression{Text: fmt.Sprintf("avg(%s)", expression)}
}

// Min 最小值函数
func Min(expression string) Expression {
	return Expression{Text: fmt.Sprintf("min(%s)", expression)}
}

// Max 最大值函数
func Max(expression string) Expression {
	return Expression{Text: fmt.Sprintf("max(%s)", expression)}
}

// Collect 收集函数
func Collect(expression string) Expression {
	return Expression{Text: fmt.Sprintf("collect(%s)", expression)}
}

// CollectDistinct 去重收集函数
func CollectDistinct(expression string) Expression {
	return Expression{Text: fmt.Sprintf("collect(DISTINCT %s)", expression)}
}

// ================================
// 字符串函数 (String Functions)
// ================================

// Lower 转小写函数
func Lower(expression string) Expression {
	return Expression{Text: fmt.Sprintf("lower(%s)", expression)}
}

// Upper 转大写函数
func Upper(expression string) Expression {
	return Expression{Text: fmt.Sprintf("upper(%s)", expression)}
}

// Trim 去除空格函数
func Trim(expression string) Expression {
	return Expression{Text: fmt.Sprintf("trim(%s)", expression)}
}

// LTrim 去除左侧空格函数
func LTrim(expression string) Expression {
	return Expression{Text: fmt.Sprintf("ltrim(%s)", expression)}
}

// RTrim 去除右侧空格函数
func RTrim(expression string) Expression {
	return Expression{Text: fmt.Sprintf("rtrim(%s)", expression)}
}

// Replace 替换字符串函数
func Replace(original, search, replace string) Expression {
	return Expression{Text: fmt.Sprintf("replace(%s, %s, %s)", original, search, replace)}
}

// Substring 子字符串函数
func Substring(str, start string, length ...string) Expression {
	if len(length) > 0 {
		return Expression{Text: fmt.Sprintf("substring(%s, %s, %s)", str, start, length[0])}
	}
	return Expression{Text: fmt.Sprintf("substring(%s, %s)", str, start)}
}

// Split 分割字符串函数
func Split(str, delimiter string) Expression {
	return Expression{Text: fmt.Sprintf("split(%s, %s)", str, delimiter)}
}

// ToString 转字符串函数
func ToString(expression string) Expression {
	return Expression{Text: fmt.Sprintf("toString(%s)", expression)}
}

// Left 左侧字符串函数
func Left(str, length string) Expression {
	return Expression{Text: fmt.Sprintf("left(%s, %s)", str, length)}
}

// Right 右侧字符串函数
func Right(str, length string) Expression {
	return Expression{Text: fmt.Sprintf("right(%s, %s)", str, length)}
}

// Reverse 反转字符串函数
func Reverse(str string) Expression {
	return Expression{Text: fmt.Sprintf("reverse(%s)", str)}
}

// ================================
// 数学函数 (Mathematical Functions)
// ================================

// Abs 绝对值函数
func Abs(expression string) Expression {
	return Expression{Text: fmt.Sprintf("abs(%s)", expression)}
}

// Ceil 向上取整函数
func Ceil(expression string) Expression {
	return Expression{Text: fmt.Sprintf("ceil(%s)", expression)}
}

// Floor 向下取整函数
func Floor(expression string) Expression {
	return Expression{Text: fmt.Sprintf("floor(%s)", expression)}
}

// Round 四舍五入函数
func Round(expression string, precision ...string) Expression {
	if len(precision) > 0 {
		return Expression{Text: fmt.Sprintf("round(%s, %s)", expression, precision[0])}
	}
	return Expression{Text: fmt.Sprintf("round(%s)", expression)}
}

// Sign 符号函数
func Sign(expression string) Expression {
	return Expression{Text: fmt.Sprintf("sign(%s)", expression)}
}

// Sqrt 平方根函数
func Sqrt(expression string) Expression {
	return Expression{Text: fmt.Sprintf("sqrt(%s)", expression)}
}

// Exp 指数函数
func Exp(expression string) Expression {
	return Expression{Text: fmt.Sprintf("exp(%s)", expression)}
}

// Log 自然对数函数
func Log(expression string) Expression {
	return Expression{Text: fmt.Sprintf("log(%s)", expression)}
}

// Log10 十进制对数函数
func Log10(expression string) Expression {
	return Expression{Text: fmt.Sprintf("log10(%s)", expression)}
}

// Sin 正弦函数
func Sin(expression string) Expression {
	return Expression{Text: fmt.Sprintf("sin(%s)", expression)}
}

// Cos 余弦函数
func Cos(expression string) Expression {
	return Expression{Text: fmt.Sprintf("cos(%s)", expression)}
}

// Tan 正切函数
func Tan(expression string) Expression {
	return Expression{Text: fmt.Sprintf("tan(%s)", expression)}
}

// Rand 随机数函数
func Rand() Expression {
	return Expression{Text: "rand()"}
}

// ================================
// 列表函数 (List Functions)
// ================================

// Size 大小函数 (适用于列表、字符串、路径)
func Size(expression string) Expression {
	return Expression{Text: fmt.Sprintf("size(%s)", expression)}
}

// Head 获取列表第一个元素函数
func Head(list string) Expression {
	return Expression{Text: fmt.Sprintf("head(%s)", list)}
}

// Last 获取列表最后一个元素函数
func Last(list string) Expression {
	return Expression{Text: fmt.Sprintf("last(%s)", list)}
}

// Tail 获取除第一个元素外的列表函数
func Tail(list string) Expression {
	return Expression{Text: fmt.Sprintf("tail(%s)", list)}
}

// Range 范围函数
func Range(start, end string, step ...string) Expression {
	if len(step) > 0 {
		return Expression{Text: fmt.Sprintf("range(%s, %s, %s)", start, end, step[0])}
	}
	return Expression{Text: fmt.Sprintf("range(%s, %s)", start, end)}
}

// Keys 获取属性键函数
func Keys(expression string) Expression {
	return Expression{Text: fmt.Sprintf("keys(%s)", expression)}
}

// Labels 获取节点标签函数
func Labels(node string) Expression {
	return Expression{Text: fmt.Sprintf("labels(%s)", node)}
}

// Type 获取关系类型函数
func Type(relationship string) Expression {
	return Expression{Text: fmt.Sprintf("type(%s)", relationship)}
}

// ================================
// 谓词函数 (Predicate Functions)
// ================================

// Exists 存在性检查函数
func Exists(expression string) Expression {
	return Expression{Text: fmt.Sprintf("exists(%s)", expression)}
}

// IsEmpty 空值检查函数
func IsEmpty(expression string) Expression {
	return Expression{Text: fmt.Sprintf("isEmpty(%s)", expression)}
}

// All 全部满足条件函数
func All(variable, list, predicate string) Expression {
	return Expression{Text: fmt.Sprintf("all(%s IN %s WHERE %s)", variable, list, predicate)}
}

// Any 任意满足条件函数
func Any(variable, list, predicate string) Expression {
	return Expression{Text: fmt.Sprintf("any(%s IN %s WHERE %s)", variable, list, predicate)}
}

// None 全部不满足条件函数
func None(variable, list, predicate string) Expression {
	return Expression{Text: fmt.Sprintf("none(%s IN %s WHERE %s)", variable, list, predicate)}
}

// Single 仅一个满足条件函数
func Single(variable, list, predicate string) Expression {
	return Expression{Text: fmt.Sprintf("single(%s IN %s WHERE %s)", variable, list, predicate)}
}

// ================================
// 标量函数 (Scalar Functions)
// ================================

// Coalesce 合并函数 (返回第一个非空值)
func Coalesce(expressions ...string) Expression {
	return Expression{Text: fmt.Sprintf("coalesce(%s)", strings.Join(expressions, ", "))}
}

// ElementId 获取元素ID函数
func ElementId(element string) Expression {
	return Expression{Text: fmt.Sprintf("elementId(%s)", element)}
}

// Id 获取ID函数 (已弃用，但仍然支持)
func Id(element string) Expression {
	return Expression{Text: fmt.Sprintf("id(%s)", element)}
}

// Properties 获取属性函数
func Properties(element string) Expression {
	return Expression{Text: fmt.Sprintf("properties(%s)", element)}
}

// StartNode 获取关系起始节点函数
func StartNode(relationship string) Expression {
	return Expression{Text: fmt.Sprintf("startNode(%s)", relationship)}
}

// EndNode 获取关系结束节点函数
func EndNode(relationship string) Expression {
	return Expression{Text: fmt.Sprintf("endNode(%s)", relationship)}
}

// ================================
// 时间函数 (Temporal Functions)
// ================================

// Date 日期函数
func Date(expression ...string) Expression {
	if len(expression) > 0 {
		return Expression{Text: fmt.Sprintf("date(%s)", expression[0])}
	}
	return Expression{Text: "date()"}
}

// DateTime 日期时间函数
func DateTime(expression ...string) Expression {
	if len(expression) > 0 {
		return Expression{Text: fmt.Sprintf("datetime(%s)", expression[0])}
	}
	return Expression{Text: "datetime()"}
}

// Time 时间函数
func Time(expression ...string) Expression {
	if len(expression) > 0 {
		return Expression{Text: fmt.Sprintf("time(%s)", expression[0])}
	}
	return Expression{Text: "time()"}
}

// LocalTime 本地时间函数
func LocalTime(expression ...string) Expression {
	if len(expression) > 0 {
		return Expression{Text: fmt.Sprintf("localtime(%s)", expression[0])}
	}
	return Expression{Text: "localtime()"}
}

// LocalDateTime 本地日期时间函数
func LocalDateTime(expression ...string) Expression {
	if len(expression) > 0 {
		return Expression{Text: fmt.Sprintf("localdatetime(%s)", expression[0])}
	}
	return Expression{Text: "localdatetime()"}
}

// Duration 持续时间函数
func Duration(expression string) Expression {
	return Expression{Text: fmt.Sprintf("duration(%s)", expression)}
}

// ================================
// 路径函数 (Path Functions)
// ================================

// Length 路径长度函数
func Length(path string) Expression {
	return Expression{Text: fmt.Sprintf("length(%s)", path)}
}

// Nodes 获取路径中所有节点函数
func Nodes(path string) Expression {
	return Expression{Text: fmt.Sprintf("nodes(%s)", path)}
}

// Relationships 获取路径中所有关系函数
func Relationships(path string) Expression {
	return Expression{Text: fmt.Sprintf("relationships(%s)", path)}
}

// ShortestPath 最短路径函数
func ShortestPath(pattern string) Expression {
	return Expression{Text: fmt.Sprintf("shortestPath(%s)", pattern)}
}

// AllShortestPaths 所有最短路径函数
func AllShortestPaths(pattern string) Expression {
	return Expression{Text: fmt.Sprintf("allShortestPaths(%s)", pattern)}
}

// ================================
// 辅助函数 (Helper Functions)
// ================================

// Case 条件表达式构建器
type CaseBuilder struct {
	parts []string
}

// NewCase 创建新的 CASE 表达式构建器
func NewCase() *CaseBuilder {
	return &CaseBuilder{
		parts: []string{"CASE"},
	}
}

// When 添加 WHEN 条件
func (cb *CaseBuilder) When(condition, result string) *CaseBuilder {
	cb.parts = append(cb.parts, fmt.Sprintf("WHEN %s THEN %s", condition, result))
	return cb
}

// Else 添加 ELSE 子句
func (cb *CaseBuilder) Else(result string) *CaseBuilder {
	cb.parts = append(cb.parts, fmt.Sprintf("ELSE %s", result))
	return cb
}

// End 结束 CASE 表达式
func (cb *CaseBuilder) End() Expression {
	cb.parts = append(cb.parts, "END")
	return Expression{Text: strings.Join(cb.parts, " ")}
}

// ================================
// 比较和逻辑运算符增强
// ================================

// Xor 异或操作
func Xor(left, right string) string {
	return fmt.Sprintf("(%s) XOR (%s)", left, right)
}

// DistinctValues 去重表达式
func DistinctValues(expression string) Expression {
	return Expression{Text: fmt.Sprintf("DISTINCT %s", expression)}
}
