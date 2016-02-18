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


Usage
-----

To query something using bootstrap strategy:

```
rdap-client 199.71.0.160
```

Or if you want to directly query a RDAP server:

```
rdap-client -H rdap.beta.registro.br nic.br
```

You can check more options with:

```
rdap-client -h
```
