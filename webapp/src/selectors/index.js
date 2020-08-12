import {id} from '../manifest';

const getPluginState = (state) => state[`plugins-${id}`] || {};

const isSubscriptionModalVisible = (state) => getPluginState(state).subscriptionModal;

const isSubscriptionEditModalVisible = (state) => getPluginState(state).subscriptionEditModal;

export default {
    isSubscriptionModalVisible,
    isSubscriptionEditModalVisible,
};
