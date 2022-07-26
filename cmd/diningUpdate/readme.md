Two URLs needed-
Avilablity - https://disneyworld.disney.go.com/finder/api/v1/explorer-service/dining-availability-list/false/wdw/80007798;entityType=destination/2022-08-15/7/?mealPeriod=80000714

curl --user-agent 'Chrome/102.0.0.0' 'https://disneyworld.disney.go.com/finder/api/v1/explorer-service/dining-availability-list/false/wdw/80007798;entityType=destination/2022-08-15/7/?mealPeriod=80000714'

mealPeriod = 
   Breakfast option-80000712
   Brunch    option-80000713
   Dinner    option-80000714
   Lunch     option-80000717
--OR--
searchTime = HH:MM:00 URL encoded 
   eg 20:30:00 is 20%3A30%3A00


Restuarnt List - https://disneyworld.disney.go.com/finder/api/v1/explorer-service/list-ancestor-entities/wdw/80007798;entityType=destination/2022-06-18/dining

