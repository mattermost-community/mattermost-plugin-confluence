import React from 'react';
import {
    FormGroup,
    FormControl,
    Modal,
    ControlLabel,
    Button,
} from 'react-bootstrap';

import PropTypes from 'prop-types';
import Select from 'react-select';

import Constants from '../../constants';
import {getStyleForReactSelect} from '../react_select_settings';

export default class ConfigModal extends React.PureComponent {
    static propTypes = {
        visibility: PropTypes.bool.isRequired,
        close: PropTypes.func.isRequired,
        theme: PropTypes.object.isRequired,
    };

    handleClose = () => {
        this.props.close();
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
                        <FormGroup>
                            <ControlLabel>{'CONFLUENCE BASE URL'}</ControlLabel>
                            <FormControl
                                type={'text'}
                                placeholder={'Enter confluence base url'}
                            />
                        </FormGroup>
                        <FormGroup>
                            <ControlLabel>{'SPACE KEY'}</ControlLabel>
                            <FormControl
                                type={'text'}
                                placeholder={'Enter space key'}
                            />
                        </FormGroup>
                        <FormGroup>
                            <ControlLabel>{'EVENTS'}</ControlLabel>
                            <Select
                                isMulti={'true'}
                                name={'events'}
                                options={Constants.CONFLUENCE_EVENTS}
                                menuPortalTarget={document.body}
                                menuPlacement='auto'
                                styles={getStyleForReactSelect(this.props.theme)}
                            />
                        </FormGroup>
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
                    >
                        {'Submit'}
                    </Button>
                </Modal.Footer>
            </Modal>
        );
    }
}
