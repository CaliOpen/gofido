Description
===========

Fido u2f HTTP API server in golang using a cassandra backend.

Why ?
=====

We [Caliopen](https://caliopen.org) use cassandra as main storage backend and a using an u2f server persisting its data into this kind of storage engine is what we need.

So this project to make such service, where your user_id drive this u2f backend to fulfill u2f request/response challenges for registration and sign actions.


Know problems
-------------

- Make AppId not a configuration value but a variable (or not)
- Missing tests
- Missing valuable documentation


Usage
-----

Bring up storage
~~~~~~~~~~~~~~~~

`
$ docker-compose cassandra
`

see (1)

Setup
~~~~~

Only on first use

`
docker-compose exec cassandra cqlsh -f /schema.cql
`

You will need a certificate and its private key in `certs` directory to run the server.
https is strictly needed to operate FIDO U2F, or strange errors occurs in the browser.

`
mkdir certs
openssl req -newkey rsa:2048 -nodes -keyout certs/key.pem -x509 -days 365 -out certs/cert.pem
`

Run
~~~~


`
docker-compose up gofido
`

*Notes*

(1) Cassandra backend take some seconds to operate, be patient
