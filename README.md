# Mattermost Plugin Confluence

[![CircleCI branch](https://img.shields.io/circleci/project/github/mattermost/mattermost-plugin-confluence/master.svg)](https://circleci.com/gh/mattermost/mattermost-plugin-confluence)


**Maintainer:** [@jfrerich](https://github.com/jfrerich)
**Co-Maintainer:** [@levb](https://github.com/levb)

A Mattermost plugin for Confluence. Supports Confluence Cloud, Server and Data Center versions.

## Installation and setup

### Platform & tools

- Make sure you have following components installed:

  - Go - v1.13 - [Getting Started](https://golang.org/doc/install)
    > **Note:** If you have installed Go to a custom location, make sure the `$GOROOT` variable is set properly. Refer [Installing to a custom location](https://golang.org/doc/install#install).
  - Make

## Setting up CircleCI

Set up CircleCI to run the build job for each branch and build-and-release for each tag.

1. Go to [CircleCI Dashboard](https://circleci.com/dashboard).
2. In the top left, you will find the Org switcher. Select your Organisation.
3. If this is your first project on CircleCI, go to the Projects page, click the **Add Projects** button, then click the **Set Up Project** button next to your project. You may also click **Start Building** to manually trigger your first build.
4. To manage GitHub releases using CircleCI, you need to add your github personal access token to your project's environment variables.
   - Follow the instructions [here](https://help.github.com/en/articles/creating-a-personal-access-token-for-the-command-line) to create a personal access token. For CircleCI releases, you would need the `repo` scope.
   - Add the environment variable to your project as `GITHUB_TOKEN` by following the instructions [here](https://circleci.com/docs/2.0/env-vars/#setting-an-environment-variable-in-a-project).

## Installation

1. Go to the [releases page of this GitHub repository](https://github.com/mattermost/mattermost-plugin-confluence/releases/latest) and download the latest release for your Mattermost server.
2. Upload this file in the Mattermost **System Console > Plugins > Management** page to install the plugin. To learn more about how to upload a plugin, [see the documentation](https://docs.mattermost.com/administration/plugins.html#plugin-uploads).
3. You can configure the Plugin from **System Console > Plugins > Confluence**.

## Configuration

- **Secret**: The generated Webhook Secret for use in Conlfuence Webhook Configuration.
