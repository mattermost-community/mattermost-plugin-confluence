import {connect} from 'react-redux';
import {bindActionCreators} from 'redux';
import {getCurrentChannelId} from 'mattermost-redux/selectors/entities/common';

import {closeConfigModal, saveChannelSubscription} from '../../actions';
import Selectors from '../../selectors';

import ConfigModal from './config_modal';

const mapStateToProps = (state) => {
    return {
        visibility: Selectors.isConfigModalvisible(state),
        currentChannelID: getCurrentChannelId(state),
    };
};

const mapDispatchToProps = (dispatch) => bindActionCreators({
    close: closeConfigModal,
    saveChannelSubscription,
}, dispatch);

export default connect(mapStateToProps, mapDispatchToProps)(ConfigModal);
