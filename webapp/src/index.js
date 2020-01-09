import Constants from './constants';

//
// Define the plugin class that will register
// our plugin components.
//
export default class PluginClass {
    // eslint-disable-next-line no-unused-vars
    initialize(registry, store) {
        // @see https://developers.mattermost.com/extend/plugins/webapp/reference/
    }
}

//
// To register your plugin you must expose it on window.
//
window.registerPlugin(Constants.PLUGIN_NAME, new PluginClass());
