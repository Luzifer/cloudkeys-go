[![Go Report Card](https://goreportcard.com/badge/github.com/Luzifer/cloudkeys-go)](https://goreportcard.com/report/github.com/Luzifer/vault2env)
![](https://badges.fyi/github/license/Luzifer/cloudkeys-go)
![](https://badges.fyi/github/downloads/Luzifer/cloudkeys-go)
![](https://badges.fyi/github/latest-release/Luzifer/cloudkeys-go)

# CloudKeys Go

CloudKeys Go is a port of the [former CloudKeys project](https://github.com/awesomecoders/cloudkeys) written in PHP. This Go port is 100% compatible to the data written by the older version but adds more storage options. Also it's faster, more stable and last but not least it's not longer based on PHP but working as a tiny webserver implemented in pure Go.

## Usage

```
# cloudkeys-go --help
Usage of ./cloudkeys-go:
      --cookie-authkey="": Key used to authenticate the session
      --cookie-encryptkey="": Key used to encrypt the session
      --listen=":3000": IP and port to listen on
      --password-salt="": A random unique salt for encrypting the passwords
      --storage="local:///./data": Configuration for storage adapter (see README.md)
      --username-salt="": A random unique salt for encrypting the usernames
```

What you definitely should set when starting the server:

- `cookie-authkey` - This flag protects the encrypted cookies you're putting on the users computers containing the session. If you don't set it yourself it will be randomly generated. In that case your users will get logged out every time you restart the server. You need to use a key with length of 16 characters (AES128) or 32 characters (AES256).
- `cookie-encryptkey` - This flag is the encryption key itself. Like the authkey it will get autogenerated with the same result. You need to use a key with length of 16 characters (AES128) or 32 characters (AES256).
- `password-salt` - [deprecated] In version <=v1.6.1 the password was hashed with a static salt. You only need to provide this if you started using Cloudkeys in one of those versions.
- `username-salt` - The usernames are the keys in the database. They are hashed but you can put an additional salt to them to make it harder to decipher them.

If you don't want to define the secrets using command line flags you also can use environment variables to set those flags:

```
FLAG                ENV-Variable

password-salt       passwordSalt
username-salt       usernameSalt
storage             storage
listen              listen
cookie-authkey      authkey
cookie-encryptkey   encryptkey
```

## Supported storage engines

### Local file storage (default)

This storage engine is used in the default config when you just start up the server as you can see in the output above. You don't have many options to set for this one. The only thing is the path where all the data is stored.

```
Schema:  local:///<your data directory>
Example: local:///./data
```

The directory can be set absolute or relative. Please ensure there are **3** slashes between `local:` and the begin of your path. (So if you're setting an absolute path you will set 4 slashes in a row.)

### Amazon Web Services S3

This is the storage engine you want to use if you're migrating from the old CloudKeys version. This option is fully compatible to every piece of data the old version stored.

```
Schema:  s3://<bucket><path>
Example: s3://mybucket/
```

You can specify the bucket and also a prefix for the storage. That way you even could use one bucket for different instances of CloudKeys Go. In case you're migrating from the old version you need to set the path to `/`.

For this to work you also need to set three environment variables: `AWS_ACCESS_KEY`, `AWS_SECRET_ACCESS_KEY` and `AWS_REGION`. When its about `AWS_REGION` pay attention to select the right region for your bucket.

### Redis

If you want to utilize a Redis storage server or even a Redis cluster you can choose this storage type. Authentication is supported as well as selecting the database to use. Aditionally you can set a prefix for the keys.

```
Schema:  redis+tcp://auth:<password>@127.0.0.1:6379/<db>?timeout=10s&maxidle=1&prefix=<prefix>
Example: redis+tcp://auth:mypass@redis.example.com:6379/5?prefix=cloudkeys::
```

## Install on Heroku

1. Create a new Heroku app

    ```
    # heroku create -b https://github.com/heroku/heroku-buildpack-go
    ```

2. Push the code to your app

    ```
    # git push heroku master
    ```

3. Set your configuration variables in the Heroku apps dashboard (see env variables in usage section above)
