import React, {useEffect, useCallback} from 'react';
import {Form, FormGroup, HelpBlock, ControlLabel, FormControl} from 'react-bootstrap';

import {ConfluenceConfig} from '../../types';

type Props = {
    state: ConfluenceConfig
    errors: ConfluenceConfig,
    setState: React.Dispatch<React.SetStateAction<ConfluenceConfig>>
    setErrors: React.Dispatch<React.SetStateAction<ConfluenceConfig>>
    reset: () => void
}

export default function TokenForm({state, errors, setState, setErrors, reset}: Props) {
    const handleURLChange = useCallback((e: React.ChangeEvent<FormControl & HTMLInputElement>) => {
        setState({...state, serverURL: e.target.value});
        setErrors({...errors, serverURL: ''});
    }, [errors, setErrors, setState, state]);

    const handleClientIDChange = useCallback((e: React.ChangeEvent<FormControl & HTMLInputElement>) => {
        setState({...state, clientID: e.target.value});
        setErrors({...errors, clientID: ''});
    }, [errors, setErrors, setState, state]);

    const handleClientSecretChange = useCallback((e: React.ChangeEvent<FormControl & HTMLInputElement>) => {
        setState({...state, clientSecret: e.target.value});
        setErrors({...errors, clientSecret: ''});
    }, [errors, setErrors, setState, state]);

    useEffect(() => {
        reset();
    }, []);

    return (
        <Form>
            <FormGroup validationState={errors.serverURL ? 'error' : null}>
                <ControlLabel>
                    {'Server URL'}
                </ControlLabel>
                <FormControl
                    type='text'
                    value={state.serverURL}
                    onChange={handleURLChange}
                    placeholder='<https://example.com>'
                />
                <HelpBlock>{errors.serverURL && <p>{`* ${errors.serverURL}`}</p>}</HelpBlock>
            </FormGroup>
            <FormGroup validationState={errors.clientID ? 'error' : null}>
                <ControlLabel >
                    {'Client ID'}
                </ControlLabel>
                <FormControl
                    type='text'
                    value={state.clientID}
                    onChange={handleClientIDChange}
                    placeholder='<client-id>'
                />
                <HelpBlock>{errors.clientID && <p>{`* ${errors.clientID}`}</p>}</HelpBlock>
            </FormGroup>
            <FormGroup validationState={errors.clientSecret ? 'error' : null}>
                <ControlLabel >
                    {'Client Secret'}
                </ControlLabel>
                <FormControl
                    type='text'
                    value={state.clientSecret}
                    onChange={handleClientSecretChange}
                    placeholder='<client-secret>'
                />
                <HelpBlock>{errors.clientSecret && <p>{`* ${errors.clientSecret}`}</p>}</HelpBlock>
            </FormGroup>
        </Form>
    );
}

