import React, {useEffect, useState, useCallback, SyntheticEvent, useMemo} from 'react';
import {DefaultRootState, useDispatch, useSelector} from 'react-redux';
import {Modal, Button} from 'react-bootstrap';

import {getCurrentChannelId} from 'mattermost-redux/selectors/entities/common';

import {Theme} from 'mattermost-redux/types/preferences';

import selectors from 'src/selectors';

import ConfluenceInstanceSelector from 'src/components/confluence_instance_selector';
import ConfluenceSpaceSelector from 'src/components/confluence_space_selector';
import Validator from 'src/components/validator';
import ConfluenceField from 'src/components/confluence_field';
import {getModalStyles} from 'src/utils/styles';
import {
    getSpacesForConfluenceURL,
    createPageForConfluence,
    closeCreateConfluencePageModal,
} from 'src/actions';

const getStyle = () => ({
    typeFormControl: {
        resize: 'none',
        height: '10em',
    },
});

const CreateConfluencePage = (theme: Theme) => {
    const dispatch = useDispatch();
    const validator = new Validator();

    const postMessage = useSelector((state:DefaultRootState) => selectors.postMessage(state));
    const channelID = useSelector((state:DefaultRootState) => getCurrentChannelId(state));

    const [modalVisible, setModalVisible] = useState<boolean>(false);
    const [instanceID, setInstanceID] = useState<string>('');
    const [pageTitle, setPageTitle] = useState<string>('');
    const [pageDescription, setPageDescription] = useState<string>(
        postMessage?.message,
    );
    const [spaceKey, setSpaceKey] = useState<string>('');
    const [saving, setSaving] = useState<boolean>(false);
    const [error, setError] = useState<string>('');

    useEffect(() => {
        if (postMessage?.message) {
            setModalVisible(true);
            setPageDescription(postMessage.message);
        } else {
            setModalVisible(false);
        }
    }, [postMessage]);

    useEffect(() => {
        if (!instanceID) {
            return;
        }
        (async () => {
            const response = await getSpacesForConfluenceURL(instanceID)(dispatch);
            if (response?.error) {
                setError(response.error.response?.text);
            }
        })();
        setSpaceKey('');
    }, [instanceID]);

    const reset = () => {
        setSaving(false);
        setInstanceID('');
        setSpaceKey('');
        setPageTitle('');
        setPageDescription('');
        setError('');
    };

    const handleClose = useCallback((e:PointerEvent) => {
        if (e?.preventDefault) {
            e.preventDefault();
        }
        reset();
        dispatch(closeCreateConfluencePageModal());
    }, []);

    const handleInstanceChange = useCallback(
        (currentInstanceID:string) => {
            setInstanceID(currentInstanceID);
            setSpaceKey('');
            setError('');
        }, []);

    const handleSpaceKeyChange = useCallback(
        (currentSpaceKey:string) => {
            setSpaceKey(currentSpaceKey);
        }, []);

    const handlePageTitle = useCallback(
        (e:SyntheticEvent) => {
            e.persist();
            setPageTitle(e.target.value);
        }, []);

    const handlePageDescription = useCallback(
        (e:SyntheticEvent) => {
            e.persist();
            setPageDescription(e.target.value);
        }, []);

    const handleSubmit = () => {
        if (!validator.validate()) {
            return;
        }

        const pageDetails = {
            title: pageTitle,
            description: pageDescription,
        };

        setSaving(true);
        (async () => {
            const response = await createPageForConfluence(
                instanceID,
                channelID,
                spaceKey,
                pageDetails,
            )(dispatch);
            if (response?.error) {
                setError(response.error.response?.text);
                setSaving(false);
            } else {
                reset();
                dispatch(closeCreateConfluencePageModal());
            }
        })();
    };

    return (
        <Modal
            dialogClassName='modal--scroll'
            show={modalVisible}
            onHide={handleClose}
            onExited={handleClose}
            backdrop={'static'}
            bsSize='large'
        >
            <Modal.Header closeButton={true}>
                <Modal.Title>{'Create Confluence Page'}</Modal.Title>
            </Modal.Header>
            <Modal.Body style={getModalStyles.modalBody}>
                <ConfluenceInstanceSelector
                    theme={theme}
                    selectedInstanceID={instanceID}
                    onInstanceChange={handleInstanceChange}
                />

                {instanceID && (
                    <ConfluenceSpaceSelector
                        theme={theme}
                        selectedSpaceKey={spaceKey}
                        onSpaceKeyChange={handleSpaceKeyChange}
                    />
                )}

                {spaceKey && (
                    <ConfluenceField
                        label={'Page Title'}
                        type={'text'}
                        fieldType={'input'}
                        required={true}
                        placeholder={'Enter Page Title.'}
                        value={pageTitle}
                        addValidation={validator.addComponent}
                        removeValidation={validator.removeComponent}
                        onChange={handlePageTitle}
                    />
                )}
                {spaceKey && (
                    <ConfluenceField
                        label={'Page Description'}
                        formControlStyle={getStyle().typeFormControl}
                        type={'textarea'}
                        fieldType={'input'}
                        required={true}
                        value={pageDescription}
                        addValidation={validator.addComponent}
                        removeValidation={validator.removeComponent}
                        onChange={handlePageDescription}
                    />
                )}
                {error && (
                    <p className='alert alert-danger'>
                        <i
                            className='fa fa-warning'
                            title='Warning Icon'
                        />
                        <span> {error}</span>
                    </p>
                )}
            </Modal.Body>

            {spaceKey && (
                <Modal.Footer style={getModalStyles.modalFooter}>
                    <Button
                        type='button'
                        bsStyle='link'
                        onClick={handleClose}
                    >
                        {'Cancel'}
                    </Button>
                    <Button
                        type='submit'
                        bsStyle='primary'
                        onClick={handleSubmit}
                        disabled={saving}
                    >
                        {saving && (
                            <span className='fa fa-spinner fa-fw fa-pulse spinner'/>
                        )}
                        {'Save Subscription'}
                    </Button>
                </Modal.Footer>
            )}
        </Modal>
    );
};

export default CreateConfluencePage;
