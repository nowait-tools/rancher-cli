## Rancher cli

This project is a rancher cli that leverages the rancher api for the various automation we need for rancher.  Rancher currently lacks a sophisticated cli at the moment so this is the solution to tailor one to our needs at Nowait.

## Usage

Update with the correct instructions.

### TODO
- [ ] Make sure service upgrade-runtime command rollback gracefully if upgrade takes too long.
- [ ] Check for valid configuration (all environment variables necessary defined in service) before initiating an upgrade.
- [ ] Provide correct feedback to the user when CATTLE environment variables are not defined.
- [ ] Provide validation on the image name and tag.  Currently expects full image name. Should probably just be the tag name like the --tag implies.
- [ ] What should the batch size and interval be when upgraded?

