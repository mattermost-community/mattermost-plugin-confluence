import {blendColors, changeOpacity, makeStyleFromTheme} from 'mattermost-redux/utils/theme_utils';

export const getReactSelectTheme = makeStyleFromTheme((mmTheme) => {
    return (originalTheme) => ({
        ...originalTheme,
        colors: {
            ...originalTheme.colors,
            primary: mmTheme.centerChannelColor,
            primary75: changeOpacity(mmTheme.centerChannelColor, 0.75),
            primary50: changeOpacity(mmTheme.centerChannelColor, 0.50),
            primary25: changeOpacity(mmTheme.centerChannelColor, 0.25),
            danger: mmTheme.errorTextColor,
            dangerLight: changeOpacity(mmTheme.errorTextColor, 0.50),
            neutral0: blendColors(mmTheme.centerChannelBg, mmTheme.centerChannelColor, 0),
            neutral5: blendColors(mmTheme.centerChannelBg, mmTheme.centerChannelColor, 0.05),
            neutral10: blendColors(mmTheme.centerChannelBg, mmTheme.centerChannelColor, 0.10),
            neutral20: blendColors(mmTheme.centerChannelBg, mmTheme.centerChannelColor, 0.20),
            neutral30: blendColors(mmTheme.centerChannelBg, mmTheme.centerChannelColor, 0.30),
            neutral40: blendColors(mmTheme.centerChannelBg, mmTheme.centerChannelColor, 0.40),
            neutral50: blendColors(mmTheme.centerChannelBg, mmTheme.centerChannelColor, 0.50),
            neutral60: blendColors(mmTheme.centerChannelBg, mmTheme.centerChannelColor, 0.60),
            neutral70: blendColors(mmTheme.centerChannelBg, mmTheme.centerChannelColor, 0.70),
            neutral80: blendColors(mmTheme.centerChannelBg, mmTheme.centerChannelColor, 0.80),
            neutral90: blendColors(mmTheme.centerChannelBg, mmTheme.centerChannelColor, 0.90),
        },
    });
});

// From: https://github.com/JedWatson/react-select/wiki/v2:-Styles-to-match-react-bootstrap-fields
export const reactSelectStyles = {
    menuPortal: (provided) => ({
        ...provided,
        zIndex: 9999,
    }),
    container: (styles) => ({
        ...styles,
        flex: 1,
    }),
    control: (styles) => ({
        ...styles,
        minHeight: '34px',
    }),
    placeholder: (styles, state) => ({
        display: state.selectProps.menuIsOpen ? 'none' : 'inline',
        paddingLeft: 3,
    }),
    clearIndicator: (styles) => ({
        ...styles,
        padding: '2px 8px',
    }),
    indicatorSeparator: () => ({
        display: 'none',
    }),
    dropdownIndicator: (styles) => ({
        ...styles,
        padding: '2px 8px',
    }),
    loadingIndicator: (styles) => ({
        ...styles,
        padding: '2px 8px',
    }),
    menu: (styles) => ({
        ...styles,
        zIndex: 3, // Without this menu will be overlapped by other fields
    }),
};
