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

export default function TokenForm(props: Props) {
    const handleURLChange = useCallback((e: React.ChangeEvent<FormControl & HTMLInputElement>) => {
        props.setState({...props.state, serverURL: e.target.value});
        props.setErrors({...props.errors, serverURL: ''});
    }, [props]);

    const handleClientIDChange = useCallback((e: React.ChangeEvent<FormControl & HTMLInputElement>) => {
        props.setState({...props.state, clientID: e.target.value});
        props.setErrors({...props.errors, clientID: ''});
    }, [props]);

    const handleClientSecretChange = useCallback((e: React.ChangeEvent<FormControl & HTMLInputElement>) => {
        props.setState({...props.state, clientSecret: e.target.value});
        props.setErrors({...props.errors, clientSecret: ''});
    }, [props]);

    useEffect(() => {
        props.reset();
    }, []);

    return (
        <Form>
            <FormGroup validationState={props.errors.serverURL ? 'error' : null}>
                <ControlLabel>
                    {'Server URL'}
                </ControlLabel>
                <FormControl
                    type='text'
                    value={props.state.serverURL}
                    onChange={handleURLChange}
                    placeholder='<https://example.com>'
                />
                <HelpBlock>{props.errors.serverURL && <p>{`* ${props.errors.serverURL}`}</p>}</HelpBlock>
            </FormGroup>
            <FormGroup validationState={props.errors.clientID ? 'error' : null}>
                <ControlLabel >
                    {'Client ID'}
                </ControlLabel>
                <FormControl
                    type='text'
                    value={props.state.clientID}
                    onChange={handleClientIDChange}
                    placeholder='<client-id>'
                />
                <HelpBlock>{props.errors.clientID && <p>{`* ${props.errors.clientID}`}</p>}</HelpBlock>
            </FormGroup>
            <FormGroup validationState={props.errors.clientSecret ? 'error' : null}>
                <ControlLabel >
                    {'Client Secret'}
                </ControlLabel>
                <FormControl
                    type='text'
                    value={props.state.clientSecret}
                    onChange={handleClientSecretChange}
                    placeholder='<client-secret>'
                />
                <HelpBlock>{props.errors.clientSecret && <p>{`* ${props.errors.clientSecret}`}</p>}</HelpBlock>
            </FormGroup>
        </Form>
    );
}

