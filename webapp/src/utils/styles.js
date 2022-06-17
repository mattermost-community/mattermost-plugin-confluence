// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import {changeOpacity} from 'mattermost-redux/utils/theme_utils';

export const getStyleForReactSelect = (theme) => {
    if (!theme) {
        return null;
    }

    return {
        menuPortal: (provided) => ({
            ...provided,
            zIndex: 9999,
        }),
        control: (provided, state) => ({
            ...provided,
            color: theme.theme.centerChannelColor,
            background: theme.theme.centerChannelBg,

            // Overwrittes the different states of border
            borderColor: state.isFocused ? changeOpacity(theme.theme.centerChannelColor, 0.25) : changeOpacity(theme.theme.centerChannelColor, 0.12),

            // Removes weird border around container
            boxShadow: 'inset 0 1px 1px ' + changeOpacity(theme.theme.centerChannelColor, 0.075),
            borderRadius: '2px',

            '&:hover': {
                borderColor: changeOpacity(theme.theme.centerChannelColor, 0.25),
            },
        }),
        option: (provided, state) => ({
            ...provided,
            background: state.isFocused ? changeOpacity(theme.theme.centerChannelColor, 0.12) : theme.theme.centerChannelBg,
            cursor: state.isDisabled ? 'not-allowed' : 'pointer',
            color: theme.theme.centerChannelColor,
            '&:hover': state.isDisabled ? {} : {
                background: changeOpacity(theme.theme.centerChannelColor, 0.12),
            },
        }),
        clearIndicator: (provided) => ({
            ...provided,
            width: '34px',
            color: changeOpacity(theme.theme.centerChannelColor, 0.4),
            transform: 'scaleX(1.15)',
            marginRight: '-10px',
            '&:hover': {
                color: theme.theme.centerChannelColor,
            },
        }),
        multiValue: (provided) => ({
            ...provided,
            background: changeOpacity(theme.theme.centerChannelColor, 0.15),
        }),
        multiValueLabel: (provided) => ({
            ...provided,
            color: theme.theme.centerChannelColor,
            paddingBottom: '4px',
            paddingLeft: '8px',
            fontSize: '90%',
        }),
        multiValueRemove: (provided) => ({
            ...provided,
            transform: 'translateX(-2px) scaleX(1.15)',
            color: changeOpacity(theme.theme.centerChannelColor, 0.4),
            '&:hover': {
                background: 'transparent',
            },
        }),
        menu: (provided) => ({
            ...provided,
            color: theme.theme.centerChannelColor,
            background: theme.theme.centerChannelBg,
            border: '1px solid ' + changeOpacity(theme.theme.centerChannelColor, 0.2),
            borderRadius: '0 0 2px 2px',
            boxShadow: changeOpacity(theme.theme.centerChannelColor, 0.2) + ' 1px 3px 12px',
            marginTop: '4px',
        }),
        input: (provided) => ({
            ...provided,
            color: theme.theme.centerChannelColor,
        }),
        placeholder: (provided) => ({
            ...provided,
            color: theme.theme.centerChannelColor,
        }),
        dropdownIndicator: (provided) => ({
            ...provided,

            '&:hover': {
                color: theme.theme.centerChannelColor,
            },
        }),
        singleValue: (provided) => ({
            ...provided,
            color: theme.theme.centerChannelColor,
        }),
        indicatorSeparator: (provided) => ({
            ...provided,
            display: 'none',
        }),
    };
};

export const getModalStyles = (theme) => ({
    modalBody: {
        padding: '2em 2em 3em',
        color: theme.centerChannelColor,
        backgroundColor: theme.centerChannelBg,
    },
    modalFooter: {
        padding: '2rem 15px',
    },
    descriptionArea: {
        height: 'auto',
        width: '100%',
        color: '#000',
    },
});
