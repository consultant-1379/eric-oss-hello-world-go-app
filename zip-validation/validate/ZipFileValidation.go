// COPYRIGHT Ericsson 2023

// The copyright to the computer program(s) herein is the property of
// Ericsson Inc. The programs may be used and/or copied only with written
// permission from Ericsson Inc. or in accordance with the terms and
// conditions stipulated in the agreement/contract under which the
// program(s) have been supplied.

package validate

import (
	"archive/zip"
	"bufio"
	"bytes"
	"fmt"
	log "gerrit-review.gic.ericsson.se/cloud-ran/src/golang-log-api/logapi"
	"io"
	"os"
	"strings"
)

/**
This is the validation code for zip file that we make available to customers, this code is only meant for internal Ericsson builds so is being
kept outside off the 'src/' directory

The PreValidationChecks.go file should always be executed first as it's the acid test for what a valid zip file is

The obligatory files and directories are defined in the 'bob-properties.yaml' file in the 'required-dirs' and 'required-files' props for example:
  - required-dirs: 'charts/,src/,vendor/,csar/'
  - required-files: 'Dockerfile-template,go.mod,go.sum,README.md,version'

This GO file will throw a panic and fail the current publish task if:
	1. any of the files listed in the 'required-files' property are missing from the zip or if the file is present but empty
	2. any of the directories listed in the 'required-dirs' property are absent or if the directory contains no files (empty sub-directories counted as empty)

NOTES:
	1. each 'requiredDir' value has to have a '/' postfix (consequence of how the zip handling code names directories internally)
	2. required files and directories are assumed to be only at the top level of the zip, nested files/directories are not validated!
*/

// ValidateZip will ensure that all the expected files and directories are present in the customer zip
func ValidateZip(zipPath string, requiredFiles string, requiredDirs string) {
	log.Info("validating zip path    : " + zipPath)
	log.Info(fmt.Sprintf("required files         : %v", requiredFiles))
	log.Info(fmt.Sprintf("required non-empty dirs: %v", requiredDirs))

	requiredDirsFound := make([]string, 0)
	requiredFilesFound := make([]string, 0)
	requiredFilesList := splitAndTrim(requiredFiles)
	requiredDirsList := splitAndTrim(requiredDirs)

	zipFile := loadZip(zipPath)
	for _, zipFile := range zipFile {
		fileName := zipFile.Name
		fileSize := zipFile.UncompressedSize64

		// directory's size is always '0' so only way to test for a non-empty directory is find a file that has it as parent
		isDir := zipFile.FileInfo().IsDir()
		if !isDir {
			if foundRequiredFile(requiredFilesList, fileName, fileSize) {
				requiredFilesFound = append(requiredFilesFound, fileName)
			} else {
				// check if the file is a child of a required directory
				index := strings.Index(fileName, "/")
				if index > 0 {
					topParent := fileName[:index+1]
					if contains(requiredDirsList, topParent) && !contains(requiredDirsFound, topParent) {
						requiredDirsFound = append(requiredDirsFound, topParent)
					}
				}
			}
		}
	}

	// now assert the contents of the zip
	if len(requiredDirsFound) != len(requiredDirsList) {
		errorMessage := fmt.Sprintf("required non-empty dirs not satisfied! required: '%v', got: '%v'", requiredDirs, requiredDirsFound)
		panic(errorMessage)
	}
	if len(requiredFilesFound) != len(requiredFilesList) {
		errorMessage := fmt.Sprintf("required files not satisfied! required: '%v', got: '%v'", requiredFiles, requiredFilesFound)
		panic(errorMessage)
	}

	log.Info("zip successfully validated!")
}

func splitAndTrim(values string) []string {
	var list = make([]string, 0)
	for _, val := range strings.Split(values, ",") {
		list = append(list, strings.TrimSpace(val))
	}
	return list
}

func foundRequiredFile(requiredList []string, fileName string, size uint64) bool {
	if contains(requiredList, fileName) {
		if size > 0 {
			return true
		}
		panic("required file: '" + fileName + "' is empty")
	}
	return false
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

func loadZip(path string) []*zip.File {
	file, err := os.Open(path)
	check(err)
	reader := bufio.NewReader(file)
	body, err := io.ReadAll(reader)
	check(err)
	zipReader, err := zip.NewReader(bytes.NewReader(body), int64(len(body)))
	check(err)
	return zipReader.File
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
