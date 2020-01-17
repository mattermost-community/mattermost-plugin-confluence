import Constants from './constants';
import Hooks from './hooks';
import reducer from './reducers';

import ConfigModal from './components/config_modal';

//
// Define the plugin class that will register
// our plugin components.
//
class PluginClass {
    initialize(registry, store) {
        registry.registerReducer(reducer);
        registry.registerRootComponent(ConfigModal);
        const hooks = new Hooks(store);
        registry.registerSlashCommandWillBePostedHook(hooks.slashCommandWillBePostedHook);
    }
}

//
// To register your plugin you must expose it on window.
//
window.registerPlugin(Constants.PLUGIN_NAME, new PluginClass());
