# zipper
A zipping tool for golang

##Functions available
* zipIt
* unzipIt

###zipIt function
```go
// You pass in the filepath of the file of folder that you want to zip
// as well as the destination folder, where you want it to be saved and
// also the name of the zipped file
//You can zip a whole folder with deeply embedded folders and files with no work
destination, err := zipIt(filePath, Destination, zipfileName string)
```
>Note that you can pass  an empty string to the destination and zipFileName if you want to zip the folder of file and have it in the same directory as the file path.

###unZipIt function
```go
// You pass in the filepath of the file of folder that you want to unzip
// as well as the destination folder, where you want it to be saved and
// also the name of the zipped file
//You can zip a whole folder with deeply embedded folders and files with no work
des, err := zipIt(filePath, Destination, zipfileName string)
```

