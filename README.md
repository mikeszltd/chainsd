Chainsd JSON processor [![Build Status](https://drone.io/github.com/mikeszltd/chainsd/status.png)](https://drone.io/github.com/mikeszltd/chainsd/latest)
======================

Chainsd is to be ultimate monitoring and data processing tool focused on JSON messages. 

Processing pipeline consists of multiple `Chains` of `Links`. There are 3 types of `Links` - `generators`, `filters` and `outputs`.

# Chain

Chain is made of `Links` - `generators`, `filters` and `outputs`. 
You can have multiple chains running simultaneously. 
Each chain has to start with `generator` followed by `output`. You can have multiple filters between `generator` and `output`.

## Generators

Generators provide source of JSON messages. There are multiple types of `generators` built-in. 

  - `heartbeat`

    Generates periodical JSON message of following format:
  
    ```
    {"timestamp": <current unix timestamp>, message: <defined message>}
    ```
  
  - `httpin`
  
    Listens on defined port for `application/json` messages.

  - `amqpin`
  
    Connects to `RabbitMQ` queue and pulls messages from specified route

  - `tail`
  
    Runs `tail` on specified file and "scraps" data into JSON using provided regular expression or pre-defined template

  - `fanin`
  
    Receives message from another `Chain` 

## Filters

Currently there is only one filter, that doesn't filter yet.

  - `fanout`
  
    Forwards the message to the `fanin` generator of defined name

## Outputs

  - `stdout`
    
    Prints message to screen

  - `fileout`
  
    Writes message to file

  - `hdfsout`
  
    Writes message to `Hadoop` filesystem. Filenames can be created using message's JSON fields

  - `mongoout`
    
    Writes message to `MongoDB` to specified Database and Collection.

  - `amqpout`
  
    Writes message to specified route of `RabbitMQ` queue

  - `nullout`
  
    Writes message into parallel universe

# Installation

  To install `Chainsd` you need Go and number of tools. Instructions below are for `Ubuntu 14.04`.
  
  Install `Go` following instructions here: http://golang.org/doc/install
  
  Install dependencies:
  
  ```
  sudo apt-get install git mercurial bzr
  ```
  
  Once installed, you can get the project to your Go environment:
  
  ```
  go get github.com/mikeszltd/chainsd
  ```
  
  And install it:
  
  ```
  go install github.com/mikeszltd/chainsd
  ```
  
  Now create sample configuration file `config.json`:
  
  ```
  {
	"chains": [{
		"heartbeat": [
			{
				"type": "heartbeat",
				"interval": "1s",
				"message": "server1"
			},
			{
				"type": "stdout"
			}
		]
	}]
  }
  ``` 
  
  and run:
  
  ```
  ./bin/chainsd
  ```
  
  You should see something like this:
  
  ```
  Chainsd by Mikesz Ltd
  {"timestamp":1403828540,"message":"server1"}
  {"timestamp":1403828541,"message":"server1"}
  {"timestamp":1403828542,"message":"server1"}
  {"timestamp":1403828543,"message":"server1"}
  ```

# Todo

  More documentation
  
  Filters
  
  
