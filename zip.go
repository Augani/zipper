//package zipper provides functions to zip and unzip files and/or folders
package zipper

import (
	"archive/zip"
	"compress/flate"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

/*
UnZipIt unzips a given zipped file to a given destination
It accepts a filepath, ie. The path of the zipped file then unzips it and
creates a new folder in a given destination which is the second parameter.
It returns a boolean and an error respectively.
Boolean is returned when the unzipping is successful and error is returned when it is not
*/
func UnZipIt(filePath, fileDestination string) (final bool, finalError error) {
	file, fileError := os.Open(filePath)
	checkError(fileError)
	defer func() {
		closeError := file.Close()
		checkError(closeError)
	}()
	fileName := file.Name()
	if fileDestination == "" {
		fileDirectory := filepath.Dir(fileName)
		fileFolder := strings.TrimSuffix(strings.TrimLeft(fileName, fileDirectory), ".zip")
		fileDestination = filepath.Join(fileDirectory, fileFolder)
	}
	closeError := os.MkdirAll(fileDestination, 0777)
	checkError(closeError)
	readZip, readError := zip.OpenReader(filePath)
	checkError(readError)
	defer func() {
		closeError := readZip.Close()
		checkError(closeError)
	}()
	for _, readFile := range readZip.File {
		folders := strings.Split(readFile.Name, "/")
		b := fileDestination
		if len(folders) > 1 {
			for r := 0; r < len(folders)-1; r++ {
				dError := os.MkdirAll(filepath.Join(b, folders[r]), 0777)
				checkError(dError)
				b = filepath.Join(b, folders[r])
			}
		}
		if os.FileMode.IsDir(readFile.FileInfo().Mode()) {
			dirError := os.MkdirAll(filepath.Join(fileDestination, readFile.Name), readFile.FileInfo().Mode())
			checkError(dirError)
			continue
		}
		reader, err := readFile.Open()
		checkError(err)
		path := filepath.Join(fileDestination, strings.TrimLeft(readFile.Name, "/"))
		fileWriter, errorFile := os.OpenFile(path, os.O_CREATE, readFile.Mode())
		checkError(errorFile)
		_, err = io.Copy(fileWriter, reader)
		checkError(err)
		closeError := fileWriter.Close()
		checkError(closeError)

	}
	return true, nil

}

func writeFolder(folder, mainPath string, theWriter *zip.Writer) {
	directory, directoryError := ioutil.ReadDir(filepath.Join(mainPath, folder))
	checkError(directoryError)
	for _, theFile := range directory {
		if theFile.IsDir() {
			_, writeError := theWriter.Create(folder + theFile.Name() + "/")
			checkError(writeError)
			writeFolder(filepath.Join(folder, theFile.Name()+"/"), mainPath, theWriter)
			continue
		}
		if strings.HasPrefix(theFile.Name(), "~") {
			continue
		}
		y, yer := theWriter.Create(filepath.Join(folder, theFile.Name()))
		checkError(yer)
		b, ero := ioutil.ReadFile(filepath.Join(mainPath, filepath.Join(folder, theFile.Name())))
		checkError(ero)
		writeFile(y, b)
	}
}
func writeFile(writer io.Writer, data []byte) {
	_, err := writer.Write(data)
	checkError(err)
}

/*
ZipIt zips a given folder or file to a given destination.
It accepts a filepath, ie. The path of the folder or file then zips it and
places it in a given destination, which is the second parameter.
It returns the destination of the zipped file as a string
*/
func ZipIt(filePath, fileDestination string, zipFileName string) (Destination string, Error error) {
	fileReader, errorReader := os.Open(filePath)
	checkError(errorReader)
	fileDirectory := filepath.Dir(fileReader.Name())
	fileNameDefault := strings.TrimPrefix(fileReader.Name(), fileDirectory) + ".zip"

	if fileDestination == "" && zipFileName == "" {
		Destination = filepath.Join(fileDirectory, fileNameDefault)
	} else if zipFileName != "" && fileDestination != "" {
		Destination = filepath.Join(fileDestination, zipFileName) + ".zip"
	} else if fileDestination != "" && zipFileName == "" {
		Destination = filepath.Join(fileDestination, fileNameDefault)
	} else {
		Destination = filepath.Join(fileDirectory, zipFileName) + ".zip"
	}
	fileWriter, errorFile := os.OpenFile(Destination, os.O_CREATE, 0666)
	fileData, errorFileData := fileReader.Stat()
	w := zip.NewWriter(fileWriter)
	w.RegisterCompressor(zip.Deflate, func(out io.Writer) (io.WriteCloser, error) {
		return flate.NewWriter(out, flate.BestCompression)
	})
	checkError(errorFileData)
	if fileData.IsDir() {
		fileContents, errorFileContents := ioutil.ReadDir(filePath)
		checkError(errorFileContents)
		for _, f := range fileContents {
			if f.IsDir() {
				_, createError := w.Create(f.Name() + "/")
				checkError(createError)
				writeFolder(f.Name()+"/", filePath, w)
				continue
			}
			y, yer := w.Create(f.Name())
			checkError(yer)
			b, ero := ioutil.ReadFile(filepath.Join(filePath, f.Name()))
			checkError(ero)
			writeFile(y, b)
		}
	} else {

		y, yer := w.Create(fileData.Name())
		checkError(yer)
		b, ero := ioutil.ReadFile(filepath.Join(filePath, fileData.Name()))
		checkError(ero)
		writeFile(y, b)
	}
	checkError(errorFile)
	defer func() {
		ert := w.Close()
		checkError(ert)
	}()
	return
}

func checkError(err error) {
	if err != nil {
		fmt.Println("An error occurred:", error.Error(err))
		os.Exit(0)
	}
}
