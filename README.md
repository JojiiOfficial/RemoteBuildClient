# RemoteBuild - Client
This is the client for the [RemoteBuild system](https://github.com/RemoteBuild/Remotebuild). You can use it to create/control jobs running on the server.<br>
If you have no idea what this is, read the README from the repository I just linked.

# Install
Make sure you have go >=1.13 installed

## Arch linux
```bash
yay -S rbuild-cli-git
```

or get it compiled from [my repo](https://repo.jojii.de)

## Compile manually
```bash
make build
sudo make install
# Run sudo make uninstall to uninstall it
```

# Setup
```bash
rbuild setup <ServerHost>
```

### Create an account<br>
(note: you have to set `allowregistration` to true inside the server config!)

```bash
rbuild register
```

For further documentation refer to the man page
