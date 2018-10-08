My attempt at recreating https://github.com/eykrehbein/strest/blob/master/SCHEMA.md#delay

it mostly works as expected.

usage 

`strest -script req.yml`

TODO

example configuratation

```
version: 1
requests:
  userRequest:
    failOnError: true
    url: 'http://localhost:8080/user'
    method: POST
    data:
        params: 
            name: testUser
        headers: 
            Authentication: Bearer
            SomemoreHeader: "asfsdfds"
        form: 
            username: kingwill101
            password: somepassword
    log: true
    validation:
        body: "Hello, World!"
        statusCode: 200
    repeat: 3
    delay: 500
```

TODO

allow the use of environment variables in strings
Allow the specifying of types