## Contents

- [Overview](#overview)
- [Features](#features)
- [Admin Guide](/docs/admin-guide.md)
  - [Prerequisites](#prerequisites)
  - [Installation](#installation)
    - [Marketplace Installation (Recommended)](#marketplace-installation-recommended)
    - [Manual Installation](#manual-installation)
  - [Install on Confluence](#install-on-confluence)
    - [Set up Confluence Server](#set-up-confluence-server)
    - [Set up Confluence Cloud](#set-up-confluence-cloud)
- [Configure Notifications](#configure-notifications)
  - [/confluence subscribe](#confluence-subscribe)
  - [/confluence list](#confluence-list)
  - [/confluence edit](#confluence-edit)
  - [/confluence unsubscribe](#confluence-unsubscribe)
- [Development](#development)
  - [Maintainers](#maintainers)
  - [Platform and Tools](#platform-and-tools)
  - [Set up CircleCI](#set-up-circleci)

# Mattermost Plugin for Confluence 

[![Build Status](https://img.shields.io/circleci/project/github/mattermost/mattermost-plugin-confluence/master)](https://circleci.com/gh/mattermost/mattermost-plugin-confluence)
[![Code Coverage](https://img.shields.io/codecov/c/github/mattermost/mattermost-plugin-confluence/master)](https://codecov.io/gh/mattermost/mattermost-plugin-confluence)
[![Release](https://img.shields.io/github/v/release/mattermost/mattermost-plugin-confluence)](https://github.com/mattermost/mattermost-plugin-confluence/releases/latest)
[![HW](https://img.shields.io/github/issues/mattermost/mattermost-plugin-confluence/Up%20For%20Grabs?color=dark%20green&label=Help%20Wanted)](https://github.com/mattermost/mattermost-plugin-confluence/issues?q=is%3Aissue+is%3Aopen+sort%3Aupdated-desc+label%3A%22Up+For+Grabs%22+label%3A%22Help+Wanted%22)

## Overview

A Mattermost plugin for Confluence. Supports Confluence Cloud, Server, and Data Center versions, as well as environments with multiple Confluence servers. This plugin helps your teams collaborate and keep in sync as Confluence Pages and Spaces get updated. Comments and activity can be pushed into specific Mattermost channels for full visibility. 

## Features

With the Confluence plugin, you can subscribe to a variety of events in Confluence, and specify which channels the associated notifications will appear in. 

- In a Mattermost channel, setup subscriptions to get updates from Confluence and manage your notification subscriptions directly within Mattermost.
- Notify a channel whenever something occurs on a Confluence object:
  - Confluence spaces, including those created, updated, deleted, and restored, and those with added comments.
  - Confluence pages, including those created, updated, deleted, restored, and those with added, deleted, or updated comments.

## Configure notifications

### /confluence subscribe

Configure what events should send notifications to your Mattermost channel. When a user types `/confluence subscribe` in a channel, they will open a modal window that lets them configure a notification from Confluence to be delivered to the channel they are currently in. 

- The ``Alias`` (Subscription Name) is intended to be an easy to remember name for the subscription. You will use this name when you need to edit the configuration again. 

- `Confluence Base URL` is the URL of the Confluence server this rule is intended to come from. The Confluence server must have been setup by an administrator with Mattermost using the `/confluence install` command prior to using it in a subscription.

- `Subscribe To` is used to specify if they want to follow events for a Page or a Space object. 

- `Space Key` is the Confluence space key used for the project, often it is 2-4 characters, such as "PROJ" or "MM" and is unique for each Space on that Confluence server.

- `Page ID` is the ID of the Page object on Confuence. Since a page name can be changed by users, the underlying PageID is used to ensure tracking continues even if the page is renamed. The pageID of a Confluence page can be found by going to the "..." menu on the page, then selecting **Page Info**. The URL will then show the PageID in the URL at the end.

    ![image](https://github.com/mattermost/mattermost-plugin-confluence/assets/74422101/9314abd2-8562-456e-9661-7f23c91db206)
    
- `Events` are the internal confluence events that will trigger a notification from Confluence. The following events are currently included:
    - Confluence spaces, including those created, updated, deleted, and restored, and those with added comments.
    - Confluence pages, inlcuidng those created, updated, deleted, restored, and those with added, deleted, or updated comments.

Example of a configured notification:

![image](https://github.com/mattermost/mattermost-plugin-confluence/assets/74422101/33bc67f8-8d36-4e79-a386-7791f4dcd1ee)

### /confluence subscribe

Show a list of all the subscriptions set for the current channel. When you need to see what subscription rules are setup for a channel, you can run `/confluence list` to see a list of the configured subscriptions.

![image](https://github.com/mattermost/mattermost-plugin-confluence/assets/74422101/33c2a456-b7d1-41a2-ba55-53a492a7483c)

### /confluence edit

Edit existing subscription settings. To change the subscription settings, you need to pass over the subscription name to edit. If you have spaces in your subscription names, you need to wrap them in quotation marks. 

For example: `/confluence edit "Project A Subscription"`

To display a list of all the subscription names in the channel, type /confluence list then use the subscription names to edit the correct rule.  

![image](https://github.com/mattermost/mattermost-plugin-confluence/assets/74422101/81ec7a75-b6c1-4513-ad11-763f92416dc8)

### /confluence unsubscribe

Stop receiving notifications to a channel. To stop receiving notifications to a channel, use the `unsubscribe` command to specify the subscription that should be unsubscribed.
example: `/confluence unsubscribe "Project A Subscription"`.

## Development 

### Maintainers 

**Maintainer:** [@jfrerich](https://github.com/jfrerich)
**Co-Maintainer:** [@levb](https://github.com/levb)

### Platform and tools

- Make sure you have following components installed:

  - Go - v1.13 - [Getting Started](https://golang.org/doc/install)
    > **Note:** If you have installed Go to a custom location, make sure the `$GOROOT` variable is set properly. Refer [Installing to a custom location](https://golang.org/doc/install#install).
  - Make

## Set up CircleCI

Set up CircleCI to run the build job for each branch and build-and-release for each tag.

1. Go to [CircleCI Dashboard](https://circleci.com/dashboard).
2. In the top left, you will find the Org switcher. Select your Organisation.
3. If this is your first project on CircleCI, go to the Projects page, select the **Add Projects** button, then select the **Set Up Project** button next to your project. You may also select **Start Building** to manually trigger your first build.
4. To manage GitHub releases using CircleCI, you need to add your github personal access token to your project's environment variables.
   - Follow the instructions [here](https://help.github.com/en/articles/creating-a-personal-access-token-for-the-command-line) to create a personal access token. For CircleCI releases, you would need the `repo` scope.
   - Add the environment variable to your project as `GITHUB_TOKEN` by following the instructions [here](https://circleci.com/docs/2.0/env-vars/#setting-an-environment-variable-in-a-project).
