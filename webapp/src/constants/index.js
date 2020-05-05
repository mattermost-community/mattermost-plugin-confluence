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

const SUBSCRIPTION_TYPE = [
    {
        value: 'space_subscription',
        label: 'Space',
    },
    {
        value: 'page_subscription',
        label: 'Page',
    },
];

const MATTERMOST_CSRF_COOKIE = 'MMCSRF';
const OPEN_EDIT_SUBSCRIPTION_MODAL_WEBSOCKET_EVENT = `custom_${PLUGIN_NAME}_open_edit_subscription_modal`;
const SPECIFY_ALIAS = 'Please specify alias.';

const COMMAND_ADMIN_ONLY = '`/confluence` commands can only be run by a system administrator.';
const SYSTEM_ADMIN_ROLE = 'system_admin';

export default {
    ACTION_TYPES,
    CONFLUENCE_EVENTS,
    MATTERMOST_CSRF_COOKIE,
    OPEN_EDIT_SUBSCRIPTION_MODAL_WEBSOCKET_EVENT,
    PLUGIN_NAME,
    SPECIFY_ALIAS,
    COMMAND_ADMIN_ONLY,
    SYSTEM_ADMIN_ROLE,
    SUBSCRIPTION_TYPE,
};
