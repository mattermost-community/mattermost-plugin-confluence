import React, {useEffect, useState} from 'react';
import {useSelector} from 'react-redux';

import PropTypes from 'prop-types';

import Validator from '../validator';
import ReactSelectSetting from '../react_select_setting';
import selectors from 'src/selectors';

const ConfluenceInstanceSelector = (props) => {
    const validator = new Validator();

    const isInstalledInstances = useSelector((state) => selectors.isInstalledInstances(state));

    const [installedInstancesOptions, setInstalledInstancesOptions] = useState([]);

    useEffect(() => {
        const installedInstancesSelectOptions = isInstalledInstances?.map((instance) => ({label: instance.instance_id, value: instance.instance_id}));
        setInstalledInstancesOptions(installedInstancesSelectOptions);
    }, [isInstalledInstances]);

    const handleEvents = (_, instanceID) => {
        if (instanceID === props.selectedInstanceID) {
            return;
        }
        props.onInstanceChange(instanceID);
    };

    return (
        <React.Fragment>
            <ReactSelectSetting
                name={'instance'}
                label={'Instance'}
                options={installedInstancesOptions}
                onChange={handleEvents}
                value={installedInstancesOptions.find((option) => option.value === props.selectedInstanceID)}
                required={true}
                theme={props.theme}
                addValidate={validator.addComponent}
                removeValidate={validator.removeComponent}
            />
        </React.Fragment>
    );
};

ConfluenceInstanceSelector.propTypes = {
    theme: PropTypes.object.isRequired,
    selectedInstanceID: PropTypes.string.isRequired,
    onInstanceChange: PropTypes.func.isRequired,
};

export default ConfluenceInstanceSelector;
