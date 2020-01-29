import request from 'superagent';
import Cookies from 'js-cookie';

import Constants from '../constants';

/**
 *  Add web utilities for interacting with servers here
 */
export default class Client {
    constructor() {
        const url = new URL(window.location.href);
        this.baseUrl = `${url.protocol}//${url.host}`;
        this.pluginUrl = `${this.baseUrl}/plugins/${Constants.PLUGIN_NAME}`;
        this.apiUrl = `${this.baseUrl}/api/v4`;
    }

    saveChannelSubscription = (channelSubscription) => {
        const url = `${this.pluginUrl}/save-channel-subscription`;
        return this.doPost(url, channelSubscription);
    };

    editChannelSubscription = (channelSubscription) => {
        const url = `${this.pluginUrl}/edit-channel-subscription`;
        return this.doPost(url, channelSubscription);
    };

    doGet = async (url, headers = {}) => {
        headers['X-Requested-With'] = 'XMLHttpRequest';

        const response = await request.
            get(url).
            set(headers).
            type('application/json').
            accept('application/json');

        return response.body;
    };

    doPost = async (url, body, headers = {}) => {
        headers['X-Requested-With'] = 'XMLHttpRequest';
        headers['X-CSRF-Token'] = Cookies.get(Constants.MATTERMOST_CSRF_COOKIE);

        const response = await request.
            post(url).
            send(body).
            set(headers).
            type('application/json').
            accept('application/json');

        return response.body;
    };

    doPatch = async (url, body, headers= {}) => {
        headers['X-Requested-With'] = 'XMLHttpRequest';
        headers['X-CSRF-Token'] = Cookies.get(Constants.MATTERMOST_CSRF_COOKIE);

        const response = await request.
            put(url).
            send(body).
            set(headers).
            type('application/json').
            accept('application/json');

        return response.body;
    }
}
