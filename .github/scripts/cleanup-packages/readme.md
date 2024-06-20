## Clean up helper script

Cleans up docker images that were using during the development and qa cycle, should only be used at the end of a workflow
when a pr has been released and relevant environments have been updated

## Usage example

```yaml
- name: "Clean Up Container Registry"
  run: ./.github/scripts/cleanup-packages/cleanup
  env:
    GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
    PACKAGES: verathread-gateway/gateway-fe,verathread-gateway/gateway-be
    TICKET: VTHP-1234
```
