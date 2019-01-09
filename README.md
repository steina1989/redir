# Redir

Url shortener API - A weekend project in Go

## Example application

Running on [redirdev.herokuapp.com](https://redirdev.herokuapp.com)

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes. See deployment for notes on how to deploy the project on a live system.

### Prerequisites for development

What things you need to install the software and how to install them

* [golang](https://golang.org/dl/)
* [postgres](https://www.postgresql.org/download/)
* [govendor](https://github.com/kardianos/govendor)*

*You can skip this one and just manually install the packages needed, but govendor works nicely with Heroku.

You will also need to create a postgres database with the code snippet inside createTables.sql

### Developing

Fetch dependencies
```
govendor migrate
govendor sync
```

Export neccessary environment variables

Build and run
```
go build && ./main
```

## Deployment

You will need the [herokuCLI](https://devcenter.heroku.com/articles/heroku-cli)


I won't get into Heroku details but when you have initialized and logged in to Heroku, and added a heroku remote:

Instead of using
```
go build && ./main
```
consider using a .env file with KEY=VALUE pairs instead of manually exporting the configuration variables and then run
```
go install && heroku local 
```
for local development 

```
git push heroku master
```

Use a heroku environmental variable to add a secret salt.

## Built With

* [Go](https://golang.org/) 
* [Gorilla/mux](https://github.com/gorilla/mux) URL router and dispatcher
* [Govendor](https://github.com/kardianos/govendor) - Dependency Management
* [Hashids](https://github.com/speps/go-hashids) - Used to generate hash tokens

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details

