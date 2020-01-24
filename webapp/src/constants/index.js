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

const CONFLUENCE_TYPE = [
    {value: 'cloud', label: 'Cloud'},
    {value: 'server', label: 'Server'},
];

const MATTERMOST_CSRF_COOKIE = 'MMCSRF';

export default {
    ACTION_TYPES,
    CONFLUENCE_TYPE,
    CONFLUENCE_EVENTS,
    PLUGIN_NAME,
    MATTERMOST_CSRF_COOKIE,
};
