import { combineReducers } from 'redux';

import redeem from './redeem';
import loading from './loading';
import home from './home';
import oauth2consent from './oauth2consent';
import userprofile from './userprofile';
import settings from './settings';
import coursesettings from './coursesettings';
import grades from './grades';

export default combineReducers({
  redeem,
  loading,
  home,
  oauth2consent,
  userprofile,
  settings,
  coursesettings,
  grades
});
