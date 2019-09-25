# poster
API testing aid from the command line.

The purpose of this tool is to help organize and send HTTP requests.
More specifically, it aims to help developers test their endpoints.
It is intended to be easy to use and intuitive.

## Usage
TODO

## Motivation
I wanted an easy way to repeatedly send curl commands for different environments.

## Project Structure
The following describes the layout of this project.

 - cmd - Contains the main package
 - internal - Contains internal logic
   - cli - Code related to parsing arguments
   - models - Various struct definitions
   - store - Database specific code
