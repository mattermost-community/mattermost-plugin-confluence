import Constants from './constants';
import Hooks from './hooks';
import reducer from './reducers';

import {receivedSubscription} from './actions';
import SubscriptionModal from './components/subscription_modal';

//
// Define the plugin class that will register
// our plugin components.
//
class PluginClass {
    initialize(registry, store) {
        registry.registerReducer(reducer);
        registry.registerRootComponent(SubscriptionModal);
        registry.registerWebSocketEventHandler(
            `custom_${Constants.PLUGIN_NAME}_open_edit_subscription_modal`,
            (payload) => {
                store.dispatch(receivedSubscription(payload.data.subscription));
            },
        );
        const hooks = new Hooks(store);
        registry.registerSlashCommandWillBePostedHook(hooks.slashCommandWillBePostedHook);
    }
}

//
// To register your plugin you must expose it on window.
//
window.registerPlugin(Constants.PLUGIN_NAME, new PluginClass());
