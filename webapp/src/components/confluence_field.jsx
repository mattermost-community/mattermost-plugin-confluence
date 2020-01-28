import React from 'react';
import PropTypes from 'prop-types';
import {ControlLabel, FormControl, FormGroup} from 'react-bootstrap';
import Select from 'react-select';

import {getReactSelectTheme, reactSelectStyles} from '../utils/react_select_styles';

export default class ConfluenceField extends React.PureComponent {
    static propTypes = {
        required: PropTypes.bool.isRequired,
        value: PropTypes.PropTypes.oneOfType([
            PropTypes.string,
            PropTypes.number,
            PropTypes.object,
            PropTypes.array,
        ]),
        label: PropTypes.string.isRequired,
        onChange: PropTypes.func.isRequired,
        addValidation: PropTypes.func.isRequired,
        removeValidation: PropTypes.func.isRequired,
        theme: PropTypes.object,
        fieldType: PropTypes.string.isRequired,
        readOnly: PropTypes.bool.isRequired,
    };

    constructor(props) {
        super(props);
        this.state = {
            valid: true,
        };
    }

    componentDidMount() {
        if (this.props.addValidation) {
            this.props.addValidation(this.isValid);
        }
    }

    handleChange = (e) => {
        if (!this.state.valid) {
            this.setState({
                valid: true,
            });
        }
        this.props.onChange(e);
    };

    componentWillUnmount() {
        if (this.props.removeValidation) {
            this.props.removeValidation(this.isValid);
        }
    }

    isValid = () => {
        if (this.props.required && !this.props.value) {
            this.setState({
                valid: false,
            });
            return false;
        }
        return true;
    };

    render() {
        const requiredErrorMsg = 'This field is required.';
        let requiredError = null;
        if (this.props.required && !this.state.valid) {
            requiredError = (
                <p className='help-text error-text'>
                    <span>{requiredErrorMsg}</span>
                </p>
            );
        }
        let field = null;
        if (this.props.fieldType === 'input') {
            field = (
                <FormControl
                    {...this.props}
                    onChange={this.handleChange}
                />
            );
        } else if (this.props.fieldType === 'dropDown') {
            field = (
                <Select
                    {...this.props}
                    menuPortalTarget={document.body}
                    menuPlacement='auto'
                    styles={reactSelectStyles}
                    theme={getReactSelectTheme(this.props.theme)}
                    onChange={this.handleChange}
                />
            );
        }
        return (
            <FormGroup>
                <ControlLabel>{this.props.label}</ControlLabel>
                {this.props.required &&
                <span
                    className='error-text'
                    style={{marginLeft: '3px'}}
                >
                    {'*'}
                </span> }
                {field}
                {requiredError}
            </FormGroup>
        );
    }
}
