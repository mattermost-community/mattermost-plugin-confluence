export default class Validator {
    constructor() {
        this.components = [];
    }

    addValidation = (validateField) => {
        this.components.push(validateField);
    };

    removeValidation = (validateField) => {
        const index = this.components.indexOf(validateField);
        if (index !== -1) {
            this.components.splice(index, 1);
        }
    };

    validate = () => {
        return Array.from(this.components.values()).reduce((accum, validateField) => {
            return validateField() && accum;
        }, true);
    };
}
