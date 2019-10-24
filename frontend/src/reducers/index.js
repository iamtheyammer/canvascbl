import { combineReducers } from 'redux';

import canvas from './canvas';
import error from './error';
import loading from './loading';
import checkout from './checkout';
import plus from './plus';

export default combineReducers({
  canvas,
  error,
  loading,
  checkout,
  plus
});
