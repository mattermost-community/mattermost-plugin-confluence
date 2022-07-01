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

export default function TokenForm(props: Props) {
    const handleURLChange = useCallback((e: React.ChangeEvent<typeof FormControl & HTMLInputElement>) => {
        props.setState({...props.state, serverURL: e.target.value});
        props.setErrors({...props.errors, serverURL: ''});
    }, [props]);

    const handleClientIDChange = useCallback((e: React.ChangeEvent<typeof FormControl & HTMLInputElement>) => {
        props.setState({...props.state, clientID: e.target.value});
        props.setErrors({...props.errors, clientID: ''});
    }, [props]);

    const handleClientSecretChange = useCallback((e: React.ChangeEvent<typeof FormControl & HTMLInputElement>) => {
        props.setState({...props.state, clientSecret: e.target.value});
        props.setErrors({...props.errors, clientSecret: ''});
    }, [props]);

    useEffect(() => {
        props.reset();
    }, []);

    return (
        <Form>
            <FormGroup>
                <label
                    className={'form-label ' + (props.errors.serverURL ? 'text-danger' : '')}
                >{'Server URL'}
                </label>
                <FormControl
                    className={props.errors.serverURL ? 'error' : ''}
                    type='text'
                    value={props.state.serverURL}
                    onChange={handleURLChange}
                    placeholder='<https://example.com>'
                />
                <small
                    className={props.errors.serverURL ? 'form-text text-danger' : ''}
                >{props.errors.serverURL && <p>{`* ${props.errors.serverURL}`}</p>}</small>
            </FormGroup>
            <FormGroup>
                <label
                    className={'form-label ' + (props.errors.clientID ? 'text-danger' : '')}
                >{'Client ID'}
                </label>
                <FormControl
                    className={props.errors.clientID ? 'error' : ''}
                    type='text'
                    value={props.state.clientID}
                    onChange={handleClientIDChange}
                    placeholder='<client-id>'
                />
                <small
                    className={props.errors.clientID ? 'form-text text-danger' : ''}
                >{props.errors.clientID && <p>{`* ${props.errors.clientID}`}</p>}</small>
            </FormGroup>
            <FormGroup>
                <label
                    className={'form-label ' + (props.errors.clientSecret ? 'text-danger' : '')}
                >{'Client Secret'}
                </label>
                <FormControl
                    className={props.errors.clientSecret ? 'error' : ''}
                    type='text'
                    value={props.state.clientSecret}
                    onChange={handleClientSecretChange}
                    placeholder='<client-secret>'
                />
                <small
                    className={props.errors.clientSecret ? 'form-text text-danger' : ''}
                >{props.errors.clientSecret && <p>{`* ${props.errors.clientSecret}`}</p>}</small>
            </FormGroup>
        </Form>
    );
}
