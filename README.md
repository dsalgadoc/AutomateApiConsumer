<h1 align="center">
  My API Consumer
</h1>

<p align="center">
  This repository does various data retrievals for trusted sources automatically. It reads data from an input, makes a request, and records the result in a different file.
</p>

## Environment Setup

### Needed tools

1. Go 1.18 or higher

### Application execution

1. Make sure to download all Needed tools
2. Clone the repository
```
  git clone https://github.com/dasalgadoc/MyApiConsumer.git
```
3. Build up go project
```
  go get .
```
4. Make a **config.yaml** file in **./configs**
In **./configs/secure_config.yaml** has an example. Setup the data source clients and provide the location and name of input data.

The config files the following structure

| Father  |     Node     | Description                                                 | Possible values / Example                        |
|:-------:|:------------:|-------------------------------------------------------------|--------------------------------------------------|
|    -    |      io      | Describes Input parameters                                  | _Father node_                                    |
|   io    |   location   | Relative path where input can be found                      | ./file_exchange/                                 |
|   io    |  file_name   | Input filename (without extension)                          | input                                            |
|    -    |   clients    | Array with following elements, to describe each rest client | _Father node_                                    |
| clients |     name     | Name to invoke the client                                   | api-shazam                                       |
| clients |     path     | HTTP URL for client (without parameters or headers)         | https://shazam.p.rapidapi.com/shazam-events/list |
| clients |   headers    | Array with following elements, to describe each header      |                                                  |
| header  | :HEADER KEY: | Header name : header value                                  | 'X-RapidAPI-Key': 'SIGN-UP-FOR-KEY'              |

5. Run the API
```
  go run main.go
```
6. You'll see various errors since not enough parameters were provided; the correct form is as follows.
```
  go run main.go <INPUT FILE TYPE> <OUTPUT FILE TYPE> <DATA SOURCE>
```
For instance:
```
  go run main.go csv json api-engine 
```
