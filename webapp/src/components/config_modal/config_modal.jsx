import React from 'react';
import {
    Modal,
    Button,
} from 'react-bootstrap';

import PropTypes from 'prop-types';

import Constants from '../../constants';
import ConfluenceField from '../confluence_field';
import Validator from '../validator';

const initialState = {
    alias: '',
    baseURL: '',
    spaceKey: '',
    events: [],
    error: '',
    saving: false,
};

export default class ConfigModal extends React.PureComponent {
    static propTypes = {
        visibility: PropTypes.bool.isRequired,
        close: PropTypes.func.isRequired,
        theme: PropTypes.object.isRequired,
        saveChannelSubscription: PropTypes.func.isRequired,
        currentChannelID: PropTypes.string.isRequired,
    };

    constructor(props) {
        super(props);
        this.state = initialState;
        this.validator = new Validator();
    }

    handleClose = (e) => {
        if (e && e.preventDefault) {
            e.preventDefault();
        }
        this.setState(initialState, this.props.close);
    };

    handleAlias = (e) => {
        this.setState({
            alias: e.target.value,
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

    handleSubmit = async () => {
        if (!this.validator.validate()) {
            return;
        }
        const {alias, baseURL, spaceKey, events} = this.state;
        const channelSubscription = {
            alias,
            baseURL,
            spaceKey,
            channelID: this.props.currentChannelID,
            events: events.map((event) => event.value),
        };
        this.setState({saving: true});
        const {error} = await this.props.saveChannelSubscription(channelSubscription);
        if (error) {
            this.setState({
                error: 'error occurred',
                saving: false,
            });
            return;
        }
        this.handleClose();
    };

    render() {
        const {visibility} = this.props;
        const {error, saving} = this.state;
        let createError = null;
        if (error) {
            createError = <span className='error'>{error}</span>;
        }

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
                            label={'Alias'}
                            type={'text'}
                            fieldType={'input'}
                            required={true}
                            placeholder={'Enter alias for this subscription'}
                            value={this.state.alias}
                            addValidation={this.validator.addValidation}
                            removeValidation={this.validator.removeValidation}
                            onChange={this.handleAlias}
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
                            theme={this.props.theme}
                            options={Constants.CONFLUENCE_EVENTS}
                            value={this.state.events}
                            addValidation={this.validator.addValidation}
                            removeValidation={this.validator.removeValidation}
                            onChange={this.handleEvents}
                        />
                        {createError}
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
                        disabled={saving}
                    >
                        {saving && <span className='fa fa-spinner fa-fw fa-pulse spinner'/>}
                        {'Save Subscription'}
                    </Button>
                </Modal.Footer>
            </Modal>
        );
    }
}
