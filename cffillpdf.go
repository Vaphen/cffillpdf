package cffillpdf

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"runtime"
)

// Fill inserts all provided values in their corresponding field in the provided pdf file
func Fill(values map[string]string, file *os.File) (filledFile *bytes.Buffer, err error) {
	libraryPath, ok := getLibraryPath()

	if !ok {
		err = errors.New("could not get path of gcf-fillpdf library")
		return
	}

	os.Setenv("PATH", libraryPath+"/pdftk")
	os.Setenv("LD_LIBRARY_PATH", libraryPath+"/pdftk")

	fdfFile, err := createFDFFile(values)
	if err != nil {
		return
	}
	defer os.Remove(fdfFile.Name())

	outputFilePath := "/tmp/filledFile.pdf"
	args := []string{
		file.Name(),
		"fill_form", fdfFile.Name(),
		"output", outputFilePath,
		"flatten",
	}

	pdftkPath, err := exec.LookPath("pdftk")
	if err != nil {
		return
	}

	var stderr bytes.Buffer
	var stdout bytes.Buffer
	cmd := exec.Command(pdftkPath, args...)
	cmd.Stderr = &stderr
	cmd.Stdout = &stdout
	cmd.Dir = "/tmp"
	// Start the command and wait for it to exit.
	if err = cmd.Run(); err != nil {
		if exitError, ok := err.(*exec.ExitError); ok && exitError.ExitCode() != 0 {
			err = fmt.Errorf("error executing pdftk\nstdout: %s\nstderr: %s", stdout.String(), stderr.String())
			return
		}
		return
	}

	defer os.Remove(outputFilePath)
	outputFile, err := os.Open(outputFilePath)
	if err != nil {
		return
	}

	filledFile = &bytes.Buffer{}
	_, err = io.Copy(filledFile, outputFile)
	return
}

func createTemporaryWritableFile(name string) (file *os.File, err error) {
	file, err = ioutil.TempFile("/tmp/", name)
	if err != nil {
		return
	}

	err = file.Chmod(0777)
	return
}

func getLibraryPath() (filename string, ok bool) {
	if _, filename, _, ok = runtime.Caller(0); ok {
		filename = path.Dir(filename)
	}
	return
}

func createFDFFile(values map[string]string) (tmpFile *os.File, err error) {
	fdfBuffer := &bytes.Buffer{}
	w := bufio.NewWriter(fdfBuffer)

	if _, err = fmt.Fprintln(w, fdfHeader); err != nil {
		return
	}

	for key, value := range values {
		if _, err = fmt.Fprintf(w, "<< /T (%s) /V (%v)>>\n", key, value); err != nil {
			return
		}
	}

	if _, err = fmt.Fprintln(w, fdfFooter); err != nil {
		return
	}

	if err = w.Flush(); err != nil {
		return
	}

	tmpFile, err = createTemporaryWritableFile("gcf-pdftk-*.fdf")
	if err != nil {
		return
	}

	_, err = tmpFile.Write(fdfBuffer.Bytes())
	return
}

const fdfHeader = `%FDF-1.2
%,,oe"
1 0 obj
<<
/FDF << /Fields [`

const fdfFooter = `]
>>
>>
endobj
trailer
<<
/Root 1 0 R
>>
%%EOF`
