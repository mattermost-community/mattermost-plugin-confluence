import React from 'react';
import {
    Modal,
    Button,
} from 'react-bootstrap';

import PropTypes from 'prop-types';

import Constants from '../../constants';
import ConfluenceField from '../confluence_field';
import Validator from '../validator';

export default class ConfigModal extends React.PureComponent {
    static propTypes = {
        visibility: PropTypes.bool.isRequired,
        close: PropTypes.func.isRequired,
        theme: PropTypes.object.isRequired,
    };

    constructor(props) {
        super(props);
        this.state = {
            type: Constants.CONFLUENCE_TYPE[0],
            baseURL: '',
            spaceKey: '',
            events: null,
        };
        this.validator = new Validator();
    }

    handleClose = () => {
        this.props.close();
    };

    handleType = (type) => {
        this.setState({
            type,
        });
    };

    handleBaseURLChange = (e) => {
        this.setState({
            baseURL: e.target.value,
        });
    };

    handleSpaceKey = (e) => {
        this.setState({
            spaceKey: e.target.value,
        });
    };

    handleEvents = (events) => {
        this.setState({
            events,
        });
    };

    handleSubmit = () => {
        if (!this.validator.validate()) {
            return;
        }

        // TODO: SAVE CONFIG
        this.handleClose();
    };

    render() {
        const {visibility} = this.props;

        return (
            <Modal
                show={visibility}
                onHide={this.handleClose}
                backdrop={'static'}
            >
                <Modal.Header closeButton={true}>
                    <Modal.Title>
                        {'Channel Settings'}
                    </Modal.Title>
                </Modal.Header>
                <Modal.Body>
                    <div>
                        <ConfluenceField
                            label={'TYPE'}
                            name={'type'}
                            fieldType={'dropDown'}
                            required={true}
                            theme={this.props.theme}
                            options={Constants.CONFLUENCE_TYPE}
                            value={this.state.type}
                            addValidation={this.validator.addValidation}
                            removeValidation={this.validator.removeValidation}
                            onChange={this.handleType}
                        />
                        <ConfluenceField
                            label={'CONFLUENCE BASE URL'}
                            type={'text'}
                            fieldType={'input'}
                            required={true}
                            placeholder={'Enter confluence base url'}
                            value={this.state.baseURL}
                            addValidation={this.validator.addValidation}
                            removeValidation={this.validator.removeValidation}
                            onChange={this.handleBaseURLChange}
                        />
                        <ConfluenceField
                            label={'SPACE KEY'}
                            type={'text'}
                            fieldType={'input'}
                            required={true}
                            placeholder={'Enter space key'}
                            value={this.state.spaceKey}
                            addValidation={this.validator.addValidation}
                            removeValidation={this.validator.removeValidation}
                            onChange={this.handleSpaceKey}
                        />
                        <ConfluenceField
                            isMulti={true}
                            label={'EVENTS'}
                            name={'events'}
                            fieldType={'dropDown'}
                            required={true}
                            theme={this.props.theme}
                            options={Constants.CONFLUENCE_EVENTS}
                            value={this.state.events}
                            addValidation={this.validator.addValidation}
                            removeValidation={this.validator.removeValidation}
                            onChange={this.handleEvents}
                        />
                    </div>
                </Modal.Body>
                <Modal.Footer>
                    <Button
                        type='button'
                        bsStyle='link'
                        onClick={this.handleClose}
                    >
                        {'Cancel'}
                    </Button>
                    <Button
                        type='submit'
                        bsStyle='primary'
                        onClick={this.handleSubmit}
                    >
                        {'Submit'}
                    </Button>
                </Modal.Footer>
            </Modal>
        );
    }
}
