import React from 'react';

import {shallow} from 'enzyme';

import {Provider} from 'react-redux';

import {Theme} from 'mattermost-redux/types/preferences';

import {configureStore} from 'tests/setup';

import ConfluenceInstanceSelector from '.';

describe('components/ConfluenceInstanceSelector', () => {
    const initialState = {};
    const baseProps = {
        theme: {} as Theme,
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
            <Provider store={store}>
                <ConfluenceInstanceSelector {...props}/>
            </Provider>,
        );
        expect(wrapper).toMatchSnapshot();
    });
});
