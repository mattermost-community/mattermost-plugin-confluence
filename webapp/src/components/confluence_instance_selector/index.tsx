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

const ConfluenceInstanceSelector = (props: Props) => {
    const validator = new Validator();

    const installedInstances = useSelector((state: DefaultRootState) =>
        selectors.isInstalledInstances(state),
    );

    const getInstanceOptions = useMemo(() => {
        return installedInstances?.map((instance: { instance_id: string; }) => ({
            label: instance.instance_id,
            value: instance.instance_id,
        }));
    }, [installedInstances]);

    const handleEvents = useCallback(
        (_, instanceID) => {
            if (instanceID !== props.selectedInstanceID) {
                props.onInstanceChange(instanceID);
            }
        },
        [props.selectedInstanceID],
    );

    return (
        <React.Fragment>
            <ReactSelectSetting
                name={'instance'}
                label={'Instance'}
                options={getInstanceOptions}
                onChange={handleEvents}
                value={getInstanceOptions.find(
                    (option: { value: string; }) => option.value === props.selectedInstanceID,
                )}
                required={true}
                theme={props.theme}
                addValidate={validator.addComponent}
                removeValidate={validator.removeComponent}
            />
        </React.Fragment>
    );
};

export default ConfluenceInstanceSelector;
