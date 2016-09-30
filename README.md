## Rancher cli

This project is a rancher cli that leverages the rancher api for the various automation we need for rancher.  Rancher currently lacks a sophisticated cli at the moment so this is the solution to tailor one to our needs at Nowait.

### TODO
- [ ] Make sure service upgrade-runtime and service upgrade-code commands rollback gracefully if upgrade takes too long.
- [ ] Check for valid configuration (all environment variables necessary defined in service) before initiating an upgrade.
- [ ] All service upgrades must use `start before stopping` to ensure that no downtime occurs during service upgrade.
- [ ] Provide correct feedback to the user when CATTLE environment variables are not defined.
