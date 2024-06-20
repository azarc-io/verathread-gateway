# chart-publish-action

A GitHub action to publish a Helm chart to a directory.

## Usage

The idea is to define a Helm chart in the same repository where your micro-service code is and use this action to publish the chart to another dedicated repository with GitHub Pages support.

To use this action for publishing a Helm chart to another repository, you need to use some other actions for doing `git clone/checkout` and `git commit/push`. In the example below actions `actions/checkout` and `EndBug/add-and-commit` are used for those purposes.

### Prerequisites

1. A GitHub repo containing a directory with a Helm chart (e.g. 'chart')
1. A GitHub repo used for publishing Helm charts via GitHub Pages.
1. A GitHub access token for pushing the changes.

### Example

Assumptions:

* The Helm chart directory in the current repository is 'chart'.
* The GitHub Pages repository for hosting published Helm charts is 'janedoe/helm-charts'
* The GitHub Pages branch is 'gh-pages'.
* The GitHub Pages directory is the root directory of 'janedoe/helm-charts' repository.
* The personal access token for accessing the Helm charts repository is in 'CHART_PUBLISH_TOKEN' secret.

```yaml
name: ci

on:
  push:

jobs:
  buildx:
    runs-on: ubuntu-latest
    steps:

      - name: Checkout helm-charts
        uses: actions/checkout@v2
        with:
          repository: janedoe/helm-charts
          ref: refs/heads/gh-pages
          path: helm-charts
          token: ${{ secrets.CHART_PUBLISH_TOKEN }}

      - name: Publish the Helm chart
        id: publish-chart
        uses: moikot/chart-publish-action@v1
        with:
          charts_dir: "helm-charts"
          charts_url: "https://janedoe.github.io/helm-charts"

      - name: Commit and push helm-charts
        uses: EndBug/add-and-commit@v6
        with:
          author_name: Jane Doe
          author_email: janedoe@example.com
          branch: gh-pages
          cwd: helm-charts
          message: "Commit chart ${{ steps.publish-chart.outputs.chart }}"
```

## Customization

### Inputs

Following inputs can be used as `step.with` keys:

| Name               | Type    | Description                       |
|--------------------|---------|-----------------------------------|
| `chart`          | String  | The relative path under $GITHUB_WORKSPACE to the chart (default: 'chart') |
| `helm_image`           | String  | The Docker image containing Helm (default: 'alpine/helm:3.3.0') |
| `charts_dir`      | String     | The relative path under $GITHUB_WORKSPACE to publish to |
| `charts_url`  | String  | The URL of the charts repository (e.g. 'https://janedoe.github.io/helm-charts') |

### Outputs

Following outputs are available:

| Name          | Type    | Description                           |
|---------------|---------|---------------------------------------|
| `chart`        | String  | The published chart file name |

## Limitation

This action is available for Linux [virtual environments](https://docs.github.com/en/actions/reference/virtual-environments-for-github-hosted-runners#supported-virtual-environments-and-hardware-resources) only.