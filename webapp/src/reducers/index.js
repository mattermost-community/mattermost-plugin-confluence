import {combineReducers} from 'redux';

import {subscriptionModal, subscriptionEditModal, installedInstances, userConnected, createConfluencePageModal, spacesForConfluenceURL} from './subscription_modal';

export default combineReducers({
    subscriptionModal,
    subscriptionEditModal,
    installedInstances,
    userConnected,
    createConfluencePageModal,
    spacesForConfluenceURL,
});
