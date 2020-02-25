import { combineReducers } from 'redux';

import canvas from './canvas';
import error from './error';
import loading from './loading';
import checkout from './checkout';
import plus from './plus';
import oauth2 from './oauth2';
import settings from './settings';
import components from './components';

export default combineReducers({
  canvas,
  error,
  loading,
  checkout,
  plus,
  oauth2,
  settings,
  components
});
