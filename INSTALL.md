## INSTALL

- Make the programs
- Move diningUpdate to `/usr/local/bin`
- Install `dining.service`, `dining.timer`, and `diningWeb.service` 
  in `/etc/systemd/system`
- Make sure `dining.service` points to the rigth place
- Setup =nginx= or some other webserver and configure it to run diningWeb
- DiningWeb needs to point to the templates
- Configure the `config.ini`, and `dining.ini`

