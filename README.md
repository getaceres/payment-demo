# payment-demo

This is a demo of a RESTful implemention for payment information. API documentation can be found in swagger JSON format in the [swagger.json](swagger.json) file and in PDF format in the [swagger.pdf](doc/swagger.pdf) file.

To build it Go 1.11+ is required since it makes use of modules. For the same reason the code must be placed outside the GOPATH or it must be built inside the GOPATH with the GO111MODULE=on environment variable.

To run, this demo relies on the availability of a MongoDB instance.

Once built the REST server can be started with the ```payment-demo serve``` command. For a list of available flags, use the command ```payment-demo help serve```. Available flags are:
- ```--mongourl``` or ```-m```: Sets the connection URL for the backend MongoDB persistence storage. Defaults to ```mongodb://localhost:27017```
- ```--port``` or ```-p```: Sets the port in which the server will listen for connections. Defaults to ```8080```