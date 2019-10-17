package models

func RequestTemplate() Request {
	return Request{
		Name:        "request template",
		Method:      "GET",
		URL:         "http://localhost",
		Environment: Environment{Name: "local"},
		Headers: []Header{
			{Key: "Content-Type", Value: "application/json"},
		},
	}
}

func EnvironmentTemplate() Environment {
	return Environment{
		Name: "environment template",
	}
}

func ConstVariableTemplate() Variable {
	return Variable{
		Name:        "const variable template",
		Value:       "value",
		Environment: Environment{Name: "local"},
		Type:        ConstType,
	}
}

func ScriptVariableTemplate() Variable {
	return Variable{
		Name:        "script variable template",
		Environment: Environment{Name: "local"},
		Type:        ScriptType,
		Generator: &VariableGenerator{
			Script: `date +'%D %T'`,
		},
	}
}

func RequestVariableTemplate() Variable {
	return Variable{
		Name:        "request variable template",
		Environment: Environment{Name: "local"},
		Type:        RequestType,
		Generator: &VariableGenerator{
			RequestName:        "request-name",
			RequestEnvironment: "environment-name",
			RequestPath:        "$",
		},
	}
}
