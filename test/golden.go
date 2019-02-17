package test

import (
  "bytes"
  "fmt"
  "io/ioutil"
)

func CompareOutToGolden(outfilename, goldenfilename string) error {
  outcontent, err := ioutil.ReadFile(outfilename)
  if err != nil {
    return fmt.Errorf("error reading back output file %s: %v", outfilename, err)
  }
  goldencontent, err := ioutil.ReadFile(goldenfilename)
  if err != nil {
    return fmt.Errorf("error reading golden file %s: %v", goldenfilename, err)
  }
  if !bytes.Equal(outcontent, goldencontent) {
    return fmt.Errorf("outfile %s does not match golden file %s", outfilename, goldenfilename)
  }
  return nil
}


