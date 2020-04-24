# cffillpdf
cffillpdf is a small library built around the pdftk utility. It was designed to work well with google cloud functions but can also be used in other applications.

#Sample

```
package main

import (
	"log"

	"github.com/Vaphen/cffillpdf"
)

func main() {
  // Create the form values.
  values := map[string]string{
    "field_1": "Hello",
    "field_2": "World",
  }

  // pdf should contain the input file
  pdf := []byte{}

  // Fill the form PDF with our values.
  // The returned bytes.Buffer contains the filled PDF.
  filled, err := cffillpdf.Fill(values, pdf)
  if err != nil {
    log.Fatal(err)
  }

  // The buffer can saved as a new file
  err = ioutil.WriteFile("filled.pdf", filled.Bytes(), 0666)
}
```

# Running on PCF
*The binary provided in this repository is compiled for the execution in google cloud functions only.* You may have to change the pdftk binary if you want to run this library under any other system.
In order to run the library in google cloud functions, you must put a copy of the pdftk folder into the cffillpdf's vendor directory of your function:
**{your_function_directory}/vendor/github.com/Vaphen/cffillpdf**

Unfortunately, this must be done after each execution of *go mod vendor*.

# Running on other systems
You can install the library on other systems by changing the content of the pdftk directory with a proper binary of pdftk for your system.





