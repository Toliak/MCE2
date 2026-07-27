# MakeConfigurationEasier

Quick start:

See https://github.com/Toliak/MCE2/releases/

For Linux AMD64 (x86_64):

```
wget https://github.com/Toliak/MCE2/releases/download/v1.1.2/mce-linux-amd64 -O /tmp/mce
chmod +x /tmp/mce
/tmp/mce
```

To select and install everything (carefully!):
```
/tmp/mce -ALL -no-ui -y
```

For the development information see [DEVELOPMENT.md](DEVELOPMENT.md)

## /etc/skel installation

If you would like to set the configs by default for any newly created user, install them into the /etc/skel directory:

```
/tmp/mce -ALL -home-dir=/etc/skel
```
