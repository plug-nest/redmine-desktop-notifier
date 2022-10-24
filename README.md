# Get it work

## Modify redmine.service
```bash
# WorkingDirectory=YOUR_DIRECTORY_YOU_STORED_THIS_PROJECT
# ExecStart=YOUR_DIRECTORY_YOU_STORED_THIS_PROJECT/vN
```

# Run command
```bash
go mod tidy
go build

cp redmine.service ~/.config/systemd/user
systemctl --user enable redmine.service 
systemctl --user start redmine.service 
journalctl --user -u redmine.service
```

# TODO

- store data in sqlite3/postgresql
- <s> store constants in .env </s>
- <s> add crontab </s>
- implement dunst
