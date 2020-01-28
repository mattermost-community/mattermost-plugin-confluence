import Constants from '../constants';

export const subscriptionModal = (state = false, action) => {
    switch (action.type) {
    case Constants.ACTION_TYPES.OPEN_SUBSCRIPTION_MODAL:
        return true;
    case Constants.ACTION_TYPES.CLOSE_SUBSCRIPTION_MODAL:
        return false;
    default:
        return state;
    }
};

export const subscriptionEditModal = (state = {}, action) => {
    switch (action.type) {
    case Constants.ACTION_TYPES.RECEIVED_SUBSCRIPTION:
        return {
            editSubscription: true,
            alias: action.data.alias,
            baseURL: action.data.baseURL,
            spaceKey: action.data.spaceKey,
            events: action.data.events,
        };
    case Constants.ACTION_TYPES.CLOSE_SUBSCRIPTION_MODAL:
        return {};
    default:
        return state;
    }
};
