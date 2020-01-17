import {connect} from 'react-redux';
import {bindActionCreators} from 'redux';

import {closeConfigModal} from '../../actions';
import Selectors from '../../selectors';

import ConfigModal from './config_modal';

const mapStateToProps = (state) => {
    return {
        visibility: Selectors.isConfigModalvisible(state),
    };
};

const mapDispatchToProps = (dispatch) => bindActionCreators({
    close: closeConfigModal,
}, dispatch);

export default connect(mapStateToProps, mapDispatchToProps)(ConfigModal);
