import React, {useCallback, useMemo} from 'react';
import {DefaultRootState, useSelector} from 'react-redux';

import {Theme} from 'mattermost-redux/types/preferences';

import Validator from 'src/components/validator';
import ReactSelectSetting from 'src/components/react_select_setting';
import selectors from 'src/selectors';

type Props = {
    selectedInstanceID: string;
    onInstanceChange: (currentInstanceID: string) => void;
    theme: Theme;
};

const ConfluenceInstanceSelector = ({selectedInstanceID, onInstanceChange, theme}: Props) => {
    const validator = useMemo(() => (new Validator()), []);

    const installedInstances = useSelector((state: DefaultRootState) =>
        selectors.installedInstances(state),
    );

    const getInstanceOptions = useMemo(() => (
        installedInstances?.map((instance: {instance_id: string}) => ({
            label: instance.instance_id,
            value: instance.instance_id,
        }))), [installedInstances]);

    const handleEvents = useCallback((_, instanceID) => {
        if (instanceID !== selectedInstanceID) {
            onInstanceChange(instanceID);
        }
    }, [selectedInstanceID, onInstanceChange],
    );

    return (
        <>
            <ReactSelectSetting
                name={'instance'}
                label={'Instance'}
                options={getInstanceOptions}
                onChange={handleEvents}
                value={getInstanceOptions.find((option: {value: string}) => option.value === selectedInstanceID)}
                required={true}
                theme={theme}
                addValidate={validator.addComponent}
                removeValidate={validator.removeComponent}
            />
        </>
    );
};

export default ConfluenceInstanceSelector;
