{{ .ConfluenceURL }} has been successfully added. To finish the configuration, add an Application Link in your Confluence instance following these steps:

1. Go to [**Settings > Applications > Application
   Links**]({{ .ConfluenceURL }}/plugins/servlet/applinks/listApplicationLinks)
2. Select **Create link**.
3. On the **Create Link** screen, select **External Application** and **Incoming** as
   `Application type` and `Direction` respectively. Select **Continue**.
4. On the **Link Applications** screen, set the following values, then select **Continue**:
   **Name**: `Mattermost`
   **Redirect URL**: ```{{ .RedirectURL }}```.
   **Application Permissions**: `Admin`
5. Copy the `clientID` and `clientSecret` and paste them into the plugin configuration.
6. In Mattermost, use the "/confluence connect" slash command to connect your Mattermost account with your
   Confluence account.

If you see an option to create a Confluence issue, you're all set! If not, refer to our [documentation](https://mattermost.gitbook.io/plugin-confluence) for troubleshooting help.
