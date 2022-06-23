import React, {useCallback} from 'react';
import {useSelector} from 'react-redux';

import PropTypes from 'prop-types';

import Validator from '../validator';
import ReactSelectSetting from '../react_select_setting';
import selectors from '../../selectors';

const ConfluenceSpaceSelector = (props) => {
    const validator = new Validator();

    const spacesForConfluenceURL = useSelector((state) => selectors.spacesForConfluenceURL(state));

    const getSpaceOptions = useCallback(() => {
        return spacesForConfluenceURL?.spaces?.map((space) => ({label: space.name, value: space.key}));
    }, [spacesForConfluenceURL]);

    const handleEvents = useCallback((_, spaceKey) => {
        if (spaceKey !== props.selectedSpaceKey) {
            props.onSpaceKeyChange(spaceKey);
        }
    }, [props.selectedSpaceKey]);

    return (
        <React.Fragment>
            <ReactSelectSetting
                name={'space'}
                label={'Space'}
                options={getSpaceOptions()}
                onChange={handleEvents}
                value={props.selectedSpaceKey ? getSpaceOptions().find((option) => option.value === props.selectedSpaceKey) : null}
                required={true}
                theme={props.theme}
                addValidate={validator.addComponent}
                removeValidate={validator.removeComponent}
            />
        </React.Fragment>
    );
};

ConfluenceSpaceSelector.propTypes = {
    theme: PropTypes.object.isRequired,
    selectedSpaceKey: PropTypes.string.isRequired,
    onSpaceKeyChange: PropTypes.func.isRequired,
};

export default ConfluenceSpaceSelector;
