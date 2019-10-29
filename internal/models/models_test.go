package models

import (
	"testing"

	"github.com/bouk/monkey"
	"github.com/mcastorina/poster/internal/cache"
	"github.com/mcastorina/poster/internal/store"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	InitLogger()
	cache.InitLogger()
	store.InitLogger()
}

func TestReplaceVariablesOne(t *testing.T) {
	defer monkey.UnpatchAll()
	env := Environment{Name: "local"}

	patchGetVariablesByEnvironment := func(environment string) []store.Variable {
		host := store.Variable{
			Name:        "host",
			Value:       "localhost:8080",
			Environment: env.Name,
		}
		return []store.Variable{host}
	}
	monkey.Patch(store.GetVariablesByEnvironment, patchGetVariablesByEnvironment)

	input := "https://:host?param=true"
	actual := env.ReplaceVariables(input)
	expected := "https://localhost:8080?param=true"

	assert.Equal(t, expected, actual)
}
func TestReplaceVariablesTwo(t *testing.T) {
	defer monkey.UnpatchAll()
	env := Environment{Name: "local"}

	patchGetVariablesByEnvironment := func(environment string) []store.Variable {
		host := store.Variable{
			Name:        "host",
			Value:       "localhost",
			Environment: env.Name,
		}
		port := store.Variable{
			Name:        "port",
			Value:       "8080",
			Environment: env.Name,
		}
		return []store.Variable{host, port}
	}
	monkey.Patch(store.GetVariablesByEnvironment, patchGetVariablesByEnvironment)

	input := "https://:host::port"
	actual := env.ReplaceVariables(input)
	expected := "https://localhost:8080"

	assert.Equal(t, expected, actual)
}
func TestReplaceVariablesOverlap(t *testing.T) {
	defer monkey.UnpatchAll()
	env := Environment{Name: "local"}

	patchGetVariablesByEnvironment := func(environment string) []store.Variable {
		host := store.Variable{
			Name:        "host",
			Value:       "localhost::port",
			Environment: env.Name,
		}
		port := store.Variable{
			Name:        "port",
			Value:       "8080",
			Environment: env.Name,
		}
		return []store.Variable{host, port}
	}
	monkey.Patch(store.GetVariablesByEnvironment, patchGetVariablesByEnvironment)

	{
		input := "https://:host"
		actual := env.ReplaceVariables(input)
		expected := "https://localhost::port"
		assert.Equal(t, expected, actual)
	}
	{
		input := "https://:host::port"
		actual := env.ReplaceVariables(input)
		expected := "https://localhost::port:8080"
		assert.Equal(t, expected, actual)
	}
}

func BenchmarkReplaceVariables(b *testing.B) {
	defer monkey.UnpatchAll()
	env := Environment{Name: "local"}

	patchGetVariablesByEnvironment := func(environment string) []store.Variable {
		host := store.Variable{
			Name:        "host",
			Value:       "localhost::port",
			Environment: env.Name,
		}
		port := store.Variable{
			Name:        "port",
			Value:       "8080",
			Environment: env.Name,
		}
		return []store.Variable{host, port}
	}
	monkey.Patch(store.GetVariablesByEnvironment, patchGetVariablesByEnvironment)

	input := "https://:host::port"
	for i := 0; i < b.N; i++ {
		env.ReplaceVariables(input)
	}
}
func BenchmarkReplaceVariablesNone(b *testing.B) {
	defer monkey.UnpatchAll()
	env := Environment{Name: "local"}

	patchGetVariablesByEnvironment := func(environment string) []store.Variable {
		host := store.Variable{
			Name:        "host",
			Value:       "localhost::port",
			Environment: env.Name,
		}
		port := store.Variable{
			Name:        "port",
			Value:       "8080",
			Environment: env.Name,
		}
		return []store.Variable{host, port}
	}
	monkey.Patch(store.GetVariablesByEnvironment, patchGetVariablesByEnvironment)

	input := "GET"
	for i := 0; i < b.N; i++ {
		env.ReplaceVariables(input)
	}
}
func BenchmarkReplaceVariablesContainsColon(b *testing.B) {
	defer monkey.UnpatchAll()
	env := Environment{Name: "local"}

	patchGetVariablesByEnvironment := func(environment string) []store.Variable {
		host := store.Variable{
			Name:        "host",
			Value:       "localhost::port",
			Environment: env.Name,
		}
		port := store.Variable{
			Name:        "port",
			Value:       "8080",
			Environment: env.Name,
		}
		return []store.Variable{host, port}
	}
	monkey.Patch(store.GetVariablesByEnvironment, patchGetVariablesByEnvironment)

	input := "http://localhost:8080"
	for i := 0; i < b.N; i++ {
		env.ReplaceVariables(input)
	}
}
