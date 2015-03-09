# Database build command line tool

Simple command line tool to create postgres databases from a database.yml file

This was built as my "so what is this go-thing" project, but it has proven itself useful for CI.

## Usage

```
usage: dbbuilder <environment>
  -c=false: Don't exit on errors, useful for CI
  -e="test": Set the database test env to create
  -p="config/database.yml": Path to yaml (otherwise config/database.yml)
  -v=false: Prints current version
```

so running `dbbuilder -e test` inside a rails directory with a postgres database will create the test database.

## TODO

- [ ] Documentation
- [ ] Tests
- [ ] All the things
