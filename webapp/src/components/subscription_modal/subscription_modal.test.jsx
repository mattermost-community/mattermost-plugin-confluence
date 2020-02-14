import React from 'react';

import {shallow} from 'enzyme';

import Constants from '../../constants';

import SubscriptionModal from './subscription_modal';

describe('components/ChannelSettingsModal', () => {
    const baseProps = {
        theme: {},
        visibility: false,
        subscription: {},
        close: jest.fn(),
        saveChannelSubscription: jest.fn().mockResolvedValue({}),
        currentChannelID: 'abcabcabcabcabc',
        editChannelSubscription: jest.fn().mockResolvedValue({}),
    };

    test('subscription modal snapshot test', async () => {
        const props = {
            ...baseProps,
            visibility: true,
        };
        const wrapper = shallow<SubscriptionModal>(
            <SubscriptionModal {...props}/>
        );
        expect(wrapper).toMatchSnapshot();
    });

    test('new subscription', async () => {
        const props = {
            ...baseProps,
            visibility: true,
        };
        const wrapper = shallow<SubscriptionModal>(
            <SubscriptionModal {...props}/>
        );
        wrapper.setState({
            alias: 'Abc',
            baseURL: 'https://test.com',
            spaceKey: 'test',
            events: Constants.CONFLUENCE_EVENTS,
            error: '',
            saving: false,
        });
        wrapper.instance().handleSubmit({preventDefault: jest.fn()});
        expect(wrapper.state().error).toBe('');
        expect(props.saveChannelSubscription).toHaveBeenCalledWith({
            alias: 'Abc',
            baseURL: 'https://test.com',
            spaceKey: 'test',
            events: Constants.CONFLUENCE_EVENTS.map((event) => event.value),
            channelID: 'abcabcabcabcabc',
        });

        expect(props.editChannelSubscription).not.toHaveBeenCalled();
    });

    test('edit subscription', async () => {
        const subscription = {
            alias: 'Abc',
            baseURL: 'https://test.com',
            spaceKey: 'test',
            events: Constants.CONFLUENCE_EVENTS,
        };
        const props = {
            ...baseProps,
            subscription,
        };
        const wrapper = shallow<SubscriptionModal>(
            <SubscriptionModal {...props}/>
        );
        wrapper.setState({
            alias: 'Abc',
            baseURL: 'https://test.com',
            spaceKey: 'test',
            events: Constants.CONFLUENCE_EVENTS,
            error: '',
            saving: false,
        });
        wrapper.instance().handleSubmit({preventDefault: jest.fn()});
        expect(wrapper.state().error).toBe('');
        expect(props.editChannelSubscription).toHaveBeenCalledWith({
            alias: 'Abc',
            baseURL: 'https://test.com',
            spaceKey: 'test',
            events: Constants.CONFLUENCE_EVENTS.map((event) => event.value),
            channelID: 'abcabcabcabcabc',
        });
        expect(props.saveChannelSubscription).not.toHaveBeenCalled();
    });

    test('subscription data clean', async () => {
        const props = {
            ...baseProps,
            visibility: true,
        };
        const wrapper = shallow<SubscriptionModal>(
            <SubscriptionModal {...props}/>
        );
        wrapper.setState({
            alias: '   Abc   ',
            baseURL: 'https://teST.com',
            spaceKey: 'test       ',
            events: Constants.CONFLUENCE_EVENTS,
            error: '',
            saving: false,
        });
        wrapper.instance().handleSubmit({preventDefault: jest.fn()});
        expect(wrapper.state().error).toBe('');
        expect(props.saveChannelSubscription).toHaveBeenCalledWith({
            alias: 'Abc',
            baseURL: 'https://test.com',
            spaceKey: 'test',
            events: Constants.CONFLUENCE_EVENTS.map((event) => event.value),
            channelID: 'abcabcabcabcabc',
        });

        expect(props.editChannelSubscription).not.toHaveBeenCalled();
    });
});
