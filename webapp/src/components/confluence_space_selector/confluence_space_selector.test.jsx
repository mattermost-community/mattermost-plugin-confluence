/*eslint max-nested-callbacks: ["error", 3]*/

import React from 'react';

import {shallow} from 'enzyme';

import configureStore from 'redux-mock-store';

import {Provider} from 'react-redux';

import ConfluenceSpaceSelector from './confluence_space_selector';

describe('components/ConfluenceSpaceSelector', () => {
    const initialState = {};
    const baseProps = {
        theme: {},
        selectedSpaceKey: 'test-spaceKey',
        onSpaceKeyChange: jest.fn(),
    };
    const mockStore = configureStore();
    test('confluence space selector snapshot test', async () => {
        const props = {
            ...baseProps,
        };
        const store = mockStore(initialState);
        const wrapper = shallow(
            <Provider store={store} >
                <ConfluenceSpaceSelector {...props}/>
            </Provider>,
        );
        expect(wrapper).toMatchSnapshot();
    });
});