import Client from '../client';
import Constants from '../constants';

export const saveChannelSubscription = (body) => {
    return async () => {
        let data = null;
        try {
            data = await Client.saveChannelSubscription(body);
        } catch (error) {
            return {data, error};
        }

        return {data, error: null};
    };
};

export const editChannelSubscription = (body) => {
    return async () => {
        let data = null;
        try {
            data = await Client.editChannelSubscription(body);
        } catch (error) {
            return {data, error};
        }

        return {data, error: null};
    };
};

export const openSubscriptionModal = () => (dispatch) => {
    dispatch({
        type: Constants.ACTION_TYPES.OPEN_SUBSCRIPTION_MODAL,
    });
};

export const closeSubscriptionModal = () => (dispatch) => {
    dispatch({
        type: Constants.ACTION_TYPES.CLOSE_SUBSCRIPTION_MODAL,
    });
};

export const receivedSubscription = (subscription) => (dispatch) => {
    dispatch({
        type: Constants.ACTION_TYPES.RECEIVED_SUBSCRIPTION,
        data: JSON.parse(subscription),
    });
};

