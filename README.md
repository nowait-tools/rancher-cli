## Rancher cli

This project is a rancher cli that leverages the rancher api for the various automation we need for rancher.  Rancher currently lacks a sophisticated cli at the moment so this is the solution to tailor one to our needs at Nowait.

## Setup

### Requirements
- Recent version of docker installed.
- Access to [Docker Hub](https://hub.docker.com/u/nowait/dashboard/) and permissions to the nowait organization.
- Access to Rancher with GitHub authorization and access to atleast one environmnet.


### Installation and Configuration

The following instructions are going to assume docker is setup properly.  If you are using docker-machine that means you have setup your environment variables correctly.  On docker for mac, docker for windows, or native it should just need to be running.

- Pull the [nowait/rancher-cli:_tag_](https://hub.docker.com/r/nowait/rancher-cli/) image
- Copy the `dockerfunc.sample` file to `~/.dockerfunc`
- Copy the `secrets.sample` file to `~/.secrets`
- Create an environment API key by opening Rancher, selecting the target environment (e.g. production or stretch) on the menu, navigating to the `API` section, and pressing the `Add Environment API Key` button. **DO NOT** use an account key unless you know what you are doing (it is under the advanced options menu).  Preserve those values in [LastPass](https://www.lastpass.com/) as you will never get to see the secret key again.
- Fill out your rancher access key, rancher secret key, docker hub username, docker hub password in your `~/.secrets` file.
- `source` the `~/.dockerfunc` and `~/.secrets` files (and maybe make this part of your `~/.bash_profile` or some other regular process)

## Usage

### Examples

An example of an upgrade where the service has a [sidekick](https://docs.rancher.com/rancher/v1.1/en/cattle/adding-services/#sidekick-services) image, and we're updating the code:

`$ ran_cli_stretch service upgrade --service-like service-name --code-tag "0.10.1" --wait`

An example of an upgrade where the service does not have a sidekick is below. Notice that the tag specification is `--runtime-tag`:

`$ ran_cli_stretch service upgrade --service-like host-cleanup --runtime-tag "1.0" --wait`

### Details

#### Subcommands

The `service` command has 2 subcommands: `upgrade` and `upgrade-finish`. `upgrade-finish` is for when you upgrade a service but don't fully finish the upgrade. A sample upgrade-finish is show below

`ran_cli_stretch service upgrade-finish --service Nowait-Server-Consumer-Api`

#### Options for the `upgrade` command.

- `--service Service-Name` - Name of the service you would like to upgrade

- `--service-like Prefix-Name-To-Match` - Prefix of the service you would like to upgrade. This option provides for partial matches against available services. For example, if services `Nowait-Server` and `Nowait-Server-Consumer-Api` are defined, `--service-like Nowait-Server` would upgrade both of these.

- `--env-file path/to/.env` - Path to a `.env` file. Will provide validation that the service in Rancher has all the environment variable keys defined in the `.env` file. Note this cli is running inside a container and you must mount your local filesystem in order for the container to see the `.env` file.

- `--env NEW_ENV_KEY=NEW_ENV_VALUE` - Key value pair like `ENV_NAME=ENV_VALUE`. Will add or update the environment variable for the services being upgraded. For multiple environment variables use the following `--env NEW_ENV_1=NEW_ENV_1_VALUE --env NEW_ENV_2=NEW_ENV_2_VALUE`.

- `--runtime-tag nowait/image-name:1.1` - Docker image tag to deploy. Upgrades the main docker image.  The following is also valid `--runtime-tag 1.1` however this assumes that you are still using the same docker image as the service was previously using (in this case nowait/image-name)

- `--code-tag nowait/image-name-code:1.0` - Docker image tag to employ. Upgrades the a sidekick's docker image.  The following is also valid `--code-tag 1.1` however this assumes that you are still using the same docker image as the service was previously using (in this case nowait/image-name-code)

- `--interval [seconds]` - **Experts only**.  The default should be sufficient in most cases and deviations could cause problems. See the documentation for more [information](https://docs.rancher.com/rancher/v1.2/en/cattle/upgrading/#in-service-upgrade).

- `--wait` - No argument value. Upgrade the service and wait until it is finished upgrading. It will then finish the upgrade. If it does not complete within a given timeframe it will rollback the upgrade.

### TODO
- [ ] Provide correct feedback to the user when CATTLE environment variables are not defined.
