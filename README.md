# Zipper
A zipping tool for golang for compressing zipped files even with deeply nested folders as well as decompressing them with one line of code

# Installation
```shell script
go get "github.com/augani/zipper"
```

# Usage
```go
import "github.com/augani/zipper"
```
## Functions available
* ZipIt
* UnZipIt

### ZipIt function
```go
// You pass in the filepath of the file of folder that you want to zip
// as well as the destination folder, where you want it to be saved and
// also the name of the zipped file
//This function returns an error or a destination of type string
//You can zip a whole folder with deeply embedded folders and files with no work
destination, err := ZipIt(filePath, Destination, zipfileName)
```
>Note that you can pass  an empty string to the destination and zipFileName if you want to zip the folder of file and have it in the same directory as the file path.

### UnZipIt function
```go
// You pass in the filepath of the file of folder that you want to unzip
// as well as the destination folder, where you want it to be saved and
// also the name of the zipped file
//This function returns an error or a boolean
//You can zip a whole folder with deeply embedded folders and files with no work
//The destination is optional so give it an empty string if you want the file to be unzipped into the current directory
done, err := UnZipIt(filePath, Destination)
```

