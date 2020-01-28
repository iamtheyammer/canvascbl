import { combineReducers } from 'redux';

import redeem from './redeem';
import loading from './loading';
import home from './home';

export default combineReducers({
  redeem,
  loading,
  home
});
