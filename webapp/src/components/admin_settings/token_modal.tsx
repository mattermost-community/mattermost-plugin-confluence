import React, {useState} from 'react';
import {Modal, Button} from 'react-bootstrap';

import {ConfluenceConfig} from '../../types';

import TokenForm from './token_form';

type Props = {
    open: boolean,
    edit: boolean,
    value: ConfluenceConfig,
    handleClose: () => void,
    onSave: (data: ConfluenceConfig) => void
    entryExists: (name: string) => boolean
}

export default function TokenModal({open, edit, value, handleClose, onSave, entryExists}: Props) {
    const [state, setState] = useState(value);
    const [errors, setErrors] = useState({
        serverURL: '',
        clientID: '',
        clientSecret: '',
    });

    const reset = () => {
        setState(value);
    };

    const onSubmit = () => {
        if (state.serverURL === '' || state.clientID === '' || state.clientSecret === '') {
            const newErrors = {...errors};
            if (state.serverURL === '') {
                newErrors.serverURL = 'This field is required';
            }
            if (state.clientID === '') {
                newErrors.clientID = 'This field is required';
            }
            if (state.clientSecret === '') {
                newErrors.clientSecret = 'This field is required';
            }
            setErrors(newErrors);
            return;
        }

        if (entryExists(state.serverURL)) {
            setErrors({...errors, serverURL: 'Server URL already exists'});
            return;
        }
        onSave(state);
    };

    const closeHandler = () => {
        setErrors({serverURL: '', clientID: '', clientSecret: ''});
        handleClose();
    };
    return (
        <Modal
            show={open}
            onHide={closeHandler}
        >
            <Modal.Header>
                <Modal.Title>{edit ? 'Update Confluence Config' : 'Add Confluence Config'}</Modal.Title>
            </Modal.Header>
            <Modal.Body>
                <TokenForm
                    state={state}
                    setState={setState}
                    errors={errors}
                    setErrors={setErrors}
                    reset={reset}
                />
            </Modal.Body>
            <Modal.Footer>
                <Button onClick={closeHandler}>{'Close'}</Button>
                <Button
                    onClick={onSubmit}
                    bsStyle='primary'
                >
                    {'Save'}
                </Button>
            </Modal.Footer>
        </Modal>
    );
}
