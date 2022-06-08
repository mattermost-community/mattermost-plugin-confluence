import {id} from './manifest';

import Hooks from './hooks';
import reducer from './reducers';

import SubscriptionModal from './components/subscription_modal';
import TokenSetting from './components/admin_settings/token_setting';

//
// Define the plugin class that will register
// our plugin components.
//
class PluginClass {
    initialize(registry, store) {
        registry.registerReducer(reducer);
        registry.registerRootComponent(SubscriptionModal);
        const hooks = new Hooks(store);
        registry.registerSlashCommandWillBePostedHook(hooks.slashCommandWillBePostedHook);
        registry.registerAdminConsoleCustomSetting('tokens', TokenSetting);
    }
}

//
// To register your plugin you must expose it on window.
//
window.registerPlugin(id, new PluginClass());
