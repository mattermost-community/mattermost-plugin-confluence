import {openSubscriptionModal, getChannelSubscription} from '../actions';

import {splitArgs} from '../utils';
import {sendEphemeralPost} from '../actions/subscription_modal';
import Constants from '../constants';

export default class Hooks {
    constructor(store) {
        this.store = store;
    }

    slashCommandWillBePostedHook = (message, contextArgs) => {
        let commandTrimmed;
        if (message) {
            commandTrimmed = message.trim();
        }
        const userID = this.store.getState().entities.users.currentUserId;
        if (commandTrimmed && commandTrimmed === '/confluence subscribe') {
            this.store.dispatch(openSubscriptionModal());
            return Promise.resolve({});
        } else if (commandTrimmed && commandTrimmed.startsWith('/confluence edit')) {
            const args = splitArgs(commandTrimmed);
            if (args.length < 3) { // eslint-disable-line
                this.store.dispatch(sendEphemeralPost(Constants.SPECIFY_ALIAS, contextArgs.channel_id, userID));
            } else {
                const alias = args.slice(2).join(' ');
                this.store.dispatch(getChannelSubscription(contextArgs.channel_id, alias, userID));
            }
            return Promise.resolve({});
        }
        return Promise.resolve({
            message,
            args: contextArgs,
        });
    }
}
