import Constants from '../constants';

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
