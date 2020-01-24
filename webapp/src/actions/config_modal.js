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

export const openConfigModal = () => (dispatch) => {
    dispatch({
        type: Constants.ACTION_TYPES.SHOW_CONFIG_MODAL,
    });
};

export const closeConfigModal = () => (dispatch) => {
    dispatch({
        type: Constants.ACTION_TYPES.CLOSE_CONFIG_MODAL,
    });
};
