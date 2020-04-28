import Constants from './constants';
import Hooks from './hooks';
import reducer from './reducers';

import SubscriptionModal from './components/subscription_modal';

const displayNameError = 'You attempted to set the key `display_name` with the value `""` on an object that is meant to be immutable and has been frozen.';

//
// Define the plugin class that will register
// our plugin components.
//
class PluginClass {
    initialize(registry, store) {
        try {
            registry.registerReducer(reducer);
        } catch (e) {
            // If it is the display name error, ignore. Otherwise re-throw error.
            if (!e.toString().includes(displayNameError)) {
                throw e;
            }
        }
        registry.registerRootComponent(SubscriptionModal);
        const hooks = new Hooks(store);
        registry.registerSlashCommandWillBePostedHook(hooks.slashCommandWillBePostedHook);
    }
}

//
// To register your plugin you must expose it on window.
//
window.registerPlugin(Constants.PLUGIN_NAME, new PluginClass());
