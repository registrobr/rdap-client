rdap-client
===========

This is a command line RDAP client.


Install
-------

First of all, you will need **Go** and **git** installed in your machine.
Instructions for intalling Go can be found in the link bellow.

http://golang.org/doc/install

You must be running Go version 1.5 or above.
Don't forget to create the $GOPATH environment variable.

Now just retrieve and install the project with the following command:
```
go get github.com/registrobr/rdap-client
```

Remember to add your **$GOPATH/bin** to your **$PATH** environment.

A Dockerfile is available and can be used as in the follow example:

```
docker build -t rdap-client:latest .

docker run -it --rm --name rdap-client rdap-client:latest registro.br
```


Usage
-----

To query something using bootstrap strategy:

```
rdap-client 199.71.0.160
```

Or if you want to directly query a RDAP server:

```
rdap-client -H rdap.registro.br nic.br
```

You can check more options with:

```
rdap-client -h
```
