# Design
poster is designed to be like an API with CRUD operations. The
following document explains the available resources and actions.

## Resources

- **Target:**
The most basic resource available is a target which is simply an
alias for an endpoint.
- **Request:**
A request resource combines a target with a path, parameters, headers,
a body, and a method. It also has an environment associated with it so
variables can be generated accordingly.
- **Authorization:**
An authorization resource defines the source of an authorization
token.
- **Suite:**
An ordered list of requests to send.
- **Environment:**
An environment is used to logically define the construction of a
request.
- **Variable:**
A variable allows for easy reuse and generation of values and
is dependent on the environment. Variables have an optional expiration
so the program can automatically regenerate when needed.

## Actions

- **Create:**
Create a resource.
- **Get:**
View a resource.
- **Edit:**
Edit a resource.
- **Delete:**
Delete a resource.
- **Run:**
Run a resource such as a request or a suite.
- **Export:**
Export a resource in a specific format. This is useful for viewing
the corresponding `curl` request.
