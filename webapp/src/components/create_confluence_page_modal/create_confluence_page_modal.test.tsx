/*eslint max-nested-callbacks: ["error", 3]*/

import React from 'react';

import {shallow} from 'enzyme';

import {Provider} from 'react-redux';

import {Theme} from 'mattermost-redux/types/preferences';

import {configureStore} from 'tests/setup';

import CreateConfluencePage from '.';

describe('components/CreateCofluencePageModal', () => {
    const initialState = {
        message: 'test-message',
    };
    const baseProps = {
        theme: {} as Theme,
    };
    const mockStore = configureStore();
    test('confluence create page modal snapshot test', async () => {
        const props = {
            ...baseProps,
        };
        const store = mockStore(initialState);
        const wrapper = shallow(
            <Provider store={store} >
                <CreateConfluencePage {...props.theme}/>
            </Provider>,
        );
        expect(wrapper).toMatchSnapshot();
    });
});
