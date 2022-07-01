import React, {useState} from 'react';
import {Button, Modal} from 'react-bootstrap';

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

export default function TokenModal(props: Props) {
    const [state, setState] = useState(props.value);
    const [errors, setErrors] = useState({
        serverURL: '',
        clientID: '',
        clientSecret: '',
    });

    const reset = () => {
        setState(props.value);
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

        if (props.entryExists(state.serverURL)) {
            setErrors({...errors, serverURL: 'Server URL already exists'});
            return;
        }
        props.onSave(state);
    };

    const handleClose = () => {
        setErrors({serverURL: '', clientID: '', clientSecret: ''});
        props.handleClose();
    };
    return (
        <Modal
            show={props.open}
            onHide={handleClose}
        >
            <Modal.Header>
                <Modal.Title>{props.edit ? 'Update Confluence Config' : 'Add Confluence Config'}</Modal.Title>
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
                <Button onClick={handleClose}>{'Close'}</Button>
                <Button
                    onClick={onSubmit}
                    className='btn btn-primary'
                >{'Save'}</Button>
            </Modal.Footer>
        </Modal>
    );
}
