[![Go Report Card](https://goreportcard.com/badge/github.com/palchukovsky/json-env)](https://goreportcard.com/report/github.com/palchukovsky/json-env)

# json-env
A tool to read and modify Base64 encoded JSON to use configurations JSON-files for CI/CD environments.

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