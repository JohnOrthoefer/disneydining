## INSTALL

- Make the programs
   - At the build root you can do `make`
   - Everything will be copied into `${buildroot}/bin`
- Copy the executables `cp bin/* /usr/local/bin`
- Install `dining.service`, `dining.timer`, and `diningWeb.service` 
  in `/etc/systemd/system`
- Make sure `dining.service` points to the rigth place
- Setup __nginx__ or some other webserver and configure it to run diningWeb
- DiningWeb needs to point to the templates
- Configure the `config.ini`, and `dining.ini`

