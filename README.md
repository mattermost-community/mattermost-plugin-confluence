# Mattermost Plugin for Confluence 

[![Build Status](https://img.shields.io/circleci/project/github/mattermost/mattermost-plugin-confluence/master)](https://circleci.com/gh/mattermost/mattermost-plugin-confluence)
[![Code Coverage](https://img.shields.io/codecov/c/github/mattermost/mattermost-plugin-confluence/master)](https://codecov.io/gh/mattermost/mattermost-plugin-confluence)
[![Release](https://img.shields.io/github/v/release/mattermost/mattermost-plugin-confluence)](https://github.com/mattermost/mattermost-plugin-confluence/releases/latest)
[![HW](https://img.shields.io/github/issues/mattermost/mattermost-plugin-confluence/Up%20For%20Grabs?color=dark%20green&label=Help%20Wanted)](https://github.com/mattermost/mattermost-plugin-confluence/issues?q=is%3Aissue+is%3Aopen+sort%3Aupdated-desc+label%3A%22Up+For+Grabs%22+label%3A%22Help+Wanted%22)

A Mattermost plugin for Confluence. Supports Confluence Cloud, Server and Data Center versions. This plugin helps your teams collaborate and keep in sync as Confluence Pages and Spaces get updated.  Comments and activity can be pushed into specific Mattermost channels for full visibility. 

# Documentation 

Installation and Usage instructions are located here: https://mattermost.gitbook.io/plugin-confluence/

# Development 

### Maintainers 

**Maintainer:** [@jfrerich](https://github.com/jfrerich)
**Co-Maintainer:** [@levb](https://github.com/levb)

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
