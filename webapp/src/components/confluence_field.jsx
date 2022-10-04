import React from 'react';
import PropTypes from 'prop-types';
import {ControlLabel, FormControl, FormGroup} from 'react-bootstrap';
import Select from 'react-select';

import {getStyleForReactSelect} from '../utils/react_select_styles';

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
        addValidation: PropTypes.func,
        removeValidation: PropTypes.func,
        theme: PropTypes.object,
        fieldType: PropTypes.string.isRequired,
        type: PropTypes.string,
        readOnly: PropTypes.bool,
        formGroupStyle: PropTypes.object,
        formControlStyle: PropTypes.object,
    };

    static defaultProps = {
        readOnly: false,
        formGroupStyle: {},
        formControlStyle: {},
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
        const {fieldType, value, required} = this.props;
        if (required &&
            (value === null ||
            (typeof value === 'string' && !value.trim()) ||
            (fieldType === 'dropDown' && value.length === 0) ||
            !value)
        ) {
            this.setState({
                valid: false,
            });
            return false;
        }
        return true;
    };

    render() {
        const {
            required, fieldType, theme, label, formGroupStyle, formControlStyle, type,
        } = this.props;
        const requiredErrorMsg = 'This field is required.';
        let requiredError = null;
        if (required && !this.state.valid) {
            requiredError = (
                <p className='help-text error-text'>
                    <span>{requiredErrorMsg}</span>
                </p>
            );
        }
        let field = null;
        if (fieldType === 'input' && type === 'textarea') {
            field = (
                <FormControl
                    style={formControlStyle}
                    componentClass='textarea'
                    {...this.props}
                    onChange={this.handleChange}
                />
            );
        } else if (fieldType === 'input' && type === 'text') {
            field = (
                <FormControl
                    style={formControlStyle}
                    {...this.props}
                    onChange={this.handleChange}
                />
            );
        } else if (fieldType === 'dropDown') {
            field = (
                <Select
                    {...this.props}
                    menuPortalTarget={document.body}
                    menuPlacement='auto'
                    styles={getStyleForReactSelect(theme)}
                    onChange={this.handleChange}
                />
            );
        }
        return (
            <FormGroup style={formGroupStyle}>
                <ControlLabel>{label}</ControlLabel>
                {required &&
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
