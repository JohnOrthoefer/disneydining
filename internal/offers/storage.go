package offers

import (
	"encoding/json"
   "compress/zlib"
   "compress/gzip"
   "os"
   "lockedfile"
   "io"
   "log"
   "path"
)

func (d DiningMap) SaveOffers(n string) {

   f, err := lockedfile.OpenFile(n, os.O_WRONLY|os.O_CREATE, 0644)
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
      f.Flush()
   }
}

func (d DiningMap) LoadOffers(n string) {
   f, err := lockedfile.OpenFile(n, os.O_RDONLY, 0644)
   defer f.Close()

   if os.IsNotExist(err) {
      d = NewOffers()
      return
   }
   checkErr(err)

   switch path.Ext(n) {
   case ".gzip":
      r, err := gzip.NewReader(f)
      checkErr(err)
      j, err := io.ReadAll(r)
      r.Close()
      if err != nil {
         log.Printf("in LoadOffers - %s\n", err)
         d = NewOffers()
      }
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
}

// vim: noai:ts=3:sw=3:set expandtab:

