// COPYRIGHT Ericsson 2023

// The copyright to the computer program(s) herein is the property of
// Ericsson Inc. The programs may be used and/or copied only with written
// permission from Ericsson Inc. or in accordance with the terms and
// conditions stipulated in the agreement/contract under which the
// program(s) have been supplied.

package validate

import (
	"fmt"
	log "gerrit-review.gic.ericsson.se/cloud-ran/src/golang-log-api/logapi"
)

// test variables
const testZipPath string = "zip-validation/test-archive.zip"
const validRequiredDirs string = "a/,b/"
const validRequiredFiles = "file_1.txt,file_2.txt"

func PreValidation() {
	log.Info("pre-validation steps...")

	validZipNoPanicExpected()
	invalidZipMissingRequiredFile()
	invalidZipMissingRequiredDirectory()
	invalidZipEmptyRequiredFile()
	invalidZipEmptyRequiredDirectory()
	invalidZipRequiredDirectoryHasEmptySubDirectory()

	log.Info("pre-validation step completed successfully")
}

/**

This func should be always run before ZipFileValidation.go's validate() function!

It's the acid test for proofing the customer's zip file

The structure of the test archive 'test-archive.zip' is:

   test-archive.zip
   |
   | - file_1.txt   (non-empty file)
   | - file_2.txt   (non-empty file)
   | - file_3.txt   (empty file)
   |
   | - a
       | - file_7.txt   (non-empty file)
   	   | - file_8.txt   (non-empty file)
       | - file_9.txt   (empty file)
   |
   | - b
       |- d
          | - file_4.txt   (non-empty file)
   	      | - file_5.txt   (non-empty file)
          | - file_6.txt   (empty file)
   |
   | - c   (empty directory)
   | - e
       |- f (empty directory)

it doesn't matter whether files or directories are empty in the zip but the files or directories declared in the required properties cannot be empty!
for example the following property specification will not raise a panic:
	requiredDirs=a/,b/
	requiredFiles=file_1.txt,file_2.txt

for example the following property specification will raise a panic (empty required directory):
	requiredDirs=a/,c/
	requiredFiles=file_1.txt,file_2.txt

for example the following property specification will raise a panic (empty required file):
	requiredDirs=a/,b/
	requiredFiles=file_1.txt,file_3.txt
*/
func validZipNoPanicExpected() {
	defer func() {
		if err := recover(); err != nil {
			errorMessage := fmt.Sprintf("unexpected panic validating good file, panic: %v", err)
			panic(errorMessage)
		}
		log.Info("validZipNoPanicExpected(): zip validation went as expected for a valid zip file!")
	}()

	ValidateZip(testZipPath, validRequiredFiles, validRequiredDirs)
}

func invalidZipMissingRequiredFile() {
	defer func() {
		if err := recover(); err != nil {
			log.Info("zip validation failed as expected on a missing required file!")
		} else {
			errorMessage := "zip validation should have failed on a missing required file!"
			panic(errorMessage)
		}
	}()

	requiredFiles := "file_1.txt,file_2.txt,missing.txt"
	ValidateZip(testZipPath, requiredFiles, validRequiredDirs)
}

func invalidZipMissingRequiredDirectory() {
	defer func() {
		if err := recover(); err != nil {
			log.Info("zip validation failed as expected on a missing required directory!")
		} else {
			errorMessage := "zip validation should have failed on a missing required directory!"
			panic(errorMessage)
		}
	}()

	requiredDirs := "a/,b/,missing/"
	ValidateZip(testZipPath, validRequiredDirs, requiredDirs)
}

func invalidZipEmptyRequiredFile() {
	defer func() {
		if err := recover(); err != nil {
			log.Info("zip validation failed as expected on a empty required file!")
		} else {
			errorMessage := "zip validation should have failed on a empty required file!"
			panic(errorMessage)
		}
	}()

	requiredFiles := "file_1.txt,file_2.txt,file_3.txt"
	ValidateZip(testZipPath, requiredFiles, validRequiredDirs)
}

func invalidZipEmptyRequiredDirectory() {
	defer func() {
		if err := recover(); err != nil {
			log.Info("zip validation failed as expected on a empty required directory!")
		} else {
			errorMessage := "zip validation should have failed on a empty required directory!"
			panic(errorMessage)
		}
	}()

	requiredDirs := "a/,b/,c/"
	ValidateZip(testZipPath, validRequiredFiles, requiredDirs)
}

func invalidZipRequiredDirectoryHasEmptySubDirectory() {
	defer func() {
		if err := recover(); err != nil {
			log.Info("zip validation failed as expected on a required directory with empty sub-directory!")
		} else {
			errorMessage := "zip validation should have failed on a required directory with empty sub-directory"
			panic(errorMessage)
		}
	}()

	requiredDirs := "a/,b/,e/"
	ValidateZip(testZipPath, validRequiredFiles, requiredDirs)
}
