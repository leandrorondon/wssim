# wssim

wssim is a simple generic webservice simulator.

## Introduction 

wssim is a simple generic webservice simulator.

You can use it to easily serve static web files (http://<HOST>:<PORT>) or simulate responses from a webservice (http://<HOST>:<PORT>/api/<FUNCTION>) developing client-side application.

## Directories

### web

Put your static html/css/js files here. You can use wssim as a simple html webserver to test your scripts against the simulated webservice.

### responses

responses' subdirectories (GET, POST, PUT, DELETE, HEAD) contain files (*.json) which content will be used to response a request to API functions.

For example, if an application makes a GET request to http://localhost:8099/api/status, wssim will try to read file "responses/GET/status.json" and return file content in the response body.

## Instructions

### Usage

wssim [OPTIONS]

wssim listens in port 8080 by default. It can be changed using -p option.

#### Options
    -p, --port set wssim to listen in a custom port
    -h, --help print the help

### Status code

The return code for any request can be set in the "statuscode.txt" file. If the file don't exists or couldn't be parsed, wssim returns 200 (OK).

### Content Type

The response's Content-Type matches the first Accept header value from request. If no Accept is specified ("*/*/"), wssim sets Content-Type to "application/json".

### Response body

Response body is read from files in responses directory.

For example, if an application makes a GET request to http://localhost:8099/api/status, wssim will try to read file "responses/GET/status.json" and will return file content in the response body.

wssim doesn't care about request's parameters or data. It neither checks if Content-Type value matches with response body content. It simply replies the request with content from corresponding file.

