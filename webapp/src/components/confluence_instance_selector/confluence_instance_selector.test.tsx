import React from 'react';

import {shallow} from 'enzyme';

import {Provider} from 'react-redux';

import {configureStore} from 'tests/setup';

import ConfluenceInstanceSelector from '.';

describe('components/ConfluenceInstanceSelector', () => {
    const initialState = {};
    const baseProps = {
        selectedInstanceID: 'test-spaceKey',
        onInstanceChange: jest.fn(),
    };
    const mockStore = configureStore();
    test('confluence instance selector snapshot test', async () => {
        const props = {
            ...baseProps,
        };
        const store = mockStore(initialState);
        const wrapper = shallow(
            <Provider store={store} >
                <ConfluenceInstanceSelector theme={{
                    type: undefined,
                    sidebarBg: '',
                    sidebarText: '',
                    sidebarUnreadText: '',
                    sidebarTextHoverBg: '',
                    sidebarTextActiveBorder: '',
                    sidebarTextActiveColor: '',
                    sidebarHeaderBg: '',
                    sidebarHeaderTextColor: '',
                    onlineIndicator: '',
                    awayIndicator: '',
                    dndIndicator: '',
                    mentionBg: '',
                    mentionBj: '',
                    mentionColor: '',
                    centerChannelBg: '',
                    centerChannelColor: '',
                    newMessageSeparator: '',
                    linkColor: '',
                    buttonBg: '',
                    buttonColor: '',
                    errorTextColor: '',
                    mentionHighlightBg: '',
                    mentionHighlightLink: '',
                    codeTheme: ''
                }} {...props}/>
            </Provider>,
        );
        expect(wrapper).toMatchSnapshot();
    });
});
