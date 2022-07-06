import React, {useEffect, useCallback} from 'react';
import {Form, FormControl, FormGroup} from 'react-bootstrap';

import {ConfluenceConfig} from '../../types';

import './style.scss';

type Props = {
    state: ConfluenceConfig
    errors: ConfluenceConfig,
    setState: React.Dispatch<React.SetStateAction<ConfluenceConfig>>
    setErrors: React.Dispatch<React.SetStateAction<ConfluenceConfig>>
    reset: () => void
}

export default function TokenForm({state, errors, setState, setErrors, reset}: Props) {
    const handleURLChange = useCallback((e: React.ChangeEvent<typeof FormControl & HTMLInputElement>) => {
        setState({...state, serverURL: e.target.value});
        setErrors({...errors, serverURL: ''});
    }, [errors, setErrors, setState, state]);

    const handleClientIDChange = useCallback((e: React.ChangeEvent<typeof FormControl & HTMLInputElement>) => {
        setState({...state, clientID: e.target.value});
        setErrors({...errors, clientID: ''});
    }, [errors, setErrors, setState, state]);

    const handleClientSecretChange = useCallback((e: React.ChangeEvent<typeof FormControl & HTMLInputElement>) => {
        setState({...state, clientSecret: e.target.value});
        setErrors({...errors, clientSecret: ''});
    }, [errors, setErrors, setState, state]);

    useEffect(() => {
        reset();
    }, []);

    return (
        <Form>
            <FormGroup>
                <label
                    className={'form-label ' + (errors.serverURL ? 'text-danger' : '')}
                >{'Server URL'}
                </label>
                <FormControl
                    className={errors.serverURL ? 'error' : ''}
                    type='text'
                    value={state.serverURL}
                    onChange={handleURLChange}
                    placeholder='<https://example.com>'
                />
                <small
                    className={errors.serverURL ? 'form-text text-danger' : ''}
                >{errors.serverURL && <p>{`* ${errors.serverURL}`}</p>}</small>
            </FormGroup>
            <FormGroup>
                <label
                    className={'form-label ' + (errors.clientID ? 'text-danger' : '')}
                >{'Client ID'}
                </label>
                <FormControl
                    className={errors.clientID ? 'error' : ''}
                    type='text'
                    value={state.clientID}
                    onChange={handleClientIDChange}
                    placeholder='<client-id>'
                />
                <small
                    className={errors.clientID ? 'form-text text-danger' : ''}
                >{errors.clientID && <p>{`* ${errors.clientID}`}</p>}</small>
            </FormGroup>
            <FormGroup>
                <label
                    className={'form-label ' + (errors.clientSecret ? 'text-danger' : '')}
                >{'Client Secret'}
                </label>
                <FormControl
                    className={errors.clientSecret ? 'error' : ''}
                    type='text'
                    value={state.clientSecret}
                    onChange={handleClientSecretChange}
                    placeholder='<client-secret>'
                />
                <small
                    className={errors.clientSecret ? 'form-text text-danger' : ''}
                >{errors.clientSecret && <p>{`* ${errors.clientSecret}`}</p>}</small>
            </FormGroup>
        </Form>
    );
}
