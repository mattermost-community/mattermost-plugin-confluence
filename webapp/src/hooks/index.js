import {openConfigModal} from '../actions';

export default class Hooks {
    constructor(store) {
        this.store = store;
    }

    slashCommandWillBePostedHook = (command, contextArgs) => {
        let commandTrimmed;
        if (command) {
            commandTrimmed = command.trim();
        }

        if (commandTrimmed && commandTrimmed === '/confluence config') {
            this.store.dispatch(openConfigModal());
            return Promise.resolve({});
        }
        return Promise.resolve({command, args: contextArgs});
    }
}
