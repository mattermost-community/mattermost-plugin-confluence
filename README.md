# Mattermost Plugin for Confluence 

[![Build Status](https://img.shields.io/circleci/project/github/mattermost/mattermost-plugin-confluence/master)](https://circleci.com/gh/mattermost/mattermost-plugin-confluence)
[![Code Coverage](https://img.shields.io/codecov/c/github/mattermost/mattermost-plugin-confluence/master)](https://codecov.io/gh/mattermost/mattermost-plugin-confluence)
[![Release](https://img.shields.io/github/v/release/mattermost/mattermost-plugin-confluence)](https://github.com/mattermost/mattermost-plugin-confluence/releases/latest)
[![HW](https://img.shields.io/github/issues/mattermost/mattermost-plugin-confluence/Up%20For%20Grabs?color=dark%20green&label=Help%20Wanted)](https://github.com/mattermost/mattermost-plugin-confluence/issues?q=is%3Aissue+is%3Aopen+sort%3Aupdated-desc+label%3A%22Up+For+Grabs%22+label%3A%22Help+Wanted%22)

A Mattermost plugin for Confluence. Supports Confluence Cloud, Server, and Data Center versions, as well as environments with multiple Confluence servers. This plugin helps your teams collaborate and keep in sync as Confluence Pages and Spaces get updated. Comments and activity can be pushed into specific Mattermost channels for full visibility. 

With the Confluence plugin, you can subscribe to a variety of events in Confluence and specify which channels the associated notifications will appear in. 

- In a Mattermost channel, setup subscriptions to get updates from Confluence and manage your notification subscriptions directly within Mattermost
- Notify a channel whenever something occurs on a Confluence object:
  - Confluence spaces, including those created, updated, deleted, and restored, and those with added comments.
  - Confluence pages, inlcuidng those created, updated, deleted, restored, and those with added, deleted, or updated comments.

## Admin Guide 

### Prerequisites

For Confluence Plugin, Mattermost server v5.19+ is required. Confluence Server version 7.x+ is supported. Confluence Data Center is not certified.

### Install the Mattermost Plugin

To get started, a plugin that is installed on the Mattermost Server is needed to receive messages from the Confluence server and the notifications route it to the correct channel based on subscriptions that are configured. 

#### Install via Plugin Marketplace (Recommended)

1. In Mattermost, go to **Main Menu > Plugin Marketplace**.
2. Search for "Confluence" or manually find the plugin from the list and click **Install**.
3. After the plugin has downloaded and been installed, click the **Configure** button.
4. Go to **Plugins Marketplace > Confluence**. 

    - Click the **Configure** button.
    - Generate a **Secret** for `Webhook Secret and Stats API Secret`.

5. Go to the top of the screen and set **Enable Plugin** to `True` and then click **Save** to enable the plugin.

#### (Alternative) Install via Manual Upload

1. Go to the [releases page of this GitHub repositiory](https://github.com/mattermost/mattermost-plugin-confluence/releases/latest), and download the latest release for your Mattermost Server.
2. Upload this file in the Mattermost **System Console > Plugins > Management** page to install the plugin.
3. Configure the Plugin from **System Console > Plugins > Confluence**.

### Install on Confluence

Now, you'll need to configure your Confluence server to communicate with the plugin on the Mattermost Server. The instructions are different for Cloud vs Server/Data Center. 

#### Set Up Confluence Server

To get started, type in `/confluence install server` in a Mattermost chat window.

1. Download the [Mattermost for Confluence OBR file](https://github.com/mattermost/mattermost-for-confluence/releases) from the "Releases" tab in the Mattermost-Confluence plugin repository. Each release has some JVM files and an .OBR file.  Download it to your local computer. You will upload it to your Confluence server later in this process. 
2. Log in as an Administrator on Confluence to configure the plugin.
3. Create a new app in your Confluence Server by navigating to **Settings > Apps > Manage Apps**. For older versions of Confluence, navigate to **Administration > Applications > Add-ons > Manage add-ons**.
4. Choose **Settings** at the bottom of the page, enable development mode, and apply the changes. Development mode allows you to install apps from outside of the Atlassian Marketplace.
5. Press **Upload app**.

    ![image](https://github.com/mattermost/mattermost-plugin-confluence/assets/74422101/158bd7b7-4e36-41d5-872a-def2f617213f)

6. Upload the OBR file you downloaded earlier. 
7. Once the app is installed, press **Configure** to open the configuration page.
8. In the **Webhook URL field**, enter the URL that is displayed after you typed `/confluence install server` - it is unique to your server. You'll need to copy/paste it from Mattermost into the **From this URL** field in Confluence. 
9. Press **Save** to finish the setup.
10. Navigate to **Settings > Apps > Manage Apps**.
11. Choose **Settings** at the bottom of the page, enable development mode, and apply the change. Development mode allows you to install apps from outside of the Atlassian Marketplace.

    ![image](https://github.com/mattermost/mattermost-plugin-confluence/assets/74422101/fa27854e-6305-4164-963f-a5692284bf95)

12. Once installed, you will see the "Installed and ready to go!" message in Confluence.
13. You can now go to Mattermost and type `/confluence subscribe` in the channel you want to get notified by Confluence.

#### Set Up Confluence Cloud

To get started, type in `/confluence install cloud` in a Mattermost chat window.

1. Log in as an Administrator on Confluence.
2. Navigate to **Settings > Apps > Manage Apps**.
3. Choose **Settings** at the bottom of the page, enable development mode, and apply the change. Development mode allows you to install apps from outside of the Atlassian Marketplace.

    ![image](https://github.com/mattermost/mattermost-plugin-confluence/assets/74422101/34471b0f-e54d-476c-b843-bf71355b6831)

4. Press **Upload App** button.

    ![image](https://github.com/mattermost/mattermost-plugin-confluence/assets/74422101/2ca5cc02-8a98-462d-b5f8-1cb20b7272d6)

5. In **From this URL**, enter the URL that is displayed after you typed `/confluence install cloud` - it is unique to your server. You'll need to copy/paste it from Mattermost into the **From this URL** field in Confluence. 
6. Once installed, you will see the "Installed and ready to go!" message in Confluence.
7. You can now go to Mattermost and type `/confluence subscribe` in the channel you want to get notified by Confluence.

## Configure Notifications

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
