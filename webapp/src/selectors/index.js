import Constants from '../constants';

const getPluginState = (state) => state[`plugins-${Constants.PLUGIN_NAME}`] || {};

const isConfigModalvisible = (state) => getPluginState(state).configModal;

export default {
    isConfigModalvisible,
};
