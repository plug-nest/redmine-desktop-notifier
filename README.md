```bash
cp redmine.service ~/.config/systemd/user
systemctl --user enable redmine.service 
systemctl --user start redmine.service 
journalctl --user -u redmine.service
```

# TODO

- store data in sqlite3/postgresql
- <s> store constants in .env </s>
- <s> add crontab </s>
