import { combineReducers } from 'redux';

import canvas from './canvas';
import filters from './filters';
import components from './components';

export default combineReducers({
  canvas,
  filters,
  components
});
