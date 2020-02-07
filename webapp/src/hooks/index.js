import {openSubscriptionModal, openEditSubscriptionModal} from '../actions';

export default class Hooks {
    constructor(store) {
        this.store = store;
    }

    slashCommandWillBePostedHook = (message, contextArgs) => {
        let commandTrimmed;
        if (message) {
            commandTrimmed = message.trim();
        }

        if (commandTrimmed && commandTrimmed === '/confluence subscribe') {
            this.store.dispatch(openSubscriptionModal());
            return Promise.resolve({});
        } else if (commandTrimmed && commandTrimmed.startsWith('/confluence edit')) {
            const data = {
                message,
                channelID: contextArgs.channel_id,
            };
            this.store.dispatch(openEditSubscriptionModal(data, this.store.getState().entities.users.currentUserId));
            return Promise.resolve({});
        }
        return Promise.resolve({
            message,
            args: contextArgs,
        });
    }
}
