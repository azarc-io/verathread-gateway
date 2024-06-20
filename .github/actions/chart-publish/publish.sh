#!/usr/bin/env bash

# Copyright (c) 2021 Sergey Anisimov
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     https://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

set -o errexit
set -o pipefail

DEFAULT_HELM_IMAGE="alpine/helm:3.4.2"

main() {
    local helm_image="${DEFAULT_HELM_IMAGE}"
    local chart="chart"
    local charts_dir=
    local charts_url=

    parse_command_line "$@"

    chart="${GITHUB_WORKSPACE}/${chart}"
    charts_dir="${GITHUB_WORKSPACE}/${charts_dir}"

    echo "Updating chart dependencies ..."
    docker run -i --rm \
      -v "${chart}":/chart \
      "${helm_image}" dependency update /chart

    echo "Linting the chart ..."
    docker run -i --rm \
      -v "${chart}":/chart \
      "${helm_image}" lint /chart

    echo "Packaging the chart ..."
    local temp_dir=$(mktemp -d)
    docker run -i --rm \
      -v "${chart}":/chart \
      -v "${temp_dir}":/temp \
      "${helm_image}" package /chart --destination /temp

    if [[ -f "${charts_dir}/index.yaml" ]]; then
      echo "Merging the chart index ..."
      docker run -i --rm \
        -v "${temp_dir}":/temp \
        -v "${charts_dir}":/charts \
        "${helm_image}" repo index /temp --url ${charts_url} --merge /charts/index.yaml
    else
      echo "Creating the chart index ..."
      docker run -i --rm \
        -v "${temp_dir}":/temp \
        -v "${charts_dir}":/charts \
        "${helm_image}" repo index /temp --url ${charts_url}
    fi

    # Set the publishied chart variable
#    echo ::set-output name=chart::$(find ${temp_dir}/*.tgz -printf "%f")
    CHART_OUT=$(find ${temp_dir}/*.tgz -printf "%f")
    echo "chart=$CHART_OUT" >> $GITHUB_OUTPUT

    mv -f "${temp_dir}"/*.tgz -f "${temp_dir}"/index.yaml "${charts_dir}"
}

show_help() {
cat << EOF
Usage: $(basename "$0") <options>
    -h, --help               Display help
    -i, --helm-image         The helm Docker image (default: $DEFAULT_HELM_IMAGE)"
    -c, --chart              The chart directory (default: chart)
    -d, --charts-dir         The directory to publish to
    -u, --charts-url         The URL of the charts repository
EOF
}

parse_command_line() {
    while :; do
        case "${1:-}" in
            -h|--help)
                show_help
                exit
                ;;
            -i|--helm-image)
                if [[ -n "${2:-}" ]]; then
                    helm_image="$2"
                    shift
                else
                    echo "ERROR: '-i|--helm-image' cannot be empty." >&2
                    show_help
                    exit 1
                fi
                ;;
            -c|--chart)
                if [[ -n "${2:-}" ]]; then
                    chart="$2"
                    shift
                else
                    echo "ERROR: '-c|--chart' cannot be empty." >&2
                    show_help
                    exit 1
                fi
                ;;
            -d|--charts-dir)
                if [[ -n "${2:-}" ]]; then
                    charts_dir="$2"
                    shift
                else
                    echo "ERROR: '-d|--charts-dir' cannot be empty." >&2
                    show_help
                    exit 1
                fi
                ;;
            -u|--charts-url)
                if [[ -n "${2:-}" ]]; then
                    charts_url="$2"
                    shift
                else
                    echo "ERROR: '-u|--charts-url' cannot be empty." >&2
                    show_help
                    exit 1
                fi
                ;;
            *)
                break
                ;;
        esac

        shift
    done

    if [[ -z "$charts_dir" ]]; then
        echo "ERROR: '-d|--charts-dir' is required." >&2
        show_help
        exit 1
    fi

    if [[ -z "$charts_url" ]]; then
        echo "ERROR: '-u|--charts-url' is required." >&2
        show_help
        exit 1
    fi
}


main "$@"
