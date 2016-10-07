## Test Plan

The following is a manual test plan to be followed in order to ensure that all functionality works correctly.  The current test coverage for this project is assuming the mocks work as expected, in most situations that will work great but this document aims to provide a list of all current functionality so all features can be tested to avoid regressions.

### Service Command
#### upgrade
- upgrading multiple services
  - command flags
    - service-like
    - env-file
      - Must verify that it will prevent an upgrade if an environment variable exists in the env-file but not in the Rancher service.  Omitting this will skip validation of the config.
    - service
    - code-tag
    - runtime-tag
    - wait
      - When using wait flag you must be able to show that it will automatically complete an upgrade that should succeed and will rollback a failed upgrade. To simulate a failed upgrade try to upgrade to an image tag that does not exist.
- upgrading single service
  - with and without env-file validation

#### finish-upgrade
- Upgrade a service manually through the Rancher UI
- Run `rancher-cli service finish-upgrade --service Service-Name` replacing Service-Name with your service's name

