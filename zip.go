//package zipper provides functions to zip and unzip files and/or folders
package zipper

import (
	"archive/zip"
	"compress/flate"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

/*
unzipFile unzips a given zipped file to a given destination
It accepts a filepath, ie. The path of the zipped file then unzips it and
creates a new folder in a given destination which is the second parameter.
It returns a boolean and an error respectively.
Boolean is returned when the unzipping is successful and error is returned when it is not
*/
func UnzipIt(filePath, fileDestination string) (final bool, finalError error) {
	file, fileError := os.Open(filePath)
	if fileError != nil {
		return false, errors.New("filepath doesn't exist")
	}
	defer func() {
		if fileError == nil {
			closeError := file.Close()
			if closeError != nil {
				log.Fatal(closeError)
			}
		}
	}()
	fileName := file.Name()
	if strings.Contains(fileName, "zip") {
		fileName = strings.TrimRight(fileName, ".zip")
	} else if strings.Contains(fileName, "tar") {
		//sucess, tarError :=	zipTar(filePath, fileDestination)
		//if tarError != nil {
		//		return false, tarError
		//}
		//return sucess, nil
	}
	fileDestination = filepath.Join(fileDestination, fileName)
	closeError := os.MkdirAll(fileDestination, 0777)
	if closeError != nil {
		log.Fatal(closeError)
	}
	readZip, readError := zip.OpenReader(filePath)
	if readError != nil {
		return false, errors.New("File cannot be read. Please check to make sure it is a zipped file.")
	}
	defer func() {
		if readError == nil {
			closeError := readZip.Close()
			if closeError != nil {
				log.Fatal(closeError)
			}
		}
	}()

	for _, readFile := range readZip.File {
		if os.FileMode.IsDir(readFile.FileInfo().Mode()) {
			dirError := os.MkdirAll(filepath.Join(fileName, readFile.Name), readFile.FileInfo().Mode())
			if dirError != nil {
				log.Fatal(dirError)
			}
			continue
		}
		reader, err := readFile.Open()
		if err != nil {
			log.Fatal(err)
		}
		path := filepath.Join(fileDestination, strings.TrimLeft(readFile.Name, "/"))
		fileWriter, errorFile := os.OpenFile(path, os.O_CREATE, readFile.Mode())
		if errorFile != nil {
			log.Fatal(errorFile)
		} else {
			closeError := fileWriter.Close()
			if closeError != nil {
				log.Fatal(closeError)
			}
		}
		_, err = io.Copy(fileWriter, reader)
		if err != nil {
			log.Fatal(err)
		}
	}

	return true, nil

}

func writeFolder(folder, mainPath string, theWriter *zip.Writer) {
	directory, directoryError := ioutil.ReadDir(filepath.Join(mainPath, folder))
	if directoryError != nil {
		log.Fatal(directoryError)
	}
	for _, theFile := range directory {
		if theFile.IsDir() {
			_, writeError := theWriter.Create(folder + theFile.Name() + "/")
			if writeError != nil {
				log.Fatal(writeError)
			}
			writeFolder(filepath.Join(folder, theFile.Name()+"/"), mainPath, theWriter)
			continue
		}
		y, yer := theWriter.Create(filepath.Join(folder, theFile.Name()))
		if yer != nil {
			log.Fatal(yer)
		}
		b, ero := ioutil.ReadFile(filepath.Join(mainPath, filepath.Join(folder, theFile.Name())))
		if ero != nil {
			log.Fatal(ero)
		}
		writeFile(y, b)


	}
}

func writeFile(writer io.Writer, data []byte) {
	_, err := writer.Write(data)
	if err != nil {
		log.Fatal(err)
	}
}

/*
zipIt zips a given folder or file to a given destination.
It accepts a filepath, ie. The path of the folder or file then zips it and
places it in a given destination, which is the second parameter.
It returns the destination of the zipped file as a string
*/
func ZipIt(filePath, fileDestination string, zipFileName string) (Destination string, Error error) {
	fileReader, errorReader := os.Open(filePath)
	if errorReader != nil {
		log.Fatal(errorReader)
		return "", errorReader
	}

	fileDirectory := filepath.Dir(fileReader.Name())
	fileNameDefault := strings.TrimLeft(fileReader.Name(), fileDirectory) + ".zip"
	fmt.Println(fileNameDefault)
	if fileDestination == "" && zipFileName == "" {
		Destination = filepath.Join(fileDirectory, fileNameDefault)
	} else if zipFileName != "" && fileDestination != "" {
		Destination = filepath.Join(fileDestination, zipFileName)
	} else if fileDestination != "" && zipFileName == "" {
		Destination = filepath.Join(fileDestination, fileNameDefault)
	} else {
		Destination = filepath.Join(fileDirectory, zipFileName)
	}
	fileWriter, errorFile := os.OpenFile(Destination, os.O_CREATE, 0666)
	fileData, errorFileData := fileReader.Stat()
	w := zip.NewWriter(fileWriter)
	w.RegisterCompressor(zip.Deflate, func(out io.Writer) (io.WriteCloser, error) {
		return flate.NewWriter(out, flate.BestCompression)
	})
	if errorFileData != nil {
		log.Fatal(errorFileData)
	}
	if fileData.IsDir() {
		fileContents, errorFileContents := ioutil.ReadDir(filePath)
		if errorFileContents != nil {
			log.Fatal(errorFileContents)
		}
		for _, f := range fileContents {
			if f.IsDir() {
				_, createError := w.Create(f.Name() + "/")
				if createError != nil {
					log.Fatal(createError)
				}
				writeFolder(f.Name()+"/", filePath, w)
				continue
			}
			y, yer := w.Create(f.Name())
			if yer != nil {
				log.Fatal(yer)
			}
			b, ero := ioutil.ReadFile(filepath.Join(filePath, f.Name()))
			if ero != nil {
				log.Fatal(ero)
			}
			writeFile(y, b)

		}
	} else {
		y, yer := w.Create(fileData.Name())
		if yer != nil {
			log.Fatal(yer)
		}
		b, ero := ioutil.ReadFile(filepath.Join(filePath, fileData.Name()))
		if ero != nil {
			log.Fatal(ero)
		}
		writeFile(y, b)
	}

	if errorFile != nil {
		log.Fatal(errorFile)
	}
	defer func() {
		ert := w.Close()
		if ert != nil {
			log.Fatal(ert)
		}
	}()

	return

}

