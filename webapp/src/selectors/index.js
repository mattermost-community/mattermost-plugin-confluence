import {id} from '../manifest';

const getPluginState = (state) => state[`plugins-${id}`] || {};

const isSubscriptionModalVisible = (state) => getPluginState(state).subscriptionModal;

const isSubscriptionEditModalVisible = (state) => getPluginState(state).subscriptionEditModal;

const postMessage = (state) => getPluginState(state).createConfluencePageModal;

const installedInstances = (state) => getPluginState(state).installedInstances;

const spacesForConfluenceURL = (state) => getPluginState(state).spacesForConfluenceURL;

const isUserConnected = (state) => getPluginState(state).userConnected;

export default {
    isSubscriptionModalVisible,
    isSubscriptionEditModalVisible,
    postMessage,
    isUserConnected,
    installedInstances,
    spacesForConfluenceURL,
};
