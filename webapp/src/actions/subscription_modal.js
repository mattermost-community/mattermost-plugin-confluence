import {PostTypes} from 'mattermost-redux/action_types';

import Client from '../client';
import Constants from '../constants';

export const saveChannelSubscription = (body) => {
    return async () => {
        let data = null;
        try {
            data = await Client.saveChannelSubscription(body);
        } catch (error) {
            return {
                data,
                error,
            };
        }

        return {
            data,
            error: null,
        };
    };
};

export const editChannelSubscription = (body) => {
    return async () => {
        let data = null;
        try {
            data = await Client.editChannelSubscription(body);
        } catch (error) {
            return {
                data,
                error,
            };
        }

        return {
            data,
            error: null,
        };
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

export const openEditSubscriptionModal = (body, userID) => async (dispatch) => {
    try {
        const response = await Client.openEditSubscriptionModal(body);
        dispatch({
            type: Constants.ACTION_TYPES.RECEIVED_SUBSCRIPTION,
            data: response,
        });
    } catch (e) {
        const timestamp = Date.now();
        const post = {
            id: 'confluencePlugin' + timestamp,
            user_id: userID,
            channel_id: body.channelID,
            message: e.response.text,
            type: 'system_ephemeral',
            create_at: timestamp,
            update_at: timestamp,
            root_id: '',
            parent_id: '',
            props: {},
        };

        dispatch({
            type: PostTypes.RECEIVED_NEW_POST,
            data: post,
            channelId: body.channelID,
        });
    }
};
