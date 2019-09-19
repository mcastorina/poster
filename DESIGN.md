# Design
poster is designed to be like an API with CRUD operations. The
following document explains the available resources and actions.

## Resources

- **Target:**
The most basic resource available is a target which is simply an
alias for an endpoint.
- **Request:**
A request resource combines a target with parameters, headers, a
body, and a method.
- **Authorization:**
An authorization resource defines the source of an authorization
token.
- **Suite:**
An ordered list of requests to send.

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
