package test

import (
  "bufio"
  "bytes"
  "fmt"
  "io/ioutil"
  "os"
  "testing"
)

type Data struct {
  TplFilePath string;
  OutFilePath string;
  GoldenFilePath string;
  OutW *bufio.Writer;
  OutF *os.File;
}

func Setup(t *testing.T, basename string) *Data {
  t.Helper()
  tplfilepath := "testdata/" + basename + ".tpl"
  outfilepath := "testdata/" + basename + ".out"
  goldenfilepath := "testdata/" + basename + ".golden"

  os.Remove(outfilepath)
  f, err := os.Create(outfilepath)
  if err != nil {
    t.Fatal(err)
  }
  w := bufio.NewWriter(f)

  return &Data{
    TplFilePath: tplfilepath,
    OutFilePath: outfilepath,
    GoldenFilePath: goldenfilepath,
    OutF: f,
    OutW: w,
  }
}

func Finish(t *testing.T, data *Data) {
  data.OutW.Flush()
  data.OutF.Close()
  if err := CompareOutToGolden(data.OutFilePath, data.GoldenFilePath); err != nil {
    t.Fatal(err)
  }
}

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
