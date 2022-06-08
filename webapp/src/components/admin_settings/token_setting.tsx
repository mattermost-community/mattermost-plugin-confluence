import React, {useState} from 'react';
import {FormGroup, Col, Table, Button, ControlLabel} from 'react-bootstrap';

import {ConfluenceConfig} from '../../types';

import TokenModal from './token_modal';
import ConfirmationModal from './confirmation_modal';

import './style.scss';

type Props = {
    id: string,
    label: string,
    value: ConfluenceConfig[],
    onChange: (id: string, value: ConfluenceConfig[]) => void
}

const defaultValue: ConfluenceConfig = {
    serverURL: '',
    clientID: '',
    clientSecret: '',
};

export default function TokenSetting(props: Props) {
    const [state, setState] = useState({
        values: props.value,
    });
    const [editModalState, setEditModalState] = useState({open: false, edit: false, value: defaultValue, index: 0});
    const [deleteModalState, setDeleteModalState] = useState({open: false, index: 0});

    const addTokenEntry = (e: { preventDefault: () => void; }) => {
        e.preventDefault();
        setEditModalState({
            open: true,
            edit: false,
            value: defaultValue,
            index: state.values.length + 1,
        });
    };

    const handleSave = (idx: number, value: ConfluenceConfig) => {
        const resultArr = [...props.value];
        if (idx >= resultArr.length) {
            resultArr.push(value);
        } else {
            resultArr[idx] = value;
        }
        props.onChange(props.id, resultArr);
        setState({values: resultArr});
        setEditModalState({
            open: false,
            edit: false,
            value: defaultValue,
            index: 0,
        });
    };

    const handleDelete = (idx: number) => {
        const resultArr = [...props.value];
        resultArr.splice(idx, 1);
        props.onChange(props.id, resultArr);
        setState({
            values: resultArr,
        });
        setDeleteModalState({open: false, index: 0});
    };

    const handleSelect = (idx: number, value: ConfluenceConfig) => {
        setEditModalState({
            open: true,
            edit: true,
            value: {...editModalState.value, ...value},
            index: idx,
        });
    };

    const handleDeleteClick = (idx: number) => {
        setDeleteModalState({open: true, index: idx});
    };

    const handleDeleteClose = () => {
        setDeleteModalState({open: false, index: 0});
    };

    const entryExists = (serverURL: string): boolean => {
        return Boolean(state.values.find((value: ConfluenceConfig, idx: number) => idx !== editModalState.index && value.serverURL.toLowerCase() === serverURL.toLowerCase()));
    };

    return (
        <FormGroup>
            <Col
                componentClass={ControlLabel}
                sm={4}
            >{props.label}</Col>
            <Col sm={8}>
                <Table
                    className='table'
                    striped={true}
                    bordered={true}
                    condensed={true}
                    hover={true}
                >
                    <thead>
                        <tr>
                            <th colSpan={4}>{'Server URL'}</th>
                            <th colSpan={4}>{'ClientID'}</th>
                            <th colSpan={4}>{'ClientSecret'}</th>
                            <th className='action-column'>{'Actions'}</th>
                        </tr>
                    </thead>
                    <tbody>
                        {(state.values.length > 0) ?
                            <>
                                {state.values.map((val, idx) => {
                                    return (
                                        <tr key={idx}>
                                            <td
                                                className={'table-data'}
                                                colSpan={4}
                                            >{val.serverURL}</td>
                                            <td
                                                className={'table-data'}
                                                colSpan={4}
                                            >{val.clientID}</td>
                                            <td
                                                className={'table-data'}
                                                colSpan={4}
                                            >{val.clientSecret}</td>
                                            <td style={{whiteSpace: 'nowrap'}}>
                                                <Button
                                                    bsClass='btn transparent-button btn-default'
                                                    onClick={() => handleSelect(idx, val)}
                                                >
                                                    <i className='button-icon fa fa-edit'/>
                                                </Button>
                                                <Button
                                                    bsClass='btn transparent-button btn-default'
                                                    onClick={() => handleDeleteClick(idx)}
                                                >
                                                    <i className='button-icon fa fa-trash'/>
                                                </Button>
                                            </td>
                                        </tr>
                                    );
                                })}
                            </> :
                            <tr>
                                <td
                                    colSpan={14}
                                    className={'no-token-bs-class'}
                                >
                                    <div className={'no-token-content'}>
                                        {'No Config Found.'}
                                    </div>
                                </td>
                            </tr>
                        }
                    </tbody>
                </Table>
                <ConfirmationModal
                    open={deleteModalState.open}
                    title={'Delete Confluence Config'}
                    body={deleteModalState.open ? `Delete the config for "${state.values[deleteModalState.index]?.serverURL}"?` : ''}
                    confirmText={'Delete'}
                    cancelText={'Cancel'}
                    approveButtonStyle={'danger'}
                    handleClose={handleDeleteClose}
                    handleConfirm={() => {
                        handleDelete(deleteModalState.index);
                    }}
                />
                <TokenModal
                    value={editModalState.value}
                    edit={editModalState.edit}
                    open={editModalState.open}
                    handleClose={() => {
                        setEditModalState({...editModalState, open: false});
                    }}
                    onSave={(values: any) => handleSave(editModalState.index, values)}
                    entryExists={entryExists}
                />
                <div style={{marginTop: '20px'}}>
                    <Button
                        onClick={addTokenEntry}
                    >
                        {'Add Config'}
                    </Button>
                </div>
            </Col>
        </FormGroup>
    );
}
