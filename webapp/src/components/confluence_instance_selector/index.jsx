import React, {useCallback} from 'react';
import {useSelector} from 'react-redux';

import PropTypes from 'prop-types';

import Validator from '../validator';
import ReactSelectSetting from '../react_select_setting';
import selectors from '../../selectors';

const ConfluenceInstanceSelector = (props) => {
    const validator = new Validator();

    const installedInstances = useSelector((state) => selectors.isInstalledInstances(state));

    const getInstanceOptions = useCallback(() => {
        return installedInstances?.map((instance) => ({label: instance.instance_id, value: instance.instance_id})); 
      }, [installedInstances])

    const handleEvents = useCallback((_, instanceID) => {
        if (instanceID === props.selectedInstanceID) {
            return;
        }
        props.onInstanceChange(instanceID);
      }, [props.selectedInstanceID])


    return (
        <React.Fragment>
            <ReactSelectSetting
                name={'instance'}
                label={'Instance'}
                options={getInstanceOptions()}
                onChange={handleEvents}
                value={getInstanceOptions().find((option) => option.value === props.selectedInstanceID)}
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
