// tests/entity_functionality_test.go
package tests

import (
	"testing"

	"norm/builder"
	"norm/types"
	"github.com/stretchr/testify/assert"
)

// TestUser struct for testing purposes
type TestUser struct {
	_        struct{} `cypher:"label:TestUser,Active"`
	Username string   `cypher:"username"`
	Email    string   `cypher:"email"`
	Age      int      `cypher:"age"`
}

// TestDefaultLabelUser struct for testing default label generation
type TestDefaultLabelUser struct {
	Username string `cypher:"username"`
}

func TestParseEntityWithMultipleLabels(t *testing.T) {
	user := &TestUser{Username: "test", Email: "test@example.com"}
	entityInfo, err := builder.ParseEntity(user)

	assert.NoError(t, err)
	assert.NotNil(t, entityInfo)
	assert.Contains(t, entityInfo.Labels.ToStrings(), "TestUser")
	assert.Contains(t, entityInfo.Labels.ToStrings(), "Active")
	assert.Equal(t, 2, len(entityInfo.Labels))
}

func TestParseEntityWithDefaultLabel(t *testing.T) {
	user := &TestDefaultLabelUser{Username: "test"}
	entityInfo, err := builder.ParseEntity(user)

	assert.NoError(t, err)
	assert.NotNil(t, entityInfo)
	assert.Contains(t, entityInfo.Labels.ToStrings(), "TestDefaultLabelUser")
	assert.Equal(t, 1, len(entityInfo.Labels))
}

func TestSetEntity(t *testing.T) {
	qb := builder.NewQueryBuilder()
	user := &User{Username: "test", Email: "test@example.com"}

	res, err := qb.Match(user).As("u").
		SetEntity(user, "u").
		Build()

	assert.NoError(t, err)
	assert.Contains(t, res.Query, "MATCH (u:User:Person)")
	assert.Contains(t, res.Query, "SET")
	assert.Contains(t, res.Query, "u.username = ")
	assert.Contains(t, res.Query, "u.email = ")

	foundUsername := false
	foundEmail := false
	for _, v := range res.Parameters {
		if v == "test" {
			foundUsername = true
		}
		if v == "test@example.com" {
			foundEmail = true
		}
	}
	assert.True(t, foundUsername, "username parameter not found")
	assert.True(t, foundEmail, "email parameter not found")
}

func TestRemoveProperties(t *testing.T) {
	qb := builder.NewQueryBuilder()
	user := &User{}

	res, err := qb.Match(user).As("u").
		RemoveProperties(user, "u", "email").
		Build()

	assert.NoError(t, err)
	assert.Contains(t, res.Query, "MATCH (u:User:Person)")
	assert.Contains(t, res.Query, "REMOVE u.email")
}

func TestRemoveAllProperties(t *testing.T) {
	qb := builder.NewQueryBuilder()
	user := &User{}

	res, err := qb.Match(user).As("u").
		RemoveProperties(user, "u").
		Build()

	assert.NoError(t, err)
	assert.Contains(t, res.Query, "MATCH (u:User:Person)")
	assert.Contains(t, res.Query, "REMOVE u.username, u.email, u.age")
}

func TestDeleteEntity(t *testing.T) {
	qb := builder.NewQueryBuilder()
	user := &User{}

	res, err := qb.Match(user).As("u").
		Delete(types.Entity{Struct: user, Alias: "u"}).
		Build()

	assert.NoError(t, err)
	assert.Contains(t, res.Query, "MATCH (u:User:Person)")
	assert.Contains(t, res.Query, "DELETE u")
}

func TestDetachDeleteEntity(t *testing.T) {
	qb := builder.NewQueryBuilder()
	user := &User{}

	res, err := qb.Match(user).As("u").
		DetachDelete(types.Entity{Struct: user, Alias: "u"}).
		Build()

	assert.NoError(t, err)
	assert.Contains(t, res.Query, "MATCH (u:User:Person)")
	assert.Contains(t, res.Query, "DETACH DELETE u")
}

func TestReturnEntity(t *testing.T) {
	qb := builder.NewQueryBuilder()
	user := &User{}

	res, err := qb.Match(user).As("u").
		Return(types.Entity{Struct: user, Alias: "u"}).
		Build()

	assert.NoError(t, err)
	assert.Contains(t, res.Query, "MATCH (u:User:Person)")
	assert.Contains(t, res.Query, "RETURN u.username, u.email, u.age")
}

func TestWithEntity(t *testing.T) {
	qb := builder.NewQueryBuilder()
	user := &User{}

	res, err := qb.Match(user).As("u").
		With(types.Entity{Struct: user, Alias: "u"}).
		Return("u.username").
		Build()

	assert.NoError(t, err)
	assert.Contains(t, res.Query, "MATCH (u:User:Person)")
	assert.Contains(t, res.Query, "WITH u.username, u.email, u.age")
	assert.Contains(t, res.Query, "RETURN u.username")
}
