import {connect} from 'react-redux';
import {bindActionCreators} from 'redux';
import {getCurrentChannelId} from 'mattermost-redux/selectors/entities/common';

import {closeSubscriptionModal, saveChannelSubscription, editChannelSubscription} from '../../actions';
import Selectors from '../../selectors';

import SubscriptionModal from './subscription_modal';

const mapStateToProps = (state) => {
    return {
        subscription: Selectors.isSubscriptionEditModalVisible(state),
        visibility: Selectors.isSubscriptionModalVisible(state),
        currentChannelID: getCurrentChannelId(state),
    };
};

const mapDispatchToProps = (dispatch) => bindActionCreators({
    close: closeSubscriptionModal,
    saveChannelSubscription,
    editChannelSubscription,
}, dispatch);

export default connect(mapStateToProps, mapDispatchToProps)(SubscriptionModal);
