import {PLUGIN_NAME} from './manifest';
import {ACTION_TYPES} from './action_types';

const CONFLUENCE_EVENTS = [
    {value: 'comment_create', label: 'Comment Create'},
    {value: 'comment_update', label: 'Comment Update'},
    {value: 'comment_delete', label: 'Comment Delete'},
    {value: 'page_create', label: 'Page Create'},
    {value: 'page_update', label: 'Page Update'},
    {value: 'page_delete', label: 'Page Delete'},
];

const MATTERMOST_CSRF_COOKIE = 'MMCSRF';
const OPEN_EDIT_SUBSCRIPTION_MODAL_WEBSOCKET_EVENT = `custom_${PLUGIN_NAME}_open_edit_subscription_modal`;

export default {
    ACTION_TYPES,
    CONFLUENCE_EVENTS,
    PLUGIN_NAME,
    MATTERMOST_CSRF_COOKIE,
    OPEN_EDIT_SUBSCRIPTION_MODAL_WEBSOCKET_EVENT,
};
