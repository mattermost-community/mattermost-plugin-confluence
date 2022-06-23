import {id} from '../manifest';

export const ACTION_TYPES = {
    OPEN_SUBSCRIPTION_MODAL: id + '_open_subscription_modal',
    CLOSE_SUBSCRIPTION_MODAL: id + '_close_subscription_modal',
    OPEN_CREATE_CONFLUENCE_PAGE_MODAL: id + '_open_create_confluence_page_modal',
    CLOSE_CREATE_CONFLUENCE_PAGE_MODAL: id + '_close_create_confluence_page_modal',
    RECEIVED_SUBSCRIPTION: id + '_received_subscription',
    RECEIVED_CONNECTED: id + '_connected',
    RECEIVED_INSTANCE_STATUS: id + '_instance_status',
    RECEIVED_CONFLUENCE_INSTANCE: id + '_received_confluence_instance',
};
