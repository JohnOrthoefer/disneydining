package squelch

import (
   "encoding/json"
   "os"
   "log"
   "time"
)

type Squelch struct {
   Items []interface{}
   Updated time.Time
}

func New() *Squelch {
   return new(Squelch)
}


func (s *Squelch)Len()int {
   return len(s.Items)
}

func (s *Squelch)Add(item interface{}) {
   s.Items = append(s.Items, item)
   s.Updated = time.Now()
   log.Printf("Added new Item: %T (%q)", item, item)
}

func (s *Squelch)Mute(item interface{}, cmp func(interface{}, interface{})bool) bool {
   for _, i := range s.Items {
      if cmp(i, item) { 
         return true 
      }
   }
   return false
}

func (s *Squelch)Load(name string) error{
   j, err := os.ReadFile(name)
   if err != nil {
      return err
   }
   json.Unmarshal(j, s)
   return nil
}

func (s *Squelch)Save(name string) {
   data, _ := json.MarshalIndent(s, "", " ")
   os.WriteFile(name, data, 0644)
}

