# poster
API testing aid from the command line.

The purpose of this tool is to help organize and send HTTP requests.
More specifically, it aims to help developers test their endpoints.
It is intended to be easy to use and intuitive.

## Design Overview
This section provides a summary of the design to aid usage. For a
more detailed explanation, see [DESIGN.md](DESIGN.md).

`poster` is designed in a similar way to a CRUD application, so
there are four subcommands to modify resources: `create`, `get`,
`edit`, and `delete`.

The three main types of resources are: `environment`, `variable`,
and `request`. Each `variable` resources belongs to an `environment`,
and a `request` runs in an `environment` as well. This allows for
automatic dependency generation for a request.

**Example**

Suppose you have a request that needs an authorization header. This
header changes when you are testing between staging and production
environments. Below are the resources to define this scenario.

```
Variables
===============================================================
Name                Environment         Generator       Type
auth-header         staging             get-token       request
auth-header         production          get-token       request
apikey              staging             STAGING_KEY     const
apikey              production          PROD_KEY        const

Requests
==========================================================================================
Name            Method          URL                             Header
get-token       GET             example.com?apikey=:apikey
check-auth      GET             example.com                     Authorization=:auth-header
```

For the purpose of this example, only the authorizaion header is
changing for the `check-auth` request. The following steps happen when
`poster` runs `check-auth` in the `staging` environment.

1. Find the `:auth-header` variable and generate its value
   1. `:auth-header` depends on the `get-token` request
   1. Run the `get-token` request
      1. Find the `:apikey` variable and replace the value with `STAGING_KEY`
      1. Send the GET request to `example.com?apikey=STAGING_KEY`
   1. Save the result to the `auth-header` variable
2. Send the GET request to `example.com`

## Usage
The following subcommands are used by `poster` to modify resources
as well as run requests. Each one has a help command to provide CLI
usage. Once the resources are created, the most common command will
most likely be the run command.

```
  create      Create a resource
  delete      Delete resources
  edit        Modify a resource
  get         Print resources
  run         Execute the named resource
```

## Motivation
I wanted an easy way to repeatedly send curl commands for different environments.

## Project Structure
The following describes the layout of this project.

 - cmd - Contains the main package
 - internal - Contains internal logic
   - cli - Code related to parsing arguments
   - models - Various struct definitions
   - store - Database specific code
   - cache - Interface between models and store
