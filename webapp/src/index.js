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
            Constants.OPEN_EDIT_SUBSCRIPTION_MODAL_WEBSOCKET_EVENT,
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
