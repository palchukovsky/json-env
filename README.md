[![Go Report Card](https://goreportcard.com/badge/github.com/palchukovsky/json-env)](https://goreportcard.com/report/github.com/palchukovsky/json-env)
# json-env
A tool to read and modify Base64 encoded JSON to use configurations JSON-files for CI/CD environments.
## To install
[Go](https://golang.org/dl/) has to be installed.
```shell
go get github.com/palchukovsky/json-env
```
## Examples
### Set value
```shell
json-env -source eyJmIjoidmFsMiIsIngiOnsieSI6eyJ6IjoidmFsMSJ9fX0 -write "x/y/z=1 2 3" z=valX
```
### Read value
```shell
json-env -source eyJmIjoidmFsMiIsIngiOnsieSI6eyJ6IjoiMSAyIDMifX0sInoiOiJ2YWxYIn0 -read x/y/z
json-env -source eyJmIjoidmFsMiIsIngiOnsieSI6eyJ6IjoiMSAyIDMifX0sInoiOiJ2YWxYIn0 -read z
json-env -source eyJmIjoidmFsMiIsIngiOnsieSI6eyJ6IjoiMSAyIDMifX0sInoiOiJ2YWxYIn0 -read zz -default "is not existent"
```
### Export source
```shell
json-env -source eyJmIjoidmFsMiIsIngiOnsieSI6eyJ6IjoiMSAyIDMifX0sInoiOiJ2YWxYIn0 -export
```
### Encode JSON-file to Base64
```shell
json-env -encode_file /path/config.json
```
### Use-case example
```Makefile
CONFIG ?= $(shell json-env -encode_file ./.config.json)
AWS_APP_GATEWAY_ID := $(shell json-env -source ${CONFIG} -read service/aws/gateway/app/id)
AWS_REGION := $(shell json-env -source ${CONFIG} -read service/aws/region)
CONFIG := $(shell json-env -source ${CONFIG} -write \
    "build/version=${VERSION}" \
    "build/commit=${COMMIT}" \
    "build/id=${BUILD}" \
    "build/builder=${BUILDER}" \
    "build/maintainer=${MAINTAINER}" \
    "service/aws/gateway/app/endpoint=${AWS_APP_GATEWAY_ID}.execute-api.${AWS_REGION}.amazonaws.com/${VERSION}")
AWS_ACCESS_KEY_ID := $(shell json-env -source ${CONFIG} -read builder/aws/accessKey/id)
AWS_SECRET_ACCESS_KEY := $(shell json-env -source ${CONFIG} -read builder/aws/accessKey/secret)
AWS_ACCOUNT_ID := $(shell json-env -source ${CONFIG} -read service/aws/accountId)
AWS_AUTH_GATEWAY_ID := $(shell json-env -source ${CONFIG} -read service/aws/gateway/auth/id)
    
define create-config-file
    -mkdir -p ${1}
    @echo "json-env -source CONFIG -export > ${1}config.json"
    @json-env -source ${CONFIG} -export > "${1}config.json"
endef
```
``` go
func newService(projectPackage string) service {
	name := os.Args[0]

	var config Config
	{
		file, err := ioutil.ReadFile("./lambda_config.json")
		if err != nil {
			log.Fatalf(`Failed to read config file: "%v".`, err)
		}
		if err := json.Unmarshal(file, &config); err != nil {
			log.Fatalf(`Failed to parse config file: "%v".`, err)
		}
	}

	return service{
		name:   name,
		log:    newProductLog(projectPackage, name, config),
		config: config.Service,
		build:  config.Build,
	}
}
```
