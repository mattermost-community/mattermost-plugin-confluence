/*eslint max-nested-callbacks: ["error", 3]*/

import React from 'react';

import {shallow} from 'enzyme';

import {Provider} from 'react-redux';

import {configureStore} from 'tests/setup';

import CreateConfluencePage from '.';

describe('components/CreateCofluencePageModal', () => {
    const initialState = {
        message: 'test-message',
    };
    const baseProps = {
        theme: {},
    };
    const mockStore = configureStore();
    test('confluence create page modal snapshot test', async () => {
        const props = {
            ...baseProps,
        };
        const store = mockStore(initialState);
        const wrapper = shallow(
            <Provider store={store} >
                <CreateConfluencePage {...props}/>
            </Provider>,
        );
        expect(wrapper).toMatchSnapshot();
    });
});
