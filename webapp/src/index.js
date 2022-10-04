import {id} from './manifest';

import Hooks from './hooks';
import reducer from './reducers';

import SubscriptionModal from './components/subscription_modal';
import TokenSetting from './components/admin_settings/token_setting';
import CreateConfluencePage from './components/create_confluence_page_modal';

//
// Define the plugin class that will register
// our plugin components.
//
class PluginClass {
    initialize(registry, store) {
        registry.registerReducer(reducer);
        registry.registerRootComponent(CreateConfluencePage);
        registry.registerRootComponent(SubscriptionModal);
        const hooks = new Hooks(store);
        registry.registerSlashCommandWillBePostedHook(hooks.slashCommandWillBePostedHook);
        registry.registerAdminConsoleCustomSetting('tokens', TokenSetting);
        registry.registerPostDropdownMenuAction('Create Confluence Page', hooks.createConfluencePage);
    }
}

//
// To register your plugin you must expose it on window.
//
window.registerPlugin(id, new PluginClass());
