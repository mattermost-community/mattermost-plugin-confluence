import React from 'react';
import {Modal, Button} from 'react-bootstrap';

type Props = {
    open: boolean,
    title: string,
    body: string,
    confirmText: string,
    cancelText: string,
    approveButtonStyle: string,
    handleClose: () => void,
    handleConfirm: () => void
}

export default function ConfirmationModal(props: Props) {
    return (
        <Modal
            show={props.open}
            onHide={props.handleClose}
        >
            <Modal.Header>
                <Modal.Title>{props.title}</Modal.Title>
            </Modal.Header>
            <Modal.Body>
                {props.body}
            </Modal.Body>
            <Modal.Footer>
                <Button onClick={props.handleClose}>{props.cancelText}</Button>
                <Button
                    onClick={props.handleConfirm}
                    bsStyle={props.approveButtonStyle}
                >{props.confirmText}</Button>
            </Modal.Footer>
        </Modal>
    );
}
