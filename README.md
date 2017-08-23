## POC


### Prerequisites and setup:

* [Docker](https://www.docker.com/products/overview) - v1.12 or higher
* [Docker Compose](https://docs.docker.com/compose/overview/) - v1.8 or higher
* [Git client](https://git-scm.com/downloads) - needed for clone commands
* **Node.js** v6.9.0 - 6.10.0 ( __Node v7+ is not supported__ )
* **jq**  apt-get install jq
* Download docker images


## Running the sample program


### Step 1: setup the fabric network

##### Terminal Window 1

* start sdk

```
cd sdk
sudo su
./runApp.sh
```
##### Terminal Window 2

* init sdk

```
cd sdk
./init.sh
```

### Step 2: setup the yiliao poc

##### Terminal Window 3

```
cd yiliao
sudo su
./runApp.sh
```

## Test the sample program
```
http://127.0.0.1:3389
```