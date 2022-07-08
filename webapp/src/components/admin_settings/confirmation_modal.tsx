import React from 'react';
import {Button, Modal} from 'react-bootstrap';

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

export default function ConfirmationModal({open, title, body, confirmText,
    cancelText, approveButtonStyle, handleClose, handleConfirm}: Props) {
    return (
        <Modal
            show={open}
            onHide={handleClose}
        >
            <Modal.Header>
                <Modal.Title>{title}</Modal.Title>
            </Modal.Header>
            <Modal.Body>
                {body}
            </Modal.Body>
            <Modal.Footer>
                <Button onClick={handleClose}>{cancelText}</Button>
                <Button
                    onClick={handleConfirm}

                    // Removed "bsStyle" prop from here, as it was used in older versions of react-bootstrap
                    // and "variant" prop was also not working, so updated it with the className prop.
                    className={`btn btn-${approveButtonStyle}`}
                >
                    {confirmText}
                </Button>
            </Modal.Footer>
        </Modal>
    );
}
