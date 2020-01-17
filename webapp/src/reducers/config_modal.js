import Constants from '../constants';

export const configModal = (state = false, action) => {
    switch (action.type) {
    case Constants.ACTION_TYPES.SHOW_CONFIG_MODAL:
        return true;
    case Constants.ACTION_TYPES.CLOSE_CONFIG_MODAL:
        return false;
    default:
        return state;
    }
};
