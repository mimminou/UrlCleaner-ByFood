## Golang Url Cleaner Micro service

## Running the service
use `go run .` at the root folder to run the server from terminal.

The server has an stdout to the console that prints incoming requests and their responses with Timestamp, Endpoint, HTTP Method as well as the body of the response requests.

## Project Structure

- `/server/*`: contains the server logic and middleware implementations

## API Endpoints

### Endpoints:
- `/`: `POST` : takes a JSON body with the Request structure (as specified below)


### Request Json Structure : 

```
{
"url" : String (valid URL),
"operation" : string (one of the following : "canonical", "redirection" "all")
}

```

### Response Json Structure :

- on success :   `{"processed_url" : string}`


- on error :  `{"msg": string}`


### Operations:

##### `canonical` : Cleans up URL, removes Query parameters and trailing slashes.

example : `https://BYFOOD.com/food-EXPeriences?query=abc/`  -->  `https://BYFOOD.com/food-EXPeriences`



##### `redirection` : Checks if the domain is `byfood.com` (as a security measure, will only work for www subdomain, and add it to the URL if it's missing), then returns the url lowercased.

example : `https://BYFOOD.com/food-EXPeriences?query=abc/` --> `https://www.byfood.com/food-experiences?query=abc/`


##### `all` : Performs both Canonical and Redirection, as such, it's bound to only work on `byfood.com` domain

example : `https://BYFOOD.com/food-EXPeriences?query=abc/` --> `https://www.byfood.com/food-experiences`


# Tests

Unit tests were written for this service

simply navigate to /server and run `go test -v`



