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

export default class SubscriptionModal extends React.PureComponent {
    static propTypes = {
        visibility: PropTypes.bool,
        editSubscription: PropTypes.bool,
        subscription: PropTypes.object,
        alias: PropTypes.string,
        baseURL: PropTypes.string,
        spaceKey: PropTypes.string,
        close: PropTypes.func,
        closeEditSubscription: PropTypes.func,
        theme: PropTypes.object.isRequired,
        saveChannelSubscription: PropTypes.func.isRequired,
        currentChannelID: PropTypes.string.isRequired,
    };

    static defaultProps = {
        visibility: false,
        editSubscription: false,
        subscription: {},
    };

    constructor(props) {
        super(props);
        this.state = initialState;
        this.validator = new Validator();
    }

    componentDidUpdate(prevProps) {
        if (this.props.subscription !== prevProps.subscription) {
            this.setData();
        }
    }

    setData = () => {
        const {alias, baseURL, spaceKey, events} = this.props.subscription;
        if (alias) {
            this.setState({
                alias,
                baseURL,
                spaceKey,
                events: Constants.CONFLUENCE_EVENTS.filter((option) => events.includes(option.value)),
            });
        }
    };

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
        const {currentChannelID, subscription} = this.props;
        const channelSubscription = {
            alias,
            baseURL,
            spaceKey,
            channelID: currentChannelID,
            events: events.map((event) => event.value),
        };
        this.setState({saving: true});
        if (subscription && subscription.alias) {
            // TODO : ADD logic to edit subscription
        } else {
            const {error} = await this.props.saveChannelSubscription(channelSubscription);
            if (error) {
                this.setState({
                    error: error.response.text,
                    saving: false,
                });
                return;
            }
        }
        this.handleClose();
    };

    render() {
        const {visibility, subscription} = this.props;
        const editSubscription = Boolean(subscription && subscription.alias);
        const isModalVisible = Boolean(visibility || editSubscription);
        const {error, saving} = this.state;
        let createError = null;
        if (error) {
            createError = (
                <p className='alert alert-danger'>
                    <i
                        className='fa fa-warning'
                        title='Warning Icon'
                    />
                    <span> {error}</span>
                </p>
            );
        }

        return (
            <Modal
                show={isModalVisible}
                onHide={this.handleClose}
                backdrop={'static'}
            >
                <Modal.Header closeButton={true}>
                    <Modal.Title>
                        {'Channel Confluence Settings'}
                    </Modal.Title>
                </Modal.Header>
                <Modal.Body>
                    <div>
                        <ConfluenceField
                            label={'Alias'}
                            type={'text'}
                            fieldType={'input'}
                            required={true}
                            readOnly={editSubscription}
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
                            readOnly={editSubscription}
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
                            readOnly={editSubscription}
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
                            required={false}
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
