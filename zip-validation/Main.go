// COPYRIGHT Ericsson 2023

// The copyright to the computer program(s) herein is the property of
// Ericsson Inc. The programs may be used and/or copied only with written
// permission from Ericsson Inc. or in accordance with the terms and
// conditions stipulated in the agreement/contract under which the
// program(s) have been supplied.

package main

import (
	"eric-oss-hello-world-go-app/zip-validation/validate"

	"flag"
	"log"
)

const testZipPath string = "zip-validation/test-archive.zip"

func main() {
	log.Println("Zip file validation...")

	// first run the sanity checks on the code itself...
	validate.PreValidation()

	// bob vars and lists don't play well together so need to split the required dirs & files variable values
	var zipPath string
	var requiredDirs string
	var requiredFiles string

	// pull in the command line vars passed in the bob rule
	flag.StringVar(&requiredDirs, "requiredDirs", "NONE", "Non-empty dirs required to be in the zip")
	flag.StringVar(&requiredFiles, "requiredFiles", "NONE", "Files required to be in the zip")
	flag.StringVar(&zipPath, "zipPath", "NONE", "Path to the zip to be evaluated")
	flag.Parse()

	validate.ValidateZip(zipPath, requiredFiles, requiredDirs)
}
