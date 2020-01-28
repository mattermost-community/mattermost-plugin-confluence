import {openSubscriptionModal} from '../actions';

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
        }
        return Promise.resolve({message, args: contextArgs});
    }
}
