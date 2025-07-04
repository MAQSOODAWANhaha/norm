# Cypher ORM 简化架构设计文档

## 概述

本文档概述了一个轻量级的 Cypher 语言 ORM（对象关系映射）包的架构设计，专注于 Go struct 到 Cypher 语句的转换和基本查询构建功能。该 ORM 将提供一个简洁的 Go 原生接口，用于构造和验证 Cypher 查询。

## 核心设计原则

1. **简洁性**: 保持最小化的依赖和简洁的 API
2. **类型安全**: 利用 Go 的类型系统提供编译时验证
3. **流畅接口**: 提供可链式调用的查询构建接口
4. **专注核心**: 只实现 struct 到 Cypher 转换和查询验证
5. **易于扩展**: 为未来功能预留扩展空间

## 简化架构图

```
┌─────────────────────────────────────────────────────────────────┐
│                         客户端应用程序                          │
└─────────────────────────┬───────────────────────────────────────┘
                          │
┌─────────────────────────▼───────────────────────────────────────┐
│                      Cypher ORM 包                             │
│  ┌─────────────────┐ ┌─────────────────┐ ┌─────────────────┐   │
│  │    查询构建器    │ │    模型系统     │ │   类型转换器    │   │
│  └─────────────────┘ └─────────────────┘ └─────────────────┘   │
│  ┌─────────────────┐ ┌─────────────────┐                     │
│  │    验证器       │ │    解析器       │                     │
│  └─────────────────┘ └─────────────────┘                     │
└─────────────────────────┬───────────────────────────────────────┘
                          │
┌─────────────────────────▼───────────────────────────────────────┐
│                    生成的 Cypher 查询                          │
│                  + 参数化查询参数                               │
└─────────────────────────────────────────────────────────────────┘
```

## 组件架构

### 1. 查询构建器 (Builder)
**目标**: 专注于第一和第二阶段功能
- **QueryBuilder**: 主查询构建器，支持流畅接口
- **NodeBuilder**: 节点模式构建
- **RelationshipBuilder**: 关系模式构建
- **ClauseBuilder**: Cypher 子句构建

### 2. 实体解析系统 (Entity Parsing)
**目标**: 直接从 Go struct 实例提取标签和属性信息
- **EntityParser**: 实时解析结构体标签和属性
- **LabelExtractor**: 从标签中提取节点标签
- **PropertyExtractor**: 从字段值中提取属性映射

### 3. 类型系统 (Types)
**目标**: Go 类型到 Cypher 类型的转换
- **Converter**: 类型转换器
- **Validator**: 类型验证
- **Registry**: 类型注册表

### 4. 验证系统 (Validator)
**目标**: 第二阶段的查询验证功能
- **QueryValidator**: 查询语法验证
- **StructureValidator**: 查询结构验证
- **ParameterValidator**: 参数验证

### 5. 解析系统 (Parser)
**目标**: 第二阶段的查询解析功能
- **CypherParser**: Cypher 语法解析
- **PatternParser**: 模式解析
- **ExpressionParser**: 表达式解析

## 简化包结构

```
norm/
├── builder/              # 查询构建器
│   ├── query.go         # 主查询构建器
│   ├── node.go          # 节点构建器
│   ├── relationship.go  # 关系构建器
│   ├── expression.go    # 表达式构建器
│   ├── entity.go        # 实体解析器
│   └── types.go         # 构建器类型定义
├── types/               # 类型系统
│   ├── core.go          # 核心类型定义
│   ├── converter.go     # 类型转换器
│   └── validator.go     # 类型验证器
├── validator/           # 验证系统 (第二阶段)
│   ├── query.go         # 查询验证
│   ├── structure.go     # 结构验证
│   └── parameter.go     # 参数验证
├── parser/              # 解析系统 (第二阶段)
│   ├── cypher.go        # Cypher 解析
│   ├── pattern.go       # 模式解析
│   └── expression.go    # 表达式解析
├── examples/            # 使用示例
│   └── main.go
├── tests/               # 测试文件
└── docs/                # 文档
```

## 核心接口设计

### 查询构建器接口
```go
type QueryBuilder interface {
    // 基本子句
    Match(pattern string) QueryBuilder
    OptionalMatch(pattern string) QueryBuilder
    Create(pattern string) QueryBuilder
    Merge(pattern string) QueryBuilder
    Where(condition string) QueryBuilder
    Return(expressions ...interface{}) QueryBuilder
    With(expressions ...interface{}) QueryBuilder
    
    // 排序和限制
    OrderBy(fields ...string) QueryBuilder
    Skip(count int) QueryBuilder
    Limit(count int) QueryBuilder
    
    // 参数操作
    SetParameter(key string, value interface{}) QueryBuilder
    
    // 构建操作
    Build() (QueryResult, error)
    Validate() []ValidationError
    
    // 实体操作 (无需预注册)
    MatchEntity(entity interface{}) QueryBuilder
    CreateEntity(entity interface{}) QueryBuilder
    MergeEntity(entity interface{}) QueryBuilder
}

// 表达式别名辅助函数
func As(expression, alias string) Expression
```

### 实体解析接口
```go
type EntityInfo struct {
    Labels     []string
    Properties map[string]interface{}
}

// 直接解析结构体实例，无需预注册
func ParseEntity(entity interface{}) (*EntityInfo, error)
    Validate(entity interface{}) error
}
```

### 查询结果接口
```go
type QueryResult struct {
    Query      string                 // 生成的 Cypher 查询
    Parameters map[string]interface{} // 查询参数
    Valid      bool                   // 查询是否有效 (第二阶段)
    Errors     []ValidationError      // 验证错误 (第二阶段)
}
```

## 实现阶段

### 第一阶段：Struct 到 Cypher 转换 (当前目标)
**时间**: 2-3 周
**目标**: 实现基本的 Go 结构体到 Cypher 语句转换

**包含功能**:
1. **核心类型定义** - 定义所有基础数据结构
2. **实体注册表** - struct 标签解析和元数据管理
3. **查询构建器** - 基本的查询构建功能
4. **节点和关系构建器** - 模式构建工具
5. **类型转换器** - Go 类型到 Cypher 类型转换
6. **基础示例** - 展示核心功能的使用方法

**输出**: 可以将 Go struct 转换为有效的 Cypher 查询语句

### 第二阶段：查询解析和验证 (当前目标)
**时间**: 2-3 周
**目标**: 增强查询构建功能，添加解析和验证

**包含功能**:
1. **查询验证器** - 验证生成的 Cypher 语法正确性
2. **结构验证器** - 验证查询结构的合理性
3. **参数验证器** - 验证查询参数的有效性
4. **Cypher 解析器** - 解析和分析 Cypher 语句
5. **模式解析器** - 解析节点和关系模式
6. **错误处理** - 详细的错误信息和建议
7. **高级示例** - 展示验证和解析功能

**输出**: 不仅能生成 Cypher 查询，还能验证其正确性

## 错误处理策略

1. **结构化错误**: 使用类型化错误区分不同错误类别
2. **详细错误信息**: 提供具体的错误位置和修复建议
3. **渐进式验证**: 从语法到语义的多层验证
4. **用户友好**: 清晰的错误消息和文档引用

## 测试策略

1. **单元测试**: 测试每个组件的独立功能
2. **集成测试**: 测试组件间的协作
3. **示例测试**: 确保示例代码正确运行
4. **边界测试**: 测试边界条件和错误情况

## 性能考虑

1. **内存效率**: 合理的对象生命周期管理
2. **构建性能**: 优化查询构建过程
3. **验证性能**: 高效的验证算法
4. **缓存策略**: 缓存元数据和验证结果

## 扩展性设计

1. **插件接口**: 为类型转换器和验证器提供插件机制
2. **中间件模式**: 支持查询构建中间件
3. **钩子函数**: 在关键点提供钩子函数
4. **配置选项**: 提供丰富的配置选项

这个简化的架构专注于核心功能，去除了数据库连接、事务管理等复杂特性，使得项目更加轻量和易于维护。第一和第二阶段的功能足以提供强大的 Cypher 查询构建和验证能力。