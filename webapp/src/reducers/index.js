import {combineReducers} from 'redux';

import {subscriptionModal, subscriptionEditModal} from './subscription_modal';

export default combineReducers({
    subscriptionModal,
    subscriptionEditModal,
});
