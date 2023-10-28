## Admin guide 

### Prerequisites

For Confluence Plugin, Mattermost server v5.19+ is required. Confluence Server version 7.x+ is supported. Confluence Data Center is not certified.

### Installation

To get started, a plugin that is installed on the Mattermost Server is needed to receive messages from the Confluence server. The notifications route it to the correct channel based on subscriptions that are configured. 

#### Marketplace installation (Recommended)

1. In Mattermost, go to **Main Menu > Plugin Marketplace**.
2. Search for "Confluence" or manually find the plugin from the list and select **Install**.
3. After the plugin has downloaded and been installed, select the **Configure** button.
4. Go to **Plugins Marketplace > Confluence**. 

    - Select the **Configure** button.
    - Generate a **Secret** for `Webhook Secret and Stats API Secret`.

5. Go to the top of the screen and set **Enable Plugin** to `True` and then select **Save** to enable the plugin.

#### Manual installation
1. Go to the [releases page of this GitHub repositiory](https://github.com/mattermost/mattermost-plugin-confluence/releases/latest), and download the latest release for your Mattermost Server.
2. Upload this file in the Mattermost **System Console > Plugins > Management** page to install the plugin.
3. Configure the Plugin from **System Console > Plugins > Confluence**.

### Install on Confluence

Now, you'll need to configure your Confluence server to communicate with the plugin on the Mattermost Server. The instructions are different for Cloud vs Server/Data Center. 

#### Set up Confluence Server

To get started, type in `/confluence install server` in a Mattermost chat window.

1. Download the [Mattermost for Confluence OBR file](https://github.com/mattermost/mattermost-for-confluence/releases) from the "Releases" tab in the Mattermost-Confluence plugin repository. Each release has some JVM files and an .OBR file.  Download it to your local computer. You will upload it to your Confluence server later in this process. 
2. Log in as an Administrator on Confluence to configure the plugin.
3. Create a new app in your Confluence Server by going to **Settings > Apps > Manage Apps**. For older versions of Confluence, go to **Administration > Applications > Add-ons > Manage add-ons**.
4. Choose **Settings** at the bottom of the page, enable development mode, and apply the changes. Development mode allows you to install apps from outside of the Atlassian Marketplace.
5. Press **Upload app**.

    ![image](https://github.com/mattermost/mattermost-plugin-confluence/assets/74422101/158bd7b7-4e36-41d5-872a-def2f617213f)

6. Upload the OBR file you downloaded earlier. 
7. Once the app is installed, press **Configure** to open the configuration page.
8. In the **Webhook URL field**, enter the URL that is displayed after you typed `/confluence install server` - it is unique to your server. You'll need to copy/paste it from Mattermost into the **From this URL** field in Confluence. 
9. Press **Save** to finish the setup.
10. Go to **Settings > Apps > Manage Apps**.
11. Choose **Settings** at the bottom of the page, enable development mode, and apply the change. Development mode allows you to install apps from outside of the Atlassian Marketplace.

    ![image](https://github.com/mattermost/mattermost-plugin-confluence/assets/74422101/fa27854e-6305-4164-963f-a5692284bf95)

12. Once installed, you will see the "Installed and ready to go!" message in Confluence.
13. You can now go to Mattermost and type `/confluence subscribe` in the channel you want to get notified by Confluence.

#### Set up Confluence Cloud

To get started, type in `/confluence install cloud` in a Mattermost chat window.

1. Log in as an Administrator on Confluence.
2. Go to **Settings > Apps > Manage Apps**.
3. Choose **Settings** at the bottom of the page, enable development mode, and apply the change. Development mode allows you to install apps from outside of the Atlassian Marketplace.

    ![image](https://github.com/mattermost/mattermost-plugin-confluence/assets/74422101/34471b0f-e54d-476c-b843-bf71355b6831)

4. Press **Upload App** button.

    ![image](https://github.com/mattermost/mattermost-plugin-confluence/assets/74422101/2ca5cc02-8a98-462d-b5f8-1cb20b7272d6)

5. In **From this URL**, enter the URL that is displayed after you typed `/confluence install cloud` - it is unique to your server. You'll need to copy/paste it from Mattermost into the **From this URL** field in Confluence. 
6. Once installed, you will see the "Installed and ready to go!" message in Confluence.
7. You can now go to Mattermost and type `/confluence subscribe` in the channel you want to get notified by Confluence.

