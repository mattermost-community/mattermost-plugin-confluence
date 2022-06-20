import React from 'react';

import {shallow} from 'enzyme';
import ConfluenceInstanceSelector from './confluence_instance_selector';
import configureStore from 'redux-mock-store'

import { Provider } from 'react-redux';

describe('components/ConfluenceInstanceSelector', () => {
    const initialState = {}
    const baseProps = {
        theme: {},
        selectedInstanceID:'test-spaceKey',
        onInstanceChange:jest.fn()
    };
    const mockStore = configureStore()
    test('confluence instance selector snapshot test', async () => {
        const props = {
            ...baseProps,
        };
        const store = mockStore(initialState)
        const wrapper = shallow(
            <Provider store= {store} >
            < ConfluenceInstanceSelector {...props}/>,
            </Provider>
        );
        expect(wrapper).toMatchSnapshot();
    });
});
