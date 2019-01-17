# GoPull

GoPull is a lightweight REST API written in go that runs a deployment script after receiving a webhook request from Bitbucket, GitHub, or GitLab. All local repositories configured to work with GoPull must be placed in the same directory known as `REPOS_DIR`. In order for this API to access the remote repository, an SSH access key must be added to the remote. A webhook can then be configured to fire on any events supported by the remote (code push by default). When the API receives a request, it checks which branch was affected by the change. If this branch matches the configured branch of the local repository, then the corresponding deployment script is run.

This API features a command line tool for easily adding new repositories or configuring current local repositories. By default, this tool assumes GoPull is installed in `$GOPATH/src/github.com/mattmeyers/gopull` and that all managed repositories are to be placed in `$HOME/repos`.  These paths can be configured with the command line tool itself.

## Installation

GoPull can be installed by running

```
go get -u github.com/mattmeyers/gopull/...
```

This command will install the `gopull` and `gopull-cli` binaries.

## Running the API

Executing the `gopull` binary will start the API as a process. There are a variety of process managers that can be used to handle the REST API. If you are using Ubuntu or another operating system that uses systemd, GoPull can be run as a systemd service. All that is needed is a `gopull.service` configuration file. A simple configuration is as follows.

```
[Unit]
Description=GoPull API Service

[Service]
User=$USER
ExecStart=/absolute/path/to/binary/gopull
Restart=always
RestartSec=10
StandardOutput=syslog
StandardError=syslog
SyslogIdentifier=gopull-api-service

[Install]
WantedBy=multi-user.target
```

After creating this file, the service can be enabled and started.

```
systemctl enable /absolute/path/to/gopull.service
systemctl start gopull.service
```

Then check that the service started properly with `systemctl status gopull`.

By default, GoPull runs on port 8080. A webserver such as Apache or Nginx can be used as a reverse proxy to avoid called the port directly. A very basic Nginx configuration without SSL is as follows

```
server {
  server_name <YOUR_DOMAIN>;

  location / {
      proxy_pass http://localhost:8080/;
  }
}
```

## Configuration

GoPull managed repository configurations are stored in a JSON file named `repos.json` that is located in the GoPull directory (referred to as `GOPULL_DIR` going forward). An entry in this file contains all of the information that GoPull needs to react properly to webhook requests. The configuration file is a map with the repository's full name (`user/repo_name`) as the keys. An example entry is as follows.

```json
"mattmeyers/gopull": {
    "user": "mattmeyers",
    "name": "gopull",
    "fullName": "mattmeyers/gopull",
    "branch": "master",
    "path": "$HOME/repos/mattmeyers/gopull",
    "deploymentScript": "$GOPATH/src/github.com/mattmeyers/gopull/deployment_scripts/gopull_deploy.sh"
}
```

Note that an entry does not store information about the remote at this time. Therefore, two repositories with the same fullname but different remotes cannot currently be managed by GoPull.

After adding an entry to the configuration file, a deployment script must be made. By default, deployment scripts go in `GOPULL_DIR/deployment_scripts`. There is an example script in this directoey named `deploy.src.sh`. This script changes the working directory to that of the managed repository, then runs the minimum number of required commands, namely `git fetch` and `git merge`. This file can be copied and renamed. By convention, the file is named `<repo_name>_deploy.sh`. Note that the commands run in the deployment script are executed by whoever ran GoPull. If any `sudo`commands are run from the script, a password might be required and the script will fail. Individual commands can be added as `NOPASSWD` fields in `/etc/sudoers`.

Next, an ssh key must be added to the remote repository. Depending on your remote, follow their guide for this.

- Bitbucket: https://confluence.atlassian.com/bitbucket/access-keys-294486051.html
- Github: https://help.github.com/articles/adding-a-new-ssh-key-to-your-github-account/
- Gitlab: https://docs.gitlab.com/ee/ssh/

Next, configure the webhook in the remote's settings. Again follow the guides for each remote.

- Bitbucket: https://confluence.atlassian.com/bitbucket/manage-webhooks-735643732.html
- Github: https://developer.github.com/webhooks/
- Gitlab: https://docs.gitlab.com/ee/user/project/integrations/webhooks.html

Use the following GoPull REST endpoints when configuring the webhooks

- Bitbucket: /webhooks/bitbucket
- Github: /webhooks/github
- Gitlab: /webhooks/gitlab

Finally, the repository can be cloned using ssh. By default, repositories are cloned into `$HOME/repos`. Because GoPull only cares about a single branch, it's a good idea to only clone that branch using

```
git clone --single-branch --branch <BRANCH> \
    git@<REMOTE>:<USER>/<REPOSITORY>
```

## gopull-cli

GoPull ships with a command line tool that makes setting up and configuring repositories much easier.

```
NAME:
   gopull-cli - Configure the GoPull REST API to pull remote repo changes

USAGE:
   gopull-cli [global options] command [command options] [arguments...]

VERSION:
   0.0.1

COMMANDS:
     config   Configure the GoPull environment
     list     List configured local repos
     add      Add a new repository
     edit     Edit an existing repository
     delete   Delete a repository
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version
```

### Configuring the GoPull Environment

Variables used by GoPull can be set using `gopull-cli config`.  Passing no flags will simply print the currently configured values.  Passing the `--repos-dir` or `-r` with a value will set the `REPOS_DIR` path.

### Adding a New Repository

In order to initialize a local managed repository, follow the instructions above for adding your ssh key to the remote as well as configuring the webhook.  Then use the command

```
gopull-cli add --uri git@<REMOTE>:<USER>/<REPOSITORY> --branch <BRANCH>
```

This will clone the repository into `REPOS_DIR/<USER>/<REPOSITORY>`, add the configuration entry to `repos.json`, and copy the deployment script template.

## Release History

* 0.0.1
    * Work in progress