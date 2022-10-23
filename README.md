```bash
cp redmine.service ~/.config/systemd/user
systemctl --user enable redmine.service 
systemctl --user start redmine.service 
journalctl --user -u redmine.service
```