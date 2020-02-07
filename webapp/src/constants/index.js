import {PLUGIN_NAME} from './manifest';
import {ACTION_TYPES} from './action_types';

const CONFLUENCE_EVENTS = [
    {
        value: 'comment_created',
        label: 'Comment Create',
    },
    {
        value: 'comment_updated',
        label: 'Comment Update',
    },
    {
        value: 'comment_removed',
        label: 'Comment Remove',
    },
    {
        value: 'page_created',
        label: 'Page Create',
    },
    {
        value: 'page_updated',
        label: 'Page Update',
    },
    {
        value: 'page_trashed',
        label: 'Page Trash',
    },
    {
        value: 'page_restored',
        label: 'Page Restore',
    },
    {
        value: 'page_removed',
        label: 'Page Remove',
    },
];

const MATTERMOST_CSRF_COOKIE = 'MMCSRF';
const OPEN_EDIT_SUBSCRIPTION_MODAL_WEBSOCKET_EVENT = `custom_${PLUGIN_NAME}_open_edit_subscription_modal`;
const SPECIFY_ALIAS = 'Please specify alias.';

export default {
    ACTION_TYPES,
    CONFLUENCE_EVENTS,
    MATTERMOST_CSRF_COOKIE,
    OPEN_EDIT_SUBSCRIPTION_MODAL_WEBSOCKET_EVENT,
    PLUGIN_NAME,
    SPECIFY_ALIAS,
};
