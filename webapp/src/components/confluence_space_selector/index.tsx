import React, {useCallback, useMemo} from 'react';
import {DefaultRootState, useSelector} from 'react-redux';

import PropTypes from 'prop-types';

import {Theme} from 'mattermost-redux/types/preferences';

import Validator from 'src/components/validator';
import ReactSelectSetting from 'src/components/react_select_setting';
import selectors from 'src/selectors';

type Props = {
    selectedSpaceKey: string;
    onSpaceKeyChange: (currentSpaceKey: string) => void;
    theme: Theme;
};

const ConfluenceSpaceSelector = ({selectedSpaceKey, onSpaceKeyChange, theme}: Props) => {
    const validator = new Validator();

    const spacesForConfluenceURL = useSelector((state: DefaultRootState) =>
        selectors.spacesForConfluenceURL(state),
    );

    const getSpaceOptions = useMemo(() => {
        return spacesForConfluenceURL?.spaces?.map((space: {name: string, key: string}) => ({
            label: space.name,
            value: space.key,
        }));
    }, [spacesForConfluenceURL]);

    const handleEvents = useCallback((_, spaceKey) => {
        if (spaceKey !== selectedSpaceKey) {
            onSpaceKeyChange(spaceKey);
        }
    }, [selectedSpaceKey, onSpaceKeyChange]);

    return (
        <>
            <ReactSelectSetting
                name={'space'}
                label={'Space'}
                options={getSpaceOptions}
                onChange={handleEvents}
                value={
                    selectedSpaceKey ?
                        getSpaceOptions.find(
                            (option: {value: string}) =>
                                option.value === selectedSpaceKey,
                        ) :
                        null
                }
                required={true}
                theme={theme}
                addValidate={validator.addComponent}
                removeValidate={validator.removeComponent}
            />
        </>
    );
};

export default ConfluenceSpaceSelector;
