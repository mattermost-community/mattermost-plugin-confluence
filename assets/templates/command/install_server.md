{{ .ConfluenceURL }} has been successfully added. To finish the configuration, add an Application Link in your Confluence instance following these steps:

1. Navigate to [**Settings > Applications > Application
   Links**]({{ .ConfluenceURL }}/plugins/servlet/applinks/listApplicationLinks)
2. Click **Create link**.
3. In **Create Link** screen, select **External Application** and **Incoming** as
   `Application type` and `Direction` respectively. Click **Continue**.
4. In **Link Applications** screen, set the following values:
**Name**: `Mattermost`
**Redirect URL**: ```{{ .RedirectURL }}```.
**Application Permissions**: `Admin`
Click **Continue**
5. Copy the `clientID` and `clientSecret` and add them into the plugin configuration.
6. Use the "/confluence connect" command to connect your Mattermost account with your
   Confluence account.
7. Click the "More Actions" (...) option of any message in the channel
   (available when you hover over a message).

If you see an option to create a Confluence issue, you're all set! If not, refer to our [documentation](https://mattermost.gitbook.io/plugin-confluence) for troubleshooting help.
