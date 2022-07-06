// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import React from 'react';
import ReactSelect from 'react-select';
import AsyncSelect, {Props as ReactSelectProps} from 'react-select/async';
import {ValueType} from 'react-select/src/types';

import Setting from 'src/components/setting';
import {getStyleForReactSelect} from 'src/utils/styles';
import {ReactSelectOption} from 'src/types';

const MAX_NUM_OPTIONS = 100;

type Omit<T, K extends keyof T> = Pick<T, Exclude<keyof T, K>>

export type Props = Omit<ReactSelectProps<ReactSelectOption>, 'theme'> & {
    theme: any;
    addValidate?: (isValid: () => boolean) => void;
    removeValidate?: (isValid: () => boolean) => void;
    limitOptions?: boolean;
};

type State = {
    invalid: boolean;
};

export default class ReactSelectSetting extends React.PureComponent<Props, State> {
    state: State = {invalid: false};

    componentDidMount() {
        if (this.props.addValidate) {
            this.props.addValidate(this.isValid);
        }
    }

    componentWillUnmount() {
        if (this.props.removeValidate) {
            this.props.removeValidate(this.isValid);
        }
    }

    componentDidUpdate(prevProps: Props, prevState: State) {
        if (prevState.invalid && this.props.value?.value !== prevProps.value?.value) {
            this.setState({invalid: false}); //eslint-disable-line react/no-did-update-set-state
        }
    }

    handleChange = (value: ValueType<ReactSelectOption> | ReactSelectOption | ReactSelectOption[]) => {
        if (this.props.onChange) {
            if (Array.isArray(value)) {
                this.props.onChange(this.props.name, value.map((x) => x.value));
            } else {
                const newValue = value ? (value as ReactSelectOption).value : null;
                this.props.onChange(this.props.name, newValue);
            }
        }
    };

    // Standard search term matching plus reducing to < 100 items
    filterOptions = (input: string) => {
        let options = this.props.options;
        if (input) {
            options = options.filter((opt: ReactSelectOption) => `${opt.label}`.toUpperCase().includes(input.toUpperCase()));
        }

        return Promise.resolve(options.slice(0, MAX_NUM_OPTIONS));
    };

    isValid = () => {
        if (!this.props.required) {
            return true;
        }

        const valid = Array.isArray(this.props.value) ? Boolean(this.props.value.length) : Boolean(this.props.value);

        this.setState({invalid: !valid});
        return valid;
    };

    render() {
        const {theme} = this.props.theme;
        const requiredMsg = 'This field is required.';
        const validationError = this.props.required && this.state.invalid ? (
            <p className='help-text error-text'>
                <span>{requiredMsg}</span>
            </p>) : null;

        const selectComponent = (this.props.limitOptions && this.props.options.length > MAX_NUM_OPTIONS) ?

            // The parent component helps us know that we may have a large number of options, and that
            // the data-set is static. In this case, we use the AsyncSelect component and synchronous func
            // "filterOptions" to limit the number of options being rendered at a given time.
            (
                <AsyncSelect
                    {...this.props}
                    loadOptions={this.filterOptions}
                    defaultOptions={true}
                    menuPortalTarget={document.body}
                    menuPlacement='auto'
                    onChange={this.handleChange}
                    styles={getStyleForReactSelect(this.props.theme)}
                />
            ) : (
                <ReactSelect
                    {...this.props}
                    menuPortalTarget={document.body}
                    menuPlacement='auto'
                    onChange={this.handleChange}
                    styles={getStyleForReactSelect(this.props.theme)}
                />
            );

        return (
            <Setting
                inputId={this.props.name}
                {...this.props as any}
            >
                {selectComponent}
                {validationError}
            </Setting>
        );
    }
}
