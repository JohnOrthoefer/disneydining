package offers

import (
	"encoding/json"
   "compress/zlib"
   "compress/gzip"
   "os"
   "io"
//   "log"
   "path"
)

func (d DiningMap) SaveOffers(n string) {

   f, err := os.OpenFile(n+".new", os.O_WRONLY|os.O_CREATE, 0644)
   checkErr(err)
   defer f.Close()

	data, _ := json.Marshal(d)

   switch path.Ext(n) {
   case ".gzip":
      w := gzip.NewWriter(f)
      w.Write(data)
      w.Flush()
      w.Close()
   case ".zz":
      w := zlib.NewWriter(f)
      w.Write(data)
      w.Flush()
      w.Close()
   default:
      f.Write(data)
   }
   err = os.Rename(n, n+".bak")
   checkErr(err)
   err = os.Rename(n+".new", n)
   checkErr(err)
}

func (d DiningMap) LoadOffers(n string) {
   f, err := os.Open(n)

   if os.IsNotExist(err) {
      d = NewOffers()
      return
   }
   checkErr(err)

//   log.Printf("Seeking...")
//   f.Seek(0, 0)

   switch path.Ext(n) {
   case ".gzip":
      r, err := gzip.NewReader(f)
      checkErr(err)
      j, _ := io.ReadAll(r)
/*      if err != nil {
         log.Printf("len(j) = %d", len(j))
         log.Printf("in LoadOffers - %s\n", err)
         d = NewOffers()
      }
*/
      r.Close()
      json.Unmarshal(j, &d)
   case ".zz":
      r, err := zlib.NewReader(f)
      checkErr(err)
      j, err := io.ReadAll(r)
      r.Close()
      checkErr(err)
	   json.Unmarshal(j, &d)
   default:
      j, err := io.ReadAll(f)
      checkErr(err)
	   json.Unmarshal(j, &d)
   }
   f.Close()
}

// vim: noai:ts=3:sw=3:set expandtab:

