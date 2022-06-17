import React, {useEffect, useState} from 'react';
import {useSelector} from 'react-redux';

import PropTypes from 'prop-types';

import Validator from '../validator';
import ReactSelectSetting from '../react_select_setting';
import selectors from 'src/selectors';

const ConfluenceSpaceSelector = (props) => {
    const validator = new Validator();

    const spacesForConfluenceURL = useSelector((state) => selectors.spacesForConfluenceURL(state));

    const [spaceOptions, setSpaceOptions] = useState([]);

    useEffect(() => {
        if (spacesForConfluenceURL && spacesForConfluenceURL?.spaces) {
            const issueOptions = spacesForConfluenceURL?.spaces.map((it) => ({label: it.name, value: it.key}));
            setSpaceOptions(issueOptions);
        }
    }, [spacesForConfluenceURL]);

    const handleEvents = (_, spaceKey) => {
        if (spaceKey === props.selectedSpaceKey) {
            return;
        }
        props.onSpaceKeyChange(spaceKey);
    };

    return (
        <React.Fragment>
            <ReactSelectSetting
                name={'space'}
                label={'Space'}
                options={spaceOptions}
                onChange={handleEvents}
                value={props.selectedSpaceKey === '' ? null : spaceOptions.find((option) => option.value === props.selectedSpaceKey)}
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
