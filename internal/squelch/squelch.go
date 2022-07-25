package squelch

import (
   "log"
   "encoding/json"
   "io/ioutil"
)

type Squelch struct {
   Items []interface{}
}

func NewSquelch() *Squelch {
   return new(Squelch)
}

func (s *Squelch)Add(item interface{}) {
   s.Items = append(s.Items, item)
}

func (s *Squelch)Mute(item interface{}) bool {
   for _, i := range s.Items {
      log.Printf("Compare %q %q", i, item)
      if i == item { return true }
   }
   return false
}

func (s *Squelch)Load(name string) error{
   j, err := ioutil.ReadFile(name)
   if err != nil {
      return err
   }
   json.Unmarshal(j, s)
   return nil
}

func (s *Squelch)Save(name string) {
   data, _ := json.MarshalIndent(s, "", " ")
   ioutil.WriteFile(name, data, 0644)
}

